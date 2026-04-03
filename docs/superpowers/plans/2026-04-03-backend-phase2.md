# Backend Phase 2 Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the Go backend and real frontend data integration for the online resume system, including SQLite persistence, Key auth, visitor tracking, and PDF export.

**Architecture:** Keep the backend as a single Gin service with focused internal packages for config, model, store, middleware, handlers, and PDF export. Reuse the existing phase-one frontend pages, replacing mock data/storage with API clients and authenticated route flows while preserving the approved UI structure.

**Tech Stack:** Go 1.25, Gin, SQLite, chromedp, React, TypeScript, React Router, Tailwind CSS, Vitest

---

## File Structure

- Create: `cmd/server/main.go`
- Create: `internal/config/config.go`
- Create: `internal/model/resume.go`
- Create: `internal/model/visitor.go`
- Create: `internal/store/sqlite.go`
- Create: `internal/store/resume_store.go`
- Create: `internal/store/visitor_store.go`
- Create: `internal/middleware/auth.go`
- Create: `internal/middleware/cors.go`
- Create: `internal/handler/auth.go`
- Create: `internal/handler/resume.go`
- Create: `internal/handler/visitor.go`
- Create: `internal/pdf/export.go`
- Create: `internal/httpapi/response.go`
- Create: `internal/seed/resume_seed.go`
- Create: `web/src/lib/api.ts`
- Create: `web/src/lib/visitors.ts`
- Modify: `web/src/lib/auth.ts`
- Modify: `web/src/lib/storage.ts`
- Modify: `web/src/hooks/useResumeDraft.ts`
- Modify: `web/src/pages/LoginPage.tsx`
- Modify: `web/src/pages/EditPage.tsx`
- Modify: `web/src/pages/ResumePage.tsx`
- Modify: `web/src/pages/VisitorsPage.tsx`
- Modify: `web/src/pages/PrintPage.tsx`

## Chunk 1: Backend Skeleton And Storage

### Task 1: Add backend dependencies and config skeleton
- [ ] Write the failing config test
- [ ] Run `go test ./internal/config/...` to watch it fail
- [ ] Add minimal config loader and server boot skeleton
- [ ] Re-run config test
- [ ] Commit

### Task 2: Add SQLite bootstrap and schema creation
- [ ] Write failing schema/bootstrap tests
- [ ] Run the store bootstrap tests and confirm failure
- [ ] Implement DB open + schema initialization
- [ ] Re-run the store bootstrap tests
- [ ] Commit

### Task 3: Seed and persist the single resume record
- [ ] Write failing tests for seed-on-empty and get/update resume
- [ ] Run `go test ./internal/store -run Resume` and confirm failure
- [ ] Implement the resume store
- [ ] Re-run the resume store tests
- [ ] Commit

## Chunk 2: Auth And Visitor Data

### Task 4: Add auth middleware and verify endpoint
- [ ] Write failing tests for valid/invalid auth key flows
- [ ] Run handler tests and confirm failure
- [ ] Implement middleware and verify handler
- [ ] Re-run auth tests
- [ ] Commit

### Task 5: Implement visitor persistence, dedupe, and stats queries
- [ ] Write failing tests for insert dedupe, stats aggregation, and range listing
- [ ] Run visitor store tests and confirm failure
- [ ] Implement visitor queries
- [ ] Re-run visitor store tests
- [ ] Commit

### Task 6: Implement visitor HTTP handlers
- [ ] Write failing handler tests for list/stats/create/update
- [ ] Run handler tests and confirm failure
- [ ] Implement visitor handlers with common response helpers
- [ ] Re-run visitor handler tests
- [ ] Commit

## Chunk 3: Resume HTTP And PDF

### Task 7: Implement resume HTTP handlers
- [ ] Write failing tests for public GET and protected PUT
- [ ] Run resume handler tests and confirm failure
- [ ] Implement resume handlers
- [ ] Re-run resume handler tests
- [ ] Commit

### Task 8: Add PDF export
- [ ] Write failing export tests for PDF response behavior
- [ ] Run PDF/export tests and confirm failure
- [ ] Implement HTML-to-PDF export using `/print`
- [ ] Re-run export tests
- [ ] Commit

### Task 9: Wire the full server entrypoint
- [ ] Write a failing smoke test for route registration
- [ ] Run server smoke tests and confirm failure
- [ ] Wire stores, handlers, middleware, static frontend, and API groups
- [ ] Re-run backend test suite
- [ ] Commit

## Chunk 4: Frontend API Integration

### Task 10: Replace fake auth and local draft helpers with API clients
- [ ] Write failing frontend tests for login verification and API-backed resume load/save
- [ ] Run targeted Vitest tests and confirm failure
- [ ] Implement the API client layer
- [ ] Re-run frontend helper tests
- [ ] Commit

### Task 11: Connect pages to real backend data
- [ ] Write failing page tests for authenticated save, public load, and visitors fetch
- [ ] Run page tests and confirm failure
- [ ] Implement real data loading/saving and auth redirect handling
- [ ] Re-run page tests
- [ ] Commit

### Task 12: Add real visitor tracking from the public resume page
- [ ] Write failing tests for track-on-open and beacon-on-leave behavior
- [ ] Run targeted tests and confirm failure
- [ ] Implement visitor create/update reporting
- [ ] Re-run tracking tests
- [ ] Commit

## Chunk 5: Final Verification

### Task 13: Verify integrated backend + frontend
- [ ] Run `go test ./...`
- [ ] Run `npm run test` in `web/`
- [ ] Run `npm run build` in `web/`
- [ ] Start the Go server and manually verify `/`, `/edit`, `/visitors`, `/login`, `/print`, `/api/resume`, `/api/visitors`, `/api/resume/pdf`
- [ ] Commit any final cleanup

## Notes

- Keep the backend single-binary and focused; no extra framework layers.
- Do not reintroduce localStorage-based fake persistence once API integration starts.
- Preserve the current frontend structure unless the integration reveals a real boundary problem.
- If chromedp proves flaky in test mode, keep handler tests narrow and verify export manually in the final verification step.
