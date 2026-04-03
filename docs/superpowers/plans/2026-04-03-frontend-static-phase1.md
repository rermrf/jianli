# Frontend Static Phase 1 Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the phase-one React frontend prototype for the online resume system with 5 static routes, responsive mobile/PC layouts, mock-data-driven interactions, local auth simulation, and a printable `/print` page.

**Architecture:** Keep the Vite app thin at the route layer and push reusable display/editor/visitor blocks into focused component folders. Use a shared resume draft source backed by `localStorage`, a tiny auth/session utility backed by `sessionStorage`, and route guards in React Router so the static prototype already mirrors the later API-backed flows.

**Tech Stack:** Vite, React 19, TypeScript, React Router, Tailwind CSS, Vitest, React Testing Library

---

## File Structure

### Planned file responsibilities

- Modify: `web/package.json`
  - Add runtime and test dependencies plus useful scripts.
- Modify: `web/vite.config.ts`
  - Register React and Tailwind plugins.
- Modify: `web/src/main.tsx`
  - Boot the router-enabled app and test-friendly global CSS entry.
- Modify: `web/src/index.css`
  - Replace template CSS with Tailwind imports and global tokens.
- Create: `web/src/app/router.tsx`
  - Define the 5 routes and protected-route wrapper.
- Create: `web/src/app/providers.tsx`
  - Optional central wrapper for shared providers if needed.
- Create: `web/src/lib/auth.ts`
  - Read/write simulated auth state and redirect targets in `sessionStorage`.
- Create: `web/src/lib/storage.ts`
  - Read/write the resume draft in `localStorage`, with safe fallback behavior.
- Create: `web/src/lib/format.ts`
  - Centralize small presentation helpers such as masked IP and duration formatting.
- Create: `web/src/lib/__tests__/auth.test.ts`
  - Cover auth storage and redirect-target behavior.
- Create: `web/src/lib/__tests__/storage.test.ts`
  - Cover resume draft persistence and bad-JSON fallback.
- Create: `web/src/lib/__tests__/format.test.ts`
  - Cover visitor formatting/filter helpers.
- Create: `web/src/types/resume.ts`
  - Define shared resume data types.
- Create: `web/src/types/visitors.ts`
  - Define visitor stats/list types.
- Create: `web/src/data/mockResume.ts`
  - Provide canonical phase-one resume seed data.
- Create: `web/src/data/mockVisitors.ts`
  - Provide stats, trend points, and visitor-list seed data.
- Create: `web/src/hooks/useResumeDraft.ts`
  - Manage editable resume state and save/reset behavior.
- Create: `web/src/components/layout/AppShell.tsx`
  - Shared app container and responsive background handling.
- Create: `web/src/components/layout/TopNav.tsx`
  - Desktop navigation for resume/edit/visitors/export.
- Create: `web/src/components/layout/ProtectedRoute.tsx`
  - Redirect unauthenticated access to `/login`.
- Create: `web/src/components/common/SectionCard.tsx`
  - Reusable card wrapper for page blocks.
- Create: `web/src/components/common/Tag.tsx`
  - Reusable skill tag chip.
- Create: `web/src/components/common/Button.tsx`
  - Shared button variants.
- Create: `web/src/components/common/EmptyHint.tsx`
  - Fallback UI for empty editable sections where useful.
- Create: `web/src/components/resume/ProfileCard.tsx`
  - Personal information summary block.
- Create: `web/src/components/resume/SkillSection.tsx`
  - Skill tags block.
- Create: `web/src/components/resume/TimelineSection.tsx`
  - Shared work/project/education/awards listing block.
- Create: `web/src/components/resume/ResumeDesktopLayout.tsx`
  - Desktop two-column composition.
- Create: `web/src/components/resume/ResumeMobileLayout.tsx`
  - Mobile stacked composition.
- Create: `web/src/components/resume/PrintResume.tsx`
  - A4-friendly printable resume composition.
- Create: `web/src/components/editor/FieldInput.tsx`
  - Shared text input and textarea styling.
- Create: `web/src/components/editor/EditableTagList.tsx`
  - Add/remove skills.
- Create: `web/src/components/editor/EditableListSection.tsx`
  - Generic add/remove list block for work/projects/education/awards.
- Create: `web/src/components/editor/SaveToast.tsx`
  - Lightweight save-success feedback.
- Create: `web/src/components/visitors/StatsCards.tsx`
  - Mobile/desktop stat-card rendering.
- Create: `web/src/components/visitors/TrendChart.tsx`
  - Lightweight SVG or CSS line chart for mock trends.
