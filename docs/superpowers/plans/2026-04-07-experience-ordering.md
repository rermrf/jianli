# Experience Ordering Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add up/down reordering for work experience, education experience, and project experience in the resume editor, with saved order reflected consistently across all displays.

**Architecture:** Keep ordering encoded directly in the existing `ResumeData` arrays and implement reordering as pure array swaps in the editor. Extend the editable list item header to expose reusable action buttons, then verify persistence by asserting the saved payload order and one display-layer rendering order.

**Tech Stack:** React, TypeScript, Vitest, Testing Library, Vite

---

## File Structure

### Primary files

- Modify: `web/src/components/editor/EditableListSection.tsx`
  Responsibility: Allow each editable item header to render optional custom actions like `上移` and `下移` without embedding list-specific business logic in the shared component.
- Modify: `web/src/pages/EditPage.tsx`
  Responsibility: Own the reorder helpers for `workExperience`, `education`, and `projects`, wire move actions into each list item, and keep save-main/save-draft payloads based on the reordered arrays.
- Modify: `web/src/pages/__tests__/edit-page.test.tsx`
  Responsibility: Add regression coverage for reorder interaction, disabled states, and persisted payload order for main save and draft save.

### Display verification files

- Modify: `web/src/pages/__tests__/routing.test.tsx`
  Responsibility: Keep existing routing behavior intact and add one focused assertion proving a reordered array is rendered in the same order on the public resume page.
- Optional reference only: `web/src/components/resume/ResumeDesktopLayout.tsx`
- Optional reference only: `web/src/components/resume/ResumeMobileLayout.tsx`
- Optional reference only: `web/src/components/resume/PrintResume.tsx`

No backend files should change for this feature.

## Chunk 1: Editor Reordering Controls

### Task 1: Add reusable move-action support to editable list items

**Files:**
- Modify: `web/src/components/editor/EditableListSection.tsx`
- Test through: `web/src/pages/__tests__/edit-page.test.tsx`

- [ ] **Step 1: Write the failing editor test for disabled move buttons**

Add a test in `web/src/pages/__tests__/edit-page.test.tsx` that loads the editor and asserts:
- the first education entry has a disabled `上移`
- the last project entry has a disabled `下移`

Example expectation shape:

```tsx
const moveUpButtons = await screen.findAllByRole('button', { name: '上移' })
expect(moveUpButtons[0]).toBeDisabled()

const moveDownProjectButton = screen.getByRole('button', { name: '下移 项目经历 3' })
expect(moveDownProjectButton).toBeDisabled()
```

- [ ] **Step 2: Run the focused test to verify it fails**

Run: `npm.cmd run test -- src/pages/__tests__/edit-page.test.tsx`
Expected: FAIL because move controls do not exist yet.

- [ ] **Step 3: Extend `EditableListItem` to accept optional header actions**

Update `web/src/components/editor/EditableListSection.tsx` so `EditableListItem` can render an optional `actions` prop alongside the existing delete button.

Implementation notes:
- keep delete behavior unchanged
- do not hardcode move semantics into the shared component
- render actions before delete so list-specific controls stay grouped

Example target shape:

```tsx
interface EditableListItemProps {
  actions?: ReactNode
  children: ReactNode
  onRemove: () => void
  title: string
}
```

- [ ] **Step 4: Re-run the focused test and confirm it still fails for the expected next reason**

Run: `npm.cmd run test -- src/pages/__tests__/edit-page.test.tsx`
Expected: FAIL because `EditPage.tsx` still does not pass move actions.

- [ ] **Step 5: Commit the shared component groundwork**

```bash
git add web/src/components/editor/EditableListSection.tsx web/src/pages/__tests__/edit-page.test.tsx
git commit -m "refactor: allow custom editable list item actions"
```

### Task 2: Implement reorder helpers and wire move buttons into the editor

**Files:**
- Modify: `web/src/pages/EditPage.tsx`
- Modify: `web/src/pages/__tests__/edit-page.test.tsx`

- [ ] **Step 1: Add a failing test for swapping work experience order**

