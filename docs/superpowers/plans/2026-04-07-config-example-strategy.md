# Config Example Strategy Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Stop tracking `config.json`, add a committed `config.example.json`, keep runtime config loading from the project-root `config.json`, and update deployment files/docs to match this workflow.

**Architecture:** Keep the runtime config path unchanged so the app still reads `./config.json`, but move the repo’s committed default into `config.example.json` and ignore `config.json`. Leave `AUTH_KEY` environment-only, and simplify deployment docs so operators copy the example file into place before starting the container.

**Tech Stack:** Go, Docker Compose, Markdown docs, Git ignore rules

---

## File Structure

### Config files

- Create: `config.example.json`
  Responsibility: Provide the committed, non-sensitive runtime config template users copy before starting the app.
- Modify: `config.json`
  Responsibility: Convert the currently tracked local config into an ignored local file, then remove it from version control.
- Modify: `.gitignore`
  Responsibility: Ignore `config.json` while continuing to track `config.example.json`.
- Modify: `internal/config/config.go`
  Responsibility: Keep reading the project-root `config.json`; likely no functional change, but verify the file path assumption still matches the new workflow.
- Modify: `internal/config/config_test.go`
  Responsibility: Ensure tests still cover project-root `config.json` loading with `AUTH_KEY` from env.

### Deployment files

- Modify: `deploy/docker-compose.yml`
  Responsibility: Keep mounting `../config.json` and `../data` directly, aligned with the root-level config workflow.
- Modify: `deploy/env/app.env.example`
  Responsibility: Keep only env-driven values (`AUTH_KEY`, optional bind/healthcheck) and ensure docs match.
- Delete: `deploy/config/config.production.json.example`
  Responsibility: Remove the now-redundant second config template to avoid two competing workflows.

### Documentation files

- Modify: `docs/deploy/manual-docker-nginx.md`
  Responsibility: Document `cp config.example.json config.json` as the required setup step and remove references to the deleted deployment config template.

## Chunk 1: Repo Config Strategy

### Task 1: Introduce `config.example.json` and stop tracking `config.json`

**Files:**
- Create: `config.example.json`
- Modify: `.gitignore`
- Modify: `config.json`

- [ ] **Step 1: Write the failing repository-state check**

Verify the current undesired state before changing files.

Run: `git ls-files config.json config.example.json`
Expected before changes:
- `config.json` is tracked
- `config.example.json` does not exist

- [ ] **Step 2: Create `config.example.json`**

Copy the current non-secret runtime shape into a committed example file.

Target content shape:

```json
{
  "browserPath": "/usr/bin/chromium",
  "port": "8088",
  "dbPath": "./data/resume.db",
  "frontendOrigin": "https://wenemoji.com"
}
```

For repo defaults, choose values that are safe and clearly example-oriented. If a Linux-oriented example is too deployment-specific for local development, document that in the file comments or deployment docs rather than adding a second example.

- [ ] **Step 3: Ignore `config.json`**

Update `.gitignore` to include:

```gitignore
config.json
```

Do not ignore `config.example.json`.

- [ ] **Step 4: Remove `config.json` from version control without deleting the local file**

Run:

```bash
git rm --cached config.json
```

Then restore or recreate the local working copy from `config.example.json` if needed for local use.

- [ ] **Step 5: Re-run the repository-state check**

Run: `git ls-files config.json config.example.json`
Expected after changes:
- `config.example.json` is tracked
- `config.json` is not tracked

- [ ] **Step 6: Commit the repo config strategy slice**

```bash
git add .gitignore config.example.json
git rm --cached config.json
git commit -m "chore: track config example instead of local config"
```

## Chunk 2: Runtime And Test Verification

### Task 2: Verify runtime still reads the root `config.json`

**Files:**
- Modify if needed: `internal/config/config.go`
- Modify: `internal/config/config_test.go`

- [ ] **Step 1: Write or refine tests that prove the runtime contract**

Ensure `internal/config/config_test.go` explicitly verifies:
- config is loaded from the working directory’s `config.json`
- `AUTH_KEY` must come from the environment
- `config.json` itself does not need `authKey`

If existing tests already cover this after the env-only change, tighten names/assertions to reflect the new repo workflow instead of adding duplicate tests.

- [ ] **Step 2: Run the focused config tests**

