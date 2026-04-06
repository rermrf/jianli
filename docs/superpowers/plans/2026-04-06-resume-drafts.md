# Resume Drafts Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add admin-only, backend-persisted resume draft versions that can be created from the edit page, listed in a drafts page, previewed, published to the current resume, and deleted.

**Architecture:** Keep the existing `resume` record as the single public source of truth and add a separate `resume_drafts` persistence path for versioned drafts. Expose draft CRUD/publish APIs behind the existing admin auth middleware, then layer focused frontend clients, hooks, and pages on top so editing, list management, and preview stay isolated.

**Tech Stack:** Go, Gin, GORM, SQLite, React, TypeScript, React Router, Vitest, Testing Library

---

## File Structure

### Backend files

- Create: `internal/model/resume_draft.go`
  Responsibility: Define the GORM model for persisted draft versions and keep draft schema isolated from the main resume record.
- Create: `internal/store/resume_draft_store.go`
  Responsibility: Encapsulate draft persistence operations: create, list, get by id, publish to main resume, delete.
- Create: `internal/store/resume_draft_store_test.go`
  Responsibility: Prove the store API covers the draft lifecycle and publish behavior.
- Create: `internal/handler/resume_draft.go`
  Responsibility: Bind and validate HTTP requests for draft create/list/detail/publish/delete endpoints.
- Create: `internal/handler/resume_draft_test.go`
  Responsibility: Verify auth, status codes, payload shape, and publish side effects through Gin handlers.
- Modify: `internal/store/sqlite.go`
  Responsibility: Add the draft model to `AutoMigrate`.
- Modify: `cmd/server/main.go`
  Responsibility: Wire the new draft store/handler into authenticated routes.

### Frontend files

- Create: `web/src/types/resumeDraft.ts`
  Responsibility: Hold draft summary/detail/request types used by the pages and API layer.
- Create: `web/src/lib/resumeDrafts.ts`
  Responsibility: Wrap draft HTTP requests so `storage.ts` can stay focused on the current resume.
- Create: `web/src/hooks/useResumeDraftList.ts`
  Responsibility: Load, refresh, publish, and delete draft summaries for the list page.
- Create: `web/src/hooks/useResumeDraftDetail.ts`
  Responsibility: Load one draft record and expose publish/delete actions for the preview page.
- Create: `web/src/components/editor/SaveDraftDialog.tsx`
  Responsibility: Collect draft name and note without bloating `EditPage.tsx`.
- Create: `web/src/pages/DraftsPage.tsx`
  Responsibility: Render the admin-only draft list and row actions.
- Create: `web/src/pages/DraftPreviewPage.tsx`
  Responsibility: Render one draft preview with publish navigation.
- Create: `web/src/pages/__tests__/drafts-page.test.tsx`
  Responsibility: Cover the draft list, preview, publish, delete, and protected-route flows.
- Modify: `web/src/pages/EditPage.tsx`
  Responsibility: Add “保存为草稿” flow while keeping “保存主简历” intact.
- Modify: `web/src/pages/__tests__/edit-page.test.tsx`
  Responsibility: Add a failing test for the save-draft flow and preserve existing save-main behavior coverage.
- Modify: `web/src/app/router.tsx`
  Responsibility: Register `/drafts` and `/drafts/:id` under `ProtectedRoute`.
- Modify: `web/src/components/layout/TopNav.tsx`
  Responsibility: Surface the drafts page in authenticated navigation.
- Modify: `web/src/pages/__tests__/routing.test.tsx`
  Responsibility: Verify drafts routes are protected and visible to admins.

### Existing code to preserve

- Preserve `GET /api/resume`, `PUT /api/resume`, and `GET /api/resume/pdf` behavior.
- Preserve `web/src/hooks/useResumeDraft.ts` as the loader for the editable current resume unless a test proves a rename/refactor is necessary.
- Preserve `ResumeDesktopLayout` and `ResumeMobileLayout` as the single rendering surfaces reused by public resume and draft preview pages.

## Chunk 1: Backend Draft Persistence And API

### Task 1: Add draft persistence with store-level publish support

**Files:**
- Create: `internal/model/resume_draft.go`
- Create: `internal/store/resume_draft_store.go`
- Create: `internal/store/resume_draft_store_test.go`
- Modify: `internal/store/sqlite.go`
- Reuse: `internal/store/resume_store.go`

- [ ] **Step 1: Write the failing store tests**