In `web/src/pages/__tests__/edit-page.test.tsx`, add a test that:
- renders the editor with at least two work experiences
- clicks the first work experience `下移`
- asserts the second item becomes first in the DOM order

Example test structure:

```tsx
const headingsBefore = await screen.findAllByRole('heading', { name: /工作经历 \d/ })
expect(headingsBefore[0]).toHaveTextContent('工作经历 1')

await user.click(screen.getByRole('button', { name: '下移 工作经历 1' }))

const companyInputs = screen.getAllByLabelText('公司名称')
expect(companyInputs[0]).toHaveValue('第二条经历的公司')
```

- [ ] **Step 2: Run the focused test to verify it fails**

Run: `npm.cmd run test -- src/pages/__tests__/edit-page.test.tsx`
Expected: FAIL because the move logic does not exist yet.

- [ ] **Step 3: Add pure reorder helpers inside `EditPage.tsx`**

Implement a small helper such as:

```ts
function moveArrayItem<T>(items: T[], fromIndex: number, toIndex: number) {
  if (
    fromIndex < 0 ||
    toIndex < 0 ||
    fromIndex >= items.length ||
    toIndex >= items.length ||
    fromIndex === toIndex
  ) {
    return items
  }

  const nextItems = [...items]
  const [item] = nextItems.splice(fromIndex, 1)
  nextItems.splice(toIndex, 0, item)
  return nextItems
}
```

- [ ] **Step 4: Wire move buttons into all three experience lists**

For `education`, `workExperience`, and `projects`:
- pass `actions` into `EditableListItem`
- render `上移` and `下移` buttons
- disable them on the first/last items respectively
- update the corresponding array in state using the helper

Example action pattern:

```tsx
<Button
  aria-label={`上移 工作经历 ${index + 1}`}
  disabled={index === 0}
  onClick={() =>
    setDraft({
      ...draft,
      workExperience: moveArrayItem(draft.workExperience, index, index - 1),
    })
  }
  variant="ghost"
>
  上移
</Button>
```

- [ ] **Step 5: Re-run the focused editor tests**

Run: `npm.cmd run test -- src/pages/__tests__/edit-page.test.tsx`
Expected: PASS for move interactions and disabled-state assertions.

- [ ] **Step 6: Commit the editor reorder feature**

```bash
git add web/src/pages/EditPage.tsx web/src/pages/__tests__/edit-page.test.tsx
git commit -m "feat: add experience reordering in editor"
```

## Chunk 2: Persistence And Display Verification

### Task 3: Verify reordered arrays are saved in the main resume and draft payloads

**Files:**
- Modify: `web/src/pages/__tests__/edit-page.test.tsx`

- [ ] **Step 1: Add a failing test for main-save payload order**

Extend `web/src/pages/__tests__/edit-page.test.tsx` with a test that:
- reorders one list, ideally `projects`
- clicks `保存主简历`
- asserts the `PUT /api/resume` request body sends the new array order

Example assertion:

```tsx
await waitFor(() =>
  expect(fetchSpy).toHaveBeenCalledWith(
    '/api/resume',
    expect.objectContaining({
      body: expect.stringMatching(/第二个项目.*第一个项目/s),
      method: 'PUT',
    }),
  ),
)
```

- [ ] **Step 2: Add a failing test for draft-save payload order**

Add a sibling test that reorders one list, opens `保存为草稿`, confirms it, and asserts the `POST /api/resume/drafts` payload preserves that new order.

- [ ] **Step 3: Run the focused tests to verify they fail**

Run: `npm.cmd run test -- src/pages/__tests__/edit-page.test.tsx`
Expected: FAIL until the payload assertions match real saved ordering.

- [ ] **Step 4: Adjust implementation only if needed**

If the tests fail because save handlers are reading stale state, fix `EditPage.tsx` so both `saveDraft()` and `createResumeDraft()` use the current reordered arrays.

Likely minimal fix: no additional code if the handlers already read `draft` directly.

- [ ] **Step 5: Re-run the focused tests and confirm green**

Run: `npm.cmd run test -- src/pages/__tests__/edit-page.test.tsx`
Expected: PASS.