Run: `go test ./internal/config -count=1`
Expected: PASS.

- [ ] **Step 3: Adjust `internal/config/config.go` only if the tests reveal ambiguity**

Likely outcome: no code change required. Only change it if the tests show the loader no longer matches the intended root `config.json` behavior.

- [ ] **Step 4: Re-run the focused config tests**

Run: `go test ./internal/config -count=1`
Expected: PASS.

- [ ] **Step 5: Run the full Go suite**

Run: `go test ./... -count=1`
Expected: PASS.

- [ ] **Step 6: Commit any backend/test adjustments**

```bash
git add internal/config/config.go internal/config/config_test.go
git commit -m "test: verify root config file workflow"
```

## Chunk 3: Deployment Asset Cleanup

### Task 3: Align Docker deployment with the root `config.json` workflow

**Files:**
- Modify: `deploy/docker-compose.yml`
- Modify: `deploy/env/app.env.example`
- Delete: `deploy/config/config.production.json.example`
- Modify: `docs/deploy/manual-docker-nginx.md`

- [ ] **Step 1: Write the failing/static deployment check**

Run and inspect:

```bash
docker compose --env-file deploy/env/app.env.example -f deploy/docker-compose.yml config
```

Expected before cleanup:
- Compose still works, but docs/templates still imply multiple config entrypoints

The failure here is conceptual rather than syntactic; document what must be removed or aligned before editing.

- [ ] **Step 2: Keep Compose pointed at the root config file**

`deploy/docker-compose.yml` should continue using:

```yaml
volumes:
  - ../config.json:/app/config.json:ro
  - ../data:/app/data
```

Do not reintroduce `APP_CONFIG_PATH` or extra config path variables.

- [ ] **Step 3: Keep env example minimal**

`deploy/env/app.env.example` should document only:
- `AUTH_KEY`
- `APP_BIND`
- `APP_HEALTHCHECK_URL`

- [ ] **Step 4: Delete the redundant production config template**

Remove `deploy/config/config.production.json.example` so there is only one repo-tracked config template source.

- [ ] **Step 5: Rewrite the manual deployment guide**

Update `docs/deploy/manual-docker-nginx.md` to instruct:

```bash
cp config.example.json config.json
```

Then edit `config.json`, set `AUTH_KEY`, and run Docker Compose. Remove any instructions that still reference `deploy/config/config.production.json.example`.

- [ ] **Step 6: Re-run Compose static validation**

Run: `docker compose --env-file deploy/env/app.env.example -f deploy/docker-compose.yml config`
Expected: PASS.

- [ ] **Step 7: Commit the deployment cleanup slice**

```bash
git add deploy/docker-compose.yml deploy/env/app.env.example docs/deploy/manual-docker-nginx.md
git add -A deploy/config/config.production.json.example
git commit -m "docs: align deploy flow with config example"
```

## Chunk 4: Final Verification

### Task 4: Verify the full workflow end to end

**Files:**
- Modify only if verification reveals issues in the files above

- [ ] **Step 1: Run backend tests**

Run: `go test ./... -count=1`
Expected: PASS.

- [ ] **Step 2: Run frontend tests**

Run: `cd web && npm run test`
Expected: PASS.

- [ ] **Step 3: Run frontend build**

Run: `cd web && npm run build`
Expected: PASS.

- [ ] **Step 4: Validate Compose config one more time**

Run: `docker compose --env-file deploy/env/app.env.example -f deploy/docker-compose.yml config`
Expected: PASS.

- [ ] **Step 5: Manual workflow checklist**

Confirm the documented flow is coherent:
- clone repo
- `cp config.example.json config.json`
- edit `config.json`
- set `AUTH_KEY`
- start app with Docker Compose

- [ ] **Step 6: Commit any final verification fixes**

```bash
git add .gitignore config.example.json internal/config deploy docs/deploy
git commit -m "test: verify config example workflow"
```

## Notes For The Implementer

- Keep `AUTH_KEY` env-only. Do not regress back to file-based auth secrets.
- Do not introduce a second config template after deleting `deploy/config/config.production.json.example`.
- Preserve the runtime assumption that the app reads `./config.json` from the project root.
- Be careful not to delete the local `config.json` content from the developer’s working copy when untracking it from git.