Add `TestResumeDraftStoreCreateAndList`, `TestResumeDraftStoreGetByID`, `TestResumeDraftStorePublishReplacesMainResume`, and `TestResumeDraftStoreDelete` in `internal/store/resume_draft_store_test.go`.

```go
func TestResumeDraftStorePublishReplacesMainResume(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "resume.db"))
	if err != nil {
		t.Fatalf("Open() returned error: %v", err)
	}

	draftStore := NewResumeDraftStore(db, NewResumeStore(db))
	draftID, err := draftStore.Create("面试版", "补充项目亮点", json.RawMessage(`{"profile":{"name":"草稿版"}}`))
	if err != nil {
		t.Fatalf("Create() returned error: %v", err)
	}

	if err := draftStore.Publish(draftID); err != nil {
		t.Fatalf("Publish() returned error: %v", err)
	}

	mainResume, err := NewResumeStore(db).Get()
	if err != nil {
		t.Fatalf("Get() returned error: %v", err)
	}

	if string(mainResume) != `{"profile":{"name":"草稿版"}}` {
		t.Fatalf("expected main resume to be replaced, got %s", string(mainResume))
	}
}
```

- [ ] **Step 2: Run the store tests to verify they fail**

Run: `go test ./internal/store -run TestResumeDraftStore -count=1`

Expected: FAIL because `ResumeDraftStore`, draft model, and methods do not exist yet.

- [ ] **Step 3: Add the draft model and migrate it**

Implement `ResumeDraftRecord` in `internal/model/resume_draft.go` with `ID`, `Name`, `Note`, `Data`, `CreatedAt`, and `UpdatedAt`, then register it in `internal/store/sqlite.go`.

```go
type ResumeDraftRecord struct {
	ID        int64     `gorm:"primaryKey"`
	Name      string    `gorm:"column:name;type:TEXT;not null"`
	Note      string    `gorm:"column:note;type:TEXT;not null;default:''"`
	Data      []byte    `gorm:"column:data;type:TEXT;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (ResumeDraftRecord) TableName() string {
	return "resume_drafts"
}
```

- [ ] **Step 4: Implement the draft store with explicit methods**

Implement `Create`, `List`, `Get`, `Publish`, and `Delete` in `internal/store/resume_draft_store.go`. Keep `Publish` responsible for copying draft JSON into `ResumeStore.Save` so the publish semantics stay in one place.

```go
type ResumeDraftStore struct {
	db          *gorm.DB
	resumeStore *ResumeStore
}

func (s *ResumeDraftStore) Publish(id int64) error {
	draft, err := s.Get(id)
	if err != nil {
		return err
	}

	return s.resumeStore.Save(draft.Data)
}
```

- [ ] **Step 5: Re-run the store tests and keep iterating until green**

Run: `go test ./internal/store -run TestResumeDraftStore -count=1`

Expected: PASS for the new draft store tests.

- [ ] **Step 6: Run the broader store regression tests**

Run: `go test ./internal/store -count=1`

Expected: PASS with existing visitor and resume store tests still green.

- [ ] **Step 7: Commit the persistence slice**

```bash
git add internal/model/resume_draft.go internal/store/resume_draft_store.go internal/store/resume_draft_store_test.go internal/store/sqlite.go
git commit -m "feat: add resume draft persistence"
```

### Task 2: Expose authenticated draft endpoints and publish behavior

**Files:**
- Create: `internal/handler/resume_draft.go`
- Create: `internal/handler/resume_draft_test.go`
- Modify: `cmd/server/main.go`
- Reuse: `internal/httpapi/response.go`
- Reuse: `internal/middleware/auth.go`

- [ ] **Step 1: Write the failing handler tests**

Add tests for:
- auth required on all draft routes
- `POST /api/resume/drafts` rejects empty `name`
- `GET /api/resume/drafts` returns summaries in `updated_at desc`
- `GET /api/resume/drafts/:id` returns one draft
- `PUT /api/resume/drafts/:id/publish` updates `GET /api/resume`
- `DELETE /api/resume/drafts/:id` removes the draft
- missing draft routes return `404`

```go
func TestResumeDraftHandlerPublishUpdatesMainResume(t *testing.T) {
	router := setupResumeDraftRouter(t)

	create := httptest.NewRequest(http.MethodPost, "/api/resume/drafts", strings.NewReader(`{
		"name":"发布版",
		"note":"准备用于线上",
		"data":{"profile":{"name":"已发布简历"}}
	}`))
	create.Header.Set("Content-Type", "application/json")
	create.Header.Set("X-Auth-Key", "resume-key")
	createRecorder := httptest.NewRecorder()
	router.ServeHTTP(createRecorder, create)

	publish := httptest.NewRequest(http.MethodPut, "/api/resume/drafts/1/publish", nil)
	publish.Header.Set("X-Auth-Key", "resume-key")
	publishRecorder := httptest.NewRecorder()
	router.ServeHTTP(publishRecorder, publish)

	mainResume := httptest.NewRecorder()
	router.ServeHTTP(mainResume, httptest.NewRequest(http.MethodGet, "/api/resume", nil))

	if !strings.Contains(mainResume.Body.String(), "已发布简历") {
		t.Fatalf("expected main resume body to contain published draft, got %s", mainResume.Body.String())
	}
}
```

- [ ] **Step 2: Run the handler tests to verify they fail**

Run: `go test ./internal/handler -run TestResumeDraftHandler -count=1`

Expected: FAIL because the draft handler, routes, and request parsing do not exist yet.

- [ ] **Step 3: Implement the draft handler**

Create `ResumeDraftHandler` with explicit request/response structs so list and detail payloads stay stable.

```go
type createResumeDraftRequest struct {
	Name string          `json:"name"`
	Note string          `json:"note"`
	Data json.RawMessage `json:"data"`
}