- Create: `web/src/components/visitors/VisitorList.tsx`
  - Mobile visitor cards.
- Create: `web/src/components/visitors/VisitorTable.tsx`
  - Desktop visitor table.
- Create: `web/src/pages/ResumePage.tsx`
  - Public resume route.
- Create: `web/src/pages/EditPage.tsx`
  - Protected editable resume route.
- Create: `web/src/pages/VisitorsPage.tsx`
  - Protected visitor analytics route.
- Create: `web/src/pages/LoginPage.tsx`
  - Login route with simulated auth flow.
- Create: `web/src/pages/PrintPage.tsx`
  - Printable resume route.
- Create: `web/src/test/setup.ts`
  - Testing Library and Jest DOM setup.
- Create: `web/src/pages/__tests__/routing.test.tsx`
  - Cover protected routing and login redirect flow.
- Create: `web/src/pages/__tests__/edit-page.test.tsx`
  - Cover save behavior and draft persistence.
- Create: `web/src/pages/__tests__/visitors-page.test.tsx`
  - Cover time-range toggling.

## Chunk 1: Tooling And App Foundation

### Task 1: Install routing, Tailwind, and test tooling

**Files:**
- Modify: `web/package.json`
- Modify: `web/package-lock.json`

- [ ] **Step 1: Add the failing dependency expectations to the plan checklist**

Record that the app cannot yet support router, Tailwind, or tests because `web/package.json` currently only includes React/Vite template dependencies.

- [ ] **Step 2: Install runtime dependencies**

Run:

```bash
npm install react-router-dom
```

Expected: `react-router-dom` appears under `dependencies`.

- [ ] **Step 3: Install styling and test dependencies**

Run:

```bash
npm install -D tailwindcss @tailwindcss/vite vitest jsdom @testing-library/react @testing-library/jest-dom @testing-library/user-event
```

Expected: all packages appear under `devDependencies`.

- [ ] **Step 4: Update scripts in `web/package.json`**

Add at least:

```json
"test": "vitest run",
"test:watch": "vitest"
```

Keep existing `dev`, `build`, `lint`, and `preview`.

- [ ] **Step 5: Commit**

```bash
git add web/package.json web/package-lock.json
git commit -m "chore: add frontend routing styling and test deps"
```

### Task 2: Replace the Vite starter shell with the project app shell

**Files:**
- Modify: `web/vite.config.ts`
- Modify: `web/src/main.tsx`
- Modify: `web/src/index.css`
- Delete/Stop using: `web/src/App.css`
- Delete/Stop using: `web/src/App.tsx`
- Delete/Stop using: `web/src/assets/react.svg`
- Delete/Stop using: `web/src/assets/vite.svg`
- Delete/Stop using: `web/src/assets/hero.png`
- Create: `web/src/app/router.tsx`
- Create: `web/src/app/providers.tsx`
- Test: `web/src/pages/__tests__/routing.test.tsx`

- [ ] **Step 1: Write the first failing routing smoke test**

Add a test like:

```tsx
it('renders the public resume route at /', () => {
  renderWithRouter('/')
  expect(screen.getByText('温庆京')).toBeInTheDocument()
})
```

Use a placeholder expectation that will fail until the real router and page exist.

- [ ] **Step 2: Run the single test to verify it fails**

Run:

```bash
npm run test -- src/pages/__tests__/routing.test.tsx
```

Expected: FAIL because router/page modules do not exist yet.

- [ ] **Step 3: Wire Tailwind and the application entry**

Implement:

- `vite.config.ts` with `tailwindcss()` plugin from `@tailwindcss/vite`
- `src/index.css` with Tailwind import plus app-level CSS variables for colors, shadows, spacing, and print helpers
- `src/main.tsx` rendering a router-based app instead of the Vite starter component
- `src/app/router.tsx` exporting a temporary minimal router with all 5 route placeholders

- [ ] **Step 4: Re-run the smoke test and dev build**

Run:

```bash
npm run test -- src/pages/__tests__/routing.test.tsx
npm run build
```

Expected: the test may still fail on missing real page content, but Vite/TypeScript should now compile the routing shell. Fix compile issues before continuing.

- [ ] **Step 5: Commit**

```bash
git add web/vite.config.ts web/src/main.tsx web/src/index.css web/src/app web/src/pages/__tests__/routing.test.tsx
git commit -m "feat: bootstrap app shell with router and tailwind"
```

## Chunk 2: Shared Data, Types, And Persistence

### Task 3: Define mock data contracts and persistence helpers