- [ ] **Step 6: Commit the persistence verification slice**

```bash
git add web/src/pages/__tests__/edit-page.test.tsx web/src/pages/EditPage.tsx
git commit -m "test: verify reordered experiences persist"
```

### Task 4: Add one display-layer regression proving render order follows the arrays

**Files:**
- Modify: `web/src/pages/__tests__/routing.test.tsx`
- Reference: `web/src/components/resume/ResumeDesktopLayout.tsx`
- Reference: `web/src/components/resume/ResumeMobileLayout.tsx`
- Reference: `web/src/components/resume/PrintResume.tsx`

- [ ] **Step 1: Write the failing display-order regression test**

Add one focused test to `web/src/pages/__tests__/routing.test.tsx` that renders `/` with a custom `resumeOverride` where one array is intentionally reversed, then asserts the rendered order matches that array.

Use a low-ambiguity target such as `projects` with unique names.

Example:

```tsx
mockAppFetch({
  ...defaultResume,
  projects: [
    { ...defaultResume.projects[1], name: '项目 B' },
    { ...defaultResume.projects[0], name: '项目 A' },
  ],
})
renderAtPath('/')

const projectAText = await screen.findAllByText('项目 A')
const projectBText = await screen.findAllByText('项目 B')
expect(projectBText[0].compareDocumentPosition(projectAText[0]) & Node.DOCUMENT_POSITION_FOLLOWING).toBeTruthy()
```

- [ ] **Step 2: Run the focused routing test to verify it fails if assumptions are wrong**

Run: `npm.cmd run test -- src/pages/__tests__/routing.test.tsx`
Expected: FAIL only if the display order is not actually following the array order.

- [ ] **Step 3: Fix only if the test reveals hidden sorting logic**

If it fails, remove any implicit sorting in display components. If it passes immediately, keep production code unchanged and preserve the regression test as proof.

- [ ] **Step 4: Re-run the focused routing test**

Run: `npm.cmd run test -- src/pages/__tests__/routing.test.tsx`
Expected: PASS.

- [ ] **Step 5: Commit the display verification slice**

```bash
git add web/src/pages/__tests__/routing.test.tsx
if git diff --quiet -- web/src/components/resume/ResumeDesktopLayout.tsx web/src/components/resume/ResumeMobileLayout.tsx web/src/components/resume/PrintResume.tsx; then
  git commit -m "test: verify experience order is rendered consistently"
else
  git add web/src/components/resume/ResumeDesktopLayout.tsx web/src/components/resume/ResumeMobileLayout.tsx web/src/components/resume/PrintResume.tsx
  git commit -m "fix: render experiences in saved order"
fi
```

## Chunk 3: Final Verification

### Task 5: Run full verification and manual checks

**Files:**
- Modify only if regressions are found during verification: `web/src/...`
- Reference: `docs/superpowers/specs/2026-04-07-experience-ordering-design.md`

- [ ] **Step 1: Run the full frontend test suite**

Run: `npm.cmd run test`
Expected: PASS with no failing tests or unhandled errors.

- [ ] **Step 2: Run the frontend production build**

Run: `npm.cmd run build`
Expected: PASS with no TypeScript or Vite build errors.

- [ ] **Step 3: Perform manual verification in the editor**

Manual checks:
- open `/edit`
- in `工作经历`, click `下移` on the first item and confirm visual order changes immediately
- in `教育经历`, confirm the first item `上移` is disabled
- in `项目经历`, confirm the last item `下移` is disabled
- save the main resume and verify `/` follows the new order
- save a draft and verify `/drafts/:id` follows the same order
- verify `/print` follows the same order

- [ ] **Step 4: Commit any final verification fixes**

```bash
git add web/src
git commit -m "test: verify experience ordering workflow"
```

## Notes For The Implementer

- Follow @test-driven-development strictly: every move-control behavior starts with a failing test.
- Keep reorder logic local and pure. Do not add sort keys or backend schema changes.
- Prefer one reusable helper over three slightly different reorder implementations.
- Resist adding drag-and-drop abstractions or generic sortable systems; the approved design is explicitly `上移 / 下移` only.