type ResumeDraftHandler struct {
	store *store.ResumeDraftStore
}
```

Return:
- `400` for invalid body or empty `name`
- `404` for missing draft id
- `401` through existing auth middleware

- [ ] **Step 4: Register the new routes under the authenticated `/api` group**

Update `cmd/server/main.go` so only admins can access the new endpoints.

```go
draftHandler := handler.NewResumeDraftHandler(draftStore)

protected.POST("/resume/drafts", draftHandler.Create)
protected.GET("/resume/drafts", draftHandler.List)
protected.GET("/resume/drafts/:id", draftHandler.Get)
protected.PUT("/resume/drafts/:id/publish", draftHandler.Publish)
protected.DELETE("/resume/drafts/:id", draftHandler.Delete)
```

- [ ] **Step 5: Re-run the focused handler tests**

Run: `go test ./internal/handler -run TestResumeDraftHandler -count=1`

Expected: PASS with draft route behavior verified.

- [ ] **Step 6: Run backend regression tests**

Run: `go test ./internal/... -count=1`

Expected: PASS across handler, middleware, store, and config packages.

- [ ] **Step 7: Commit the API slice**

```bash
git add internal/handler/resume_draft.go internal/handler/resume_draft_test.go cmd/server/main.go
git commit -m "feat: add resume draft api"
```

## Chunk 2: Frontend Draft Creation, Management, And Preview

### Task 3: Add the save-draft client API and edit-page flow

**Files:**
- Create: `web/src/types/resumeDraft.ts`
- Create: `web/src/lib/resumeDrafts.ts`
- Create: `web/src/components/editor/SaveDraftDialog.tsx`
- Modify: `web/src/pages/EditPage.tsx`
- Modify: `web/src/pages/__tests__/edit-page.test.tsx`
- Reuse: `web/src/lib/api.ts`
- Reuse: `web/src/lib/auth.ts`

- [ ] **Step 1: Write the failing edit-page test for saving a draft**

Add a test that opens the save-draft UI, fills `name` and `note`, confirms the action, and expects a `POST /api/resume/drafts` call with the current form data plus admin auth header.

```tsx
it('saves the current edit state as a named draft', async () => {
  loginWithKey('resume-key')
  const fetchSpy = vi.spyOn(globalThis, 'fetch').mockImplementation(async (input, init) => {
    const url = String(input)

    if (url === '/api/resume' && (!init || init.method === 'GET')) {
      return new Response(JSON.stringify({ code: 0, data: defaultResume }), { status: 200 })
    }

    if (url === '/api/resume/drafts' && init?.method === 'POST') {
      return new Response(JSON.stringify({ code: 0, data: { id: 8, name: '面试前调整版', note: '补充项目亮点', data: JSON.parse(String(init.body))?.data } }), { status: 200 })
    }

    return new Response(JSON.stringify({ code: 0, data: [] }), { status: 200 })
  })

  renderEditPage()

  const user = userEvent.setup()
  await user.click(await screen.findByRole('button', { name: '保存为草稿' }))
  await user.type(screen.getByLabelText('草稿名称'), '面试前调整版')
  await user.type(screen.getByLabelText('草稿备注'), '补充项目亮点')
  await user.click(screen.getByRole('button', { name: '确认保存草稿' }))

  await waitFor(() =>
    expect(fetchSpy).toHaveBeenCalledWith(
      '/api/resume/drafts',
      expect.objectContaining({ method: 'POST' }),
    ),
  )
})
```

- [ ] **Step 2: Run the focused edit-page test to verify it fails**

Run: `npm.cmd run test -- --runInBand web/src/pages/__tests__/edit-page.test.tsx`

Expected: FAIL because the save-draft UI and API client do not exist yet.

- [ ] **Step 3: Add draft types and API helpers**

Create `web/src/types/resumeDraft.ts` and `web/src/lib/resumeDrafts.ts`.

```ts
export interface ResumeDraftSummary {
  id: number
  name: string
  note: string
  createdAt: string
  updatedAt: string
}