**Files:**
- Create: `web/src/types/resume.ts`
- Create: `web/src/types/visitors.ts`
- Create: `web/src/data/mockResume.ts`
- Create: `web/src/data/mockVisitors.ts`
- Create: `web/src/lib/storage.ts`
- Create: `web/src/lib/format.ts`
- Create: `web/src/lib/__tests__/storage.test.ts`
- Create: `web/src/lib/__tests__/format.test.ts`

- [ ] **Step 1: Write failing tests for storage fallback and visitor formatting**

Add tests covering:

```ts
it('returns the seed resume when localStorage contains invalid JSON', () => {})
it('masks visitor IP addresses for list display', () => {})
it('filters visitor trend data by selected range', () => {})
```

- [ ] **Step 2: Run only those tests to verify they fail**

Run:

```bash
npm run test -- src/lib/__tests__/storage.test.ts src/lib/__tests__/format.test.ts
```

Expected: FAIL because types, mocks, and helpers do not exist yet.

- [ ] **Step 3: Implement the minimal shared data layer**

Create focused modules for:

- resume types and seed data mirroring the approved phase-one content
- visitor types and seed data for stats, trend points, and list/table rows
- `loadResumeDraft`, `saveResumeDraft`, `resetResumeDraft`
- `maskIp`, `formatDuration`, `getVisitorRangeData`

Keep functions pure where possible.

- [ ] **Step 4: Re-run the targeted tests**

Run:

```bash
npm run test -- src/lib/__tests__/storage.test.ts src/lib/__tests__/format.test.ts
```

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add web/src/types web/src/data web/src/lib
git commit -m "feat: add mock data and persistence helpers"
```

### Task 4: Implement auth/session helpers and draft hook

**Files:**
- Create: `web/src/lib/auth.ts`
- Create: `web/src/hooks/useResumeDraft.ts`
- Create: `web/src/lib/__tests__/auth.test.ts`

- [ ] **Step 1: Write failing tests for auth storage and redirect target handling**

Add tests covering:

```ts
it('marks the session as authenticated after login')
it('stores and consumes a redirect target once')
```

- [ ] **Step 2: Run the auth tests to verify they fail**

Run:

```bash
npm run test -- src/lib/__tests__/auth.test.ts
```

Expected: FAIL because the auth helper does not exist yet.

- [ ] **Step 3: Implement the auth helper and resume draft hook**

Implement small APIs such as:

```ts
isAuthenticated()
loginWithKey(key: string)
logout()
setRedirectPath(path: string)
consumeRedirectPath()
```

And a `useResumeDraft` hook that exposes current draft, update helpers, save, and reset behavior on top of `storage.ts`.

- [ ] **Step 4: Re-run auth tests**

Run:

```bash
npm run test -- src/lib/__tests__/auth.test.ts
```

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add web/src/lib/auth.ts web/src/hooks/useResumeDraft.ts web/src/lib/__tests__/auth.test.ts
git commit -m "feat: add local auth and draft state helpers"
```

## Chunk 3: Shared UI Composition

### Task 5: Build common layout and protected-route primitives

**Files:**
- Create: `web/src/components/layout/AppShell.tsx`
- Create: `web/src/components/layout/TopNav.tsx`
- Create: `web/src/components/layout/ProtectedRoute.tsx`
- Create: `web/src/components/common/Button.tsx`
- Create: `web/src/components/common/SectionCard.tsx`
- Create: `web/src/components/common/Tag.tsx`
- Create: `web/src/components/common/EmptyHint.tsx`
- Modify: `web/src/app/router.tsx`
- Test: `web/src/pages/__tests__/routing.test.tsx`

- [ ] **Step 1: Extend the routing test to cover protection**

Add a test like:

```tsx
it('redirects unauthenticated users from /edit to /login', async () => {
  renderWithRouter('/edit')
  expect(await screen.findByText('管理后台')).toBeInTheDocument()
})
```

- [ ] **Step 2: Run the routing test to verify it fails**

Run:

```bash
npm run test -- src/pages/__tests__/routing.test.tsx
```

Expected: FAIL because protection/layout logic is not implemented yet.

- [ ] **Step 3: Implement shared layout primitives**

Implement:

- `AppShell` for page background, max width, and mobile/desktop section spacing
- `TopNav` for desktop-only public/protected navigation
- `ProtectedRoute` using `auth.ts` redirect helpers
- shared buttons/cards/tags so pages do not duplicate low-level Tailwind classes

- [ ] **Step 4: Re-run the routing test**

Run:

```bash
npm run test -- src/pages/__tests__/routing.test.tsx
```

