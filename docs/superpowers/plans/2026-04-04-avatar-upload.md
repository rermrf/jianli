# Avatar Upload Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add avatar upload, crop, persistence, and rendering across the edit page, public resume page, and print page.

**Architecture:** Keep avatar upload as a dedicated authenticated backend endpoint that stores the final cropped image under `data/uploads/avatars/` and returns a stable URL. The frontend handles image selection, local preview, square crop output, upload, and then saves the returned `avatarUrl` into the existing resume payload so all pages consume the same field.

**Tech Stack:** Go 1.26, Gin, GORM, SQLite, React, TypeScript, React Router, Tailwind CSS, Vitest

---

## File Structure

- Modify: `internal/model/resume.go`
  - Add `profile.avatarUrl` support to persisted resume data shape if needed by helpers.
- Create: `internal/handler/upload.go`
  - Authenticated avatar upload endpoint.
- Create: `internal/handler/upload_test.go`
  - Upload auth/file validation/success tests.
- Modify: `cmd/server/main.go`
  - Register upload route and static `/uploads/*` file serving.
- Create: `web/src/lib/upload.ts`
  - Avatar upload client.
- Create: `web/src/components/editor/AvatarUploader.tsx`
  - File picker, preview, crop UI, confirm/cancel flow.
- Modify: `web/src/pages/EditPage.tsx`
  - Wire avatar uploader into resume editing flow.
- Modify: `web/src/components/resume/ProfileCard.tsx`
  - Show real avatar when `avatarUrl` exists.
- Modify: `web/src/components/resume/PrintResume.tsx`
  - Show real avatar in print layout.
- Modify: `web/src/types/resume.ts`
  - Add `avatarUrl` to profile type.
- Modify: `web/src/data/mockResume.ts`
  - Seed empty `avatarUrl`.
- Create or modify tests:
  - `web/src/pages/__tests__/edit-page.test.tsx`
  - `web/src/pages/__tests__/routing.test.tsx`
  - avatar-uploader focused test if needed

## Chunk 1: Backend Upload Endpoint

### Task 1: Add failing backend upload tests

**Files:**
- Create: `internal/handler/upload_test.go`

- [ ] **Step 1: Write failing tests for upload auth, invalid file type, and success response**
- [ ] **Step 2: Run `go test ./internal/handler -run Upload` and confirm failure**
- [ ] **Step 3: Implement minimal upload handler behavior contract**
- [ ] **Step 4: Re-run `go test ./internal/handler -run Upload` and confirm pass**
- [ ] **Step 5: Commit**

### Task 2: Wire upload route and static file serving

**Files:**
- Modify: `cmd/server/main.go`
- Modify: `internal/store/sqlite.go` only if upload directory bootstrap helper belongs there

- [ ] **Step 1: Write a failing server smoke test for `/api/upload/avatar` and `/uploads/avatars/...` route registration**
- [ ] **Step 2: Run `go test ./cmd/server` and confirm failure**
- [ ] **Step 3: Register authenticated upload route and static serving path**
- [ ] **Step 4: Re-run `go test ./cmd/server` and confirm pass**
- [ ] **Step 5: Commit**

## Chunk 2: Frontend Upload Flow

### Task 3: Add failing frontend upload client test

**Files:**
- Create: `web/src/lib/upload.ts`
- Add tests in a new or existing `web/src/lib/__tests__/upload.test.ts`

- [ ] **Step 1: Write failing tests for multipart avatar upload success and error handling**
- [ ] **Step 2: Run targeted Vitest tests and confirm failure**
- [ ] **Step 3: Implement `uploadAvatar(file)` API helper**
- [ ] **Step 4: Re-run targeted tests and confirm pass**
- [ ] **Step 5: Commit**

### Task 4: Add failing edit-page avatar interaction test

**Files:**
- Modify: `web/src/pages/__tests__/edit-page.test.tsx`
- Create if needed: `web/src/components/editor/AvatarUploader.tsx`

- [ ] **Step 1: Write failing test for selecting an image, confirming crop/upload, and persisting returned `avatarUrl`**
- [ ] **Step 2: Run the edit-page test and confirm failure**
- [ ] **Step 3: Implement avatar uploader UI with square crop flow and hook it into `EditPage`**
- [ ] **Step 4: Re-run the edit-page test and confirm pass**
- [ ] **Step 5: Commit**

## Chunk 3: Resume And Print Rendering

### Task 5: Add failing rendering tests for avatar display

**Files:**
- Modify: `web/src/pages/__tests__/routing.test.tsx`
- Modify: `web/src/types/resume.ts`
- Modify: `web/src/data/mockResume.ts`
- Modify: `web/src/components/resume/ProfileCard.tsx`
- Modify: `web/src/components/resume/PrintResume.tsx`

- [ ] **Step 1: Write failing tests for avatar rendering on public resume and print page when `avatarUrl` exists**
- [ ] **Step 2: Run routing/print-related tests and confirm failure**
- [ ] **Step 3: Implement avatar rendering with fallback to placeholder**
- [ ] **Step 4: Re-run the tests and confirm pass**
- [ ] **Step 5: Commit**

## Chunk 4: Final Verification

### Task 6: Verify integrated avatar flow

**Files:**
- Modify docs only if needed

- [ ] **Step 1: Run `go test ./...`**
- [ ] **Step 2: Run `npm run test` in `web/`**
- [ ] **Step 3: Run `npm run build` and `npm run lint` in `web/`**
- [ ] **Step 4: Start the backend and frontend locally, upload an avatar, save resume, and verify avatar on `/`, `/print`, and direct `/uploads/...` access**
- [ ] **Step 5: Commit final cleanup**

## Notes

- Keep the first version simple: square crop only, no rotation/filter pipeline.
- Do not add a new database table for avatars.
- Store only the final avatar URL in resume data.
- If a crop library feels too heavy, prefer a small canvas-based implementation over overbuilding.