export interface ResumeDraftDetail extends ResumeDraftSummary {
  data: ResumeData
}

export async function createResumeDraft(input: CreateResumeDraftInput) {
  return apiFetch<ResumeDraftDetail>('/api/resume/drafts', {
    method: 'POST',
    headers: { 'X-Auth-Key': getAuthKey() ?? '' },
    body: JSON.stringify(input),
  })
}
```

- [ ] **Step 4: Add a focused save-draft dialog and wire it into the edit page**

Create `SaveDraftDialog.tsx` for the name/note form and keep `EditPage.tsx` responsible only for opening it, passing `draft`, and showing success/error state.

Implementation constraints:
- keep existing “保存” main-resume path intact, but relabel its visible copy to “保存主简历”
- do not reuse the main save toast text for draft saves; show a distinct success message
- keep mobile and desktop action groups functionally equivalent

- [ ] **Step 5: Re-run the focused edit-page tests**

Run: `npm.cmd run test -- --runInBand web/src/pages/__tests__/edit-page.test.tsx`

Expected: PASS for the new save-draft flow and existing edit-page assertions.

- [ ] **Step 6: Commit the edit-page draft flow**

```bash
git add web/src/types/resumeDraft.ts web/src/lib/resumeDrafts.ts web/src/components/editor/SaveDraftDialog.tsx web/src/pages/EditPage.tsx web/src/pages/__tests__/edit-page.test.tsx
git commit -m "feat: add save resume draft flow"
```

### Task 4: Add protected drafts list and preview pages

**Files:**
- Create: `web/src/hooks/useResumeDraftList.ts`
- Create: `web/src/hooks/useResumeDraftDetail.ts`
- Create: `web/src/pages/DraftsPage.tsx`
- Create: `web/src/pages/DraftPreviewPage.tsx`
- Create: `web/src/pages/__tests__/drafts-page.test.tsx`
- Modify: `web/src/app/router.tsx`
- Modify: `web/src/components/layout/TopNav.tsx`
- Modify: `web/src/pages/__tests__/routing.test.tsx`
- Reuse: `web/src/components/resume/ResumeDesktopLayout.tsx`
- Reuse: `web/src/components/resume/ResumeMobileLayout.tsx`

- [ ] **Step 1: Write the failing routing and drafts-page tests**

Add tests that verify:
- authenticated nav includes “草稿”
- unauthenticated `/drafts` and `/drafts/:id` redirect to `/login`
- `/drafts` renders server-returned draft summaries
- clicking “预览” navigates to `/drafts/:id`
- clicking “设为主简历” sends `PUT /api/resume/drafts/:id/publish`
- clicking “删除” sends `DELETE /api/resume/drafts/:id`

```tsx
it('publishes a draft from the drafts list', async () => {
  loginWithKey('resume-key')
  const user = userEvent.setup()
  const fetchSpy = mockDraftFetch()

  renderAtPath('/drafts')

  await user.click(await screen.findByRole('button', { name: '设为主简历 面试前调整版' }))

  await waitFor(() =>
    expect(fetchSpy).toHaveBeenCalledWith(
      '/api/resume/drafts/8/publish',
      expect.objectContaining({ method: 'PUT' }),
    ),
  )
})
```

- [ ] **Step 2: Run the focused frontend tests to verify they fail**

Run: `npm.cmd run test -- --runInBand web/src/pages/__tests__/drafts-page.test.tsx web/src/pages/__tests__/routing.test.tsx`

Expected: FAIL because the routes, hooks, pages, and nav item do not exist yet.

- [ ] **Step 3: Implement draft list/detail hooks and pages**

Keep data loading out of the page components:

```ts
export function useResumeDraftList() {
  const [drafts, setDrafts] = useState<ResumeDraftSummary[]>([])

  async function publishDraft(id: number) {
    await publishResumeDraft(id)
    await refresh()
  }

  return { drafts, error, loading, refresh, publishDraft, deleteDraft }
}
```

Page behavior:
- `DraftsPage.tsx` renders summaries plus row actions
- `DraftPreviewPage.tsx` loads one draft and feeds `draft.data` into the existing resume layouts
- publish success should navigate back to `/drafts` or stay in place with explicit success text; choose one behavior and cover it in tests

- [ ] **Step 4: Register routes and expose the drafts entry in admin navigation**

Update `web/src/app/router.tsx` and `web/src/components/layout/TopNav.tsx`.

```tsx
<Route
  path="/drafts"
  element={
    <ProtectedRoute>
      <DraftsPage />
    </ProtectedRoute>
  }