Expected: PASS for redirect behavior and basic public route rendering.

- [ ] **Step 5: Commit**

```bash
git add web/src/components/common web/src/components/layout web/src/app/router.tsx web/src/pages/__tests__/routing.test.tsx
git commit -m "feat: add shared layout primitives and route protection"
```

### Task 6: Build reusable resume and editor section components

**Files:**
- Create: `web/src/components/resume/ProfileCard.tsx`
- Create: `web/src/components/resume/SkillSection.tsx`
- Create: `web/src/components/resume/TimelineSection.tsx`
- Create: `web/src/components/resume/ResumeDesktopLayout.tsx`
- Create: `web/src/components/resume/ResumeMobileLayout.tsx`
- Create: `web/src/components/resume/PrintResume.tsx`
- Create: `web/src/components/editor/FieldInput.tsx`
- Create: `web/src/components/editor/EditableTagList.tsx`
- Create: `web/src/components/editor/EditableListSection.tsx`
- Create: `web/src/components/editor/SaveToast.tsx`

- [ ] **Step 1: Write a focused component test for editable tags**

Add a test like:

```tsx
it('adds a new skill tag and removes an existing one')
```

Place it in `web/src/pages/__tests__/edit-page.test.tsx` or a nearby component test if you prefer page-driven coverage.

- [ ] **Step 2: Run the targeted test to verify it fails**

Run:

```bash
npm run test -- src/pages/__tests__/edit-page.test.tsx
```

Expected: FAIL because the editable components do not exist yet.

- [ ] **Step 3: Implement the reusable content blocks**

Keep components focused:

- resume-side components only read structured data
- editor-side components only render and emit changes
- `PrintResume` renders a simplified A4-friendly version without app nav chrome

- [ ] **Step 4: Re-run the edit-page test**

Run:

```bash
npm run test -- src/pages/__tests__/edit-page.test.tsx
```

Expected: still partially failing until the page container is wired, but the component-level compile path should be clear.

- [ ] **Step 5: Commit**

```bash
git add web/src/components/resume web/src/components/editor web/src/pages/__tests__/edit-page.test.tsx
git commit -m "feat: add reusable resume and editor sections"
```

## Chunk 4: Route Pages

### Task 7: Implement `/`, `/print`, and the shared resume presentation flow

**Files:**
- Create: `web/src/pages/ResumePage.tsx`
- Create: `web/src/pages/PrintPage.tsx`
- Modify: `web/src/app/router.tsx`
- Test: `web/src/pages/__tests__/routing.test.tsx`

- [ ] **Step 1: Expand the routing test to assert print navigation**

Add a test like:

```tsx
it('navigates from resume export action to /print')
```

- [ ] **Step 2: Run the routing test to verify it fails**

Run:

```bash
npm run test -- src/pages/__tests__/routing.test.tsx
```

Expected: FAIL because the real pages and export link are not fully wired.

- [ ] **Step 3: Implement the resume pages**

Requirements:

- `ResumePage` chooses mobile vs desktop composition based on responsive CSS, not duplicate route logic
- reads the current draft from `useResumeDraft`/storage-backed loader
- desktop nav exposes edit, visitors, and print/export
- mobile page matches the Pencil information hierarchy
- `PrintPage` renders printable content only

- [ ] **Step 4: Re-run the routing test**

Run:

```bash
npm run test -- src/pages/__tests__/routing.test.tsx
```

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add web/src/pages/ResumePage.tsx web/src/pages/PrintPage.tsx web/src/app/router.tsx web/src/pages/__tests__/routing.test.tsx
git commit -m "feat: add resume and print routes"
```

### Task 8: Implement `/login` and `/edit`

**Files:**
- Create: `web/src/pages/LoginPage.tsx`
- Create: `web/src/pages/EditPage.tsx`
- Modify: `web/src/app/router.tsx`
- Modify: `web/src/pages/__tests__/edit-page.test.tsx`
- Modify: `web/src/pages/__tests__/routing.test.tsx`

- [ ] **Step 1: Add failing tests for login flow and edit persistence**

Cover:

```tsx
it('redirects back to the protected page after successful login')
it('saves edited resume data to localStorage')
```

- [ ] **Step 2: Run the login/edit tests to verify they fail**

Run:

```bash
npm run test -- src/pages/__tests__/routing.test.tsx src/pages/__tests__/edit-page.test.tsx
```

Expected: FAIL because page behavior is incomplete.

- [ ] **Step 3: Implement the login and edit pages**

Requirements:

- `LoginPage` validates non-empty input and calls `loginWithKey`
- `EditPage` uses `useResumeDraft`
- save action persists to `localStorage` and triggers `SaveToast`
- add/remove interactions work for skills and at least one repeated list section pattern that can be reused for all collections
- mobile and desktop forms both reflect the approved spec

- [ ] **Step 4: Re-run the login/edit tests**

Run:

```bash
npm run test -- src/pages/__tests__/routing.test.tsx src/pages/__tests__/edit-page.test.tsx
```

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add web/src/pages/LoginPage.tsx web/src/pages/EditPage.tsx web/src/app/router.tsx web/src/pages/__tests__/routing.test.tsx web/src/pages/__tests__/edit-page.test.tsx
git commit -m "feat: add login and edit flows"
```