/>
<Route
  path="/drafts/:id"
  element={
    <ProtectedRoute>
      <DraftPreviewPage />
    </ProtectedRoute>
  }
/>
```

- [ ] **Step 5: Re-run the focused drafts and routing tests**

Run: `npm.cmd run test -- --runInBand web/src/pages/__tests__/drafts-page.test.tsx web/src/pages/__tests__/routing.test.tsx`

Expected: PASS with authenticated routing and draft-management flows covered.

- [ ] **Step 6: Run the full frontend regression suite**

Run: `npm.cmd run test -- --runInBand`

Expected: PASS across edit, routing, resume, print, and drafts page tests.

- [ ] **Step 7: Commit the drafts UI slice**

```bash
git add web/src/hooks/useResumeDraftList.ts web/src/hooks/useResumeDraftDetail.ts web/src/pages/DraftsPage.tsx web/src/pages/DraftPreviewPage.tsx web/src/pages/__tests__/drafts-page.test.tsx web/src/app/router.tsx web/src/components/layout/TopNav.tsx web/src/pages/__tests__/routing.test.tsx
git commit -m "feat: add resume draft management pages"
```

## Chunk 3: Final Verification

### Task 5: Verify the complete draft workflow end to end

**Files:**
- Modify only if regressions are found while verifying: `internal/...`, `web/src/...`
- Reference: `docs/superpowers/specs/2026-04-06-resume-drafts-design.md`

- [ ] **Step 1: Run all backend tests**

Run: `go test ./... -count=1`

Expected: PASS with no new failures outside the draft feature.

- [ ] **Step 2: Run all frontend tests**

Run: `npm.cmd run test -- --runInBand`

Expected: PASS.

- [ ] **Step 3: Run the frontend production build**

Run: `npm.cmd run build`

Expected: PASS and emit a Vite production bundle without type errors.

- [ ] **Step 4: Manually verify the admin workflow**

Manual checks:
- log in as admin and open `/edit`
- change the resume, save a named draft, and confirm `/` still shows the old main resume
- open `/drafts`, preview the new draft, publish it, and confirm `/` now shows the draft content
- confirm the published draft still exists in `/drafts`
- delete the draft and confirm it disappears from the list

- [ ] **Step 5: Commit any verification fixes**

```bash
git add .
git commit -m "test: verify resume draft workflow"
```

## Notes For The Implementer

- Follow @test-driven-development discipline exactly for each task: write the failing test first, confirm the failure reason, then add the minimum code.
- Do not collapse main-resume and draft APIs into one endpoint; the separation is the core boundary that keeps the public resume stable.
- Keep new frontend helpers narrow. Avoid putting list, detail, publish, and delete logic back into `EditPage.tsx`.
- If a route or page starts accumulating too much state, split view logic into a local component rather than expanding one file further.

## Local Review Checklist

- The new draft store owns draft lifecycle operations and publish semantics.
- All draft APIs sit behind existing admin auth middleware.
- The public resume page still reads only from `GET /api/resume`.
- Draft preview reuses existing resume rendering components instead of introducing a second layout system.
- Tests cover both success paths and not-found/auth failures.