### Task 9: Implement `/visitors`

**Files:**
- Create: `web/src/components/visitors/StatsCards.tsx`
- Create: `web/src/components/visitors/TrendChart.tsx`
- Create: `web/src/components/visitors/VisitorList.tsx`
- Create: `web/src/components/visitors/VisitorTable.tsx`
- Create: `web/src/pages/VisitorsPage.tsx`
- Create: `web/src/pages/__tests__/visitors-page.test.tsx`
- Modify: `web/src/app/router.tsx`

- [ ] **Step 1: Write the failing visitor range-switch test**

Add a test like:

```tsx
it('switches between 7-day and 30-day visitor data views')
```

- [ ] **Step 2: Run the visitors test to verify it fails**

Run:

```bash
npm run test -- src/pages/__tests__/visitors-page.test.tsx
```

Expected: FAIL because the visitors page/components do not exist yet.

- [ ] **Step 3: Implement the visitors page and its focused components**

Requirements:

- stat cards support 3-card mobile layout and 4-card desktop layout
- trend chart is lightweight and dependency-free
- mobile uses visitor cards
- desktop uses a table
- range toggle swaps mock dataset slices using `getVisitorRangeData`

- [ ] **Step 4: Re-run the visitors test**

Run:

```bash
npm run test -- src/pages/__tests__/visitors-page.test.tsx
```

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add web/src/components/visitors web/src/pages/VisitorsPage.tsx web/src/pages/__tests__/visitors-page.test.tsx web/src/app/router.tsx
git commit -m "feat: add visitor analytics page"
```

## Chunk 5: Final Polish And Verification

### Task 10: Remove template leftovers and finish responsive polish

**Files:**
- Delete: `web/src/App.tsx`
- Delete: `web/src/App.css`
- Delete: `web/src/assets/react.svg`
- Delete: `web/src/assets/vite.svg`
- Delete: `web/src/assets/hero.png`
- Modify: any page/component files that still need spacing, colors, print styles, or responsive fixes

- [ ] **Step 1: Search for starter-template leftovers**

Run:

```bash
rg "vite|reactLogo|heroImg|Get started|Count is" web/src web/public
```

Expected: identify any remaining template content.

- [ ] **Step 2: Remove or replace all starter leftovers**

Delete unused starter files and clean imports/usages.

- [ ] **Step 3: Run the full verification suite**

Run:

```bash
npm run test
npm run build
npm run lint
```

Expected: all commands pass.

- [ ] **Step 4: Manually verify responsive routes**

Run:

```bash
npm run dev
```

Check in the browser:

- `/` on mobile and desktop
- `/edit` redirect and post-login flow
- `/visitors` range switching on mobile and desktop
- `/login` empty-state validation
- `/print` printable layout without nav chrome

- [ ] **Step 5: Commit**

```bash
git add web
git commit -m "feat: complete phase-one static frontend prototype"
```

## Notes For Execution

- Prefer small files with one clear purpose; do not collapse all UI into route files.
- Keep mock data centralized; do not inline large data objects inside page components.
- Use semantic Tailwind utility groupings and shared components to avoid repeated class soup.
- Avoid introducing a heavy chart library for the static visitor trend; a lightweight SVG path is enough.
- Avoid overbuilding abstractions beyond the 5 approved routes and the current static prototype scope.
- If a page-specific component becomes too generic during implementation, stop and keep it local unless another page truly needs it.
- Preserve the approved phase-one scope: no real API clients, no backend integration, no true PDF export.

## Plan Review Notes

- This environment prohibits me from spawning a review subagent without explicit user authorization, so the normal plan-review loop is replaced here with a local self-review against the approved spec.
- Before execution, re-read:
  - `docs/superpowers/specs/2026-04-03-frontend-static-phase1-design.md`
  - `docs/superpowers/specs/2026-04-03-online-resume-design.md`

