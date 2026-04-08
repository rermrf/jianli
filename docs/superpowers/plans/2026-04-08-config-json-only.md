# Config JSON Only Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Move all runtime configuration back into `config.json`, stop using environment variables for `authKey`, and make `config.example.json` the only tracked template file.

**Architecture:** Keep runtime config loading from the project-root `config.json`, but restore `authKey` to file-based loading and remove the `AUTH_KEY` environment dependency from code and deployment assets. Replace the tracked root config file with `config.example.json`, ignore the real `config.json`, and simplify manual deployment so operators only copy the example file and run Docker Compose.

**Tech Stack:** Go, Docker Compose, Markdown docs, Git ignore rules, Git file tracking

---

## File Structure

### Config files

- Create: `config.example.json`
  Responsibility: Provide the tracked, complete runtime config template including `authKey`.
- Modify: `config.json`
  Responsibility: Keep a local working copy for development, but remove it from git tracking.
- Modify: `.gitignore`
  Responsibility: Ignore `config.json` while continuing to track `config.example.json`.

### Backend config files

- Modify: `internal/config/config.go`
  Responsibility: Read `authKey` from `config.json` again and remove the `AUTH_KEY` environment fallback.
- Modify: `internal/config/config_test.go`
  Responsibility: Rewrite tests to validate config-file-only `authKey` loading.

### Deployment files

- Modify: `deploy/docker-compose.yml`
  Responsibility: Remove `AUTH_KEY` from service environment and rely only on the mounted `config.json`.
- Delete: `deploy/env/app.env.example`
  Responsibility: Remove the now-unneeded deployment env template.
- Modify: `deploy/scripts/deploy.sh`
  Responsibility: Remove the `env-file` dependency and run Compose directly.
- Modify: `deploy/scripts/healthcheck.sh`
  Responsibility: Keep the healthcheck default aligned with the file-based deployment and port `8088`.
- Delete: `deploy/config/config.production.json.example`
  Responsibility: Remove the redundant second config template source.

### Documentation files

- Modify: `docs/deploy/manual-docker-nginx.md`
  Responsibility: Rewrite the deployment flow to use `cp config.example.json config.json` and no environment variables.

## Chunk 1: Backend Config Source Switch

### Task 1: Move `authKey` back into `config.json`

**Files:**
- Modify: `internal/config/config.go`
- Modify: `internal/config/config_test.go`

- [ ] **Step 1: Write the failing backend config tests**

Rewrite `internal/config/config_test.go` to assert:
- `authKey` is loaded from `config.json`
- missing `authKey` in `config.json` causes failure
- no environment variable setup is needed
- BOM-prefixed `config.json` with `authKey` still loads correctly

Example target test:

```go
func TestLoadUsesAuthKeyFromConfigFile(t *testing.T) {
	withTempConfigDir(t)

	if err := os.WriteFile("config.json", []byte(`{
  "authKey": "resume-key",
  "port": "9090"
}`), 0o644); err != nil {
		t.Fatalf("WriteFile(config.json) returned error: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if cfg.AuthKey != "resume-key" {
		t.Fatalf("expected auth key from config file, got %q", cfg.AuthKey)
	}
}
```

- [ ] **Step 2: Run the focused config tests to verify they fail**

Run: `go test ./internal/config -count=1`
Expected: FAIL because `Load()` still requires `AUTH_KEY` from environment.

- [ ] **Step 3: Implement file-based `authKey` loading**

Update `internal/config/config.go` so:
- `Config.AuthKey` is read directly from JSON again
- the `json:"-"` tag is removed
- `os.Getenv("AUTH_KEY")` logic is removed
- the error becomes `authKey is required in config.json`

Example target shape:

```go
type Config struct {
	AuthKey        string `json:"authKey"`
	BrowserPath    string `json:"browserPath"`
	DBPath         string `json:"dbPath"`
	FrontendOrigin string `json:"frontendOrigin"`
	Port           string `json:"port"`
}
```

- [ ] **Step 4: Re-run the focused config tests**

Run: `go test ./internal/config -count=1`
Expected: PASS.

- [ ] **Step 5: Run the full Go suite**

Run: `go test ./... -count=1`
Expected: PASS.

- [ ] **Step 6: Commit the backend config slice**

```bash
git add internal/config/config.go internal/config/config_test.go
git commit -m "refactor: read auth key from config file"
```

## Chunk 2: Config Template And Git Tracking

### Task 2: Replace tracked `config.json` with `config.example.json`

**Files:**
- Create: `config.example.json`
- Modify: `.gitignore`
- Remove from tracking: `config.json`

- [ ] **Step 1: Capture the current tracked state**

Run: `git ls-files config.json config.example.json`
Expected before changes:
- `config.json` is tracked
- `config.example.json` is absent

- [ ] **Step 2: Create `config.example.json`**

Populate it with the full config shape, including `authKey`.

Suggested content:

```json
{
  "authKey": "replace-with-strong-secret",
  "browserPath": "/usr/bin/chromium",
  "port": "8088",
  "dbPath": "./data/resume.db",
  "frontendOrigin": "https://wenemoji.com"
}
```

- [ ] **Step 3: Ignore `config.json`**

Add to `.gitignore`:

```gitignore
config.json
```

- [ ] **Step 4: Stop tracking `config.json` without deleting the local copy**

Run:

```bash
git rm --cached config.json
```

Then restore or preserve the working copy so local development still works.

- [ ] **Step 5: Re-run the tracked-state check**

Run: `git ls-files config.json config.example.json`
Expected after changes:
- `config.example.json` is tracked
- `config.json` is not tracked

- [ ] **Step 6: Commit the config-template slice**

```bash
git add .gitignore config.example.json
git rm --cached config.json
git commit -m "chore: track config example file"
```

## Chunk 3: Deployment Simplification

### Task 3: Remove environment-variable deployment dependency

**Files:**
- Modify: `deploy/docker-compose.yml`
- Delete: `deploy/env/app.env.example`
- Delete: `deploy/config/config.production.json.example`
- Modify: `deploy/scripts/deploy.sh`
- Modify: `deploy/scripts/healthcheck.sh`
- Modify: `docs/deploy/manual-docker-nginx.md`

- [ ] **Step 1: Write the failing/static deployment check**

Run and inspect:

```bash
docker compose -f deploy/docker-compose.yml config
```

Expected before changes:
- the `app` service still depends on `AUTH_KEY`
- the docs still reference `.env`

- [ ] **Step 2: Simplify `deploy/docker-compose.yml`**

Remove:
- `environment.AUTH_KEY`

Keep:
- `../config.json:/app/config.json:ro`
- `../data:/app/data`

If you want the simplest deploy flow, hardcode:

```yaml
ports:
  - '127.0.0.1:8088:8088'
```

- [ ] **Step 3: Delete the unused env and redundant config template files**

Delete:
- `deploy/env/app.env.example`
- `deploy/config/config.production.json.example`

- [ ] **Step 4: Simplify deployment scripts**

`deploy/scripts/deploy.sh` should become roughly:

```bash
docker compose -f deploy/docker-compose.yml up -d --build app
```

`deploy/scripts/healthcheck.sh` should default to:

```bash
HEALTHCHECK_URL="${HEALTHCHECK_URL:-http://127.0.0.1:8088/api/resume}"
```

- [ ] **Step 5: Rewrite the manual deployment guide**

The final documented flow should be:

```bash
cp config.example.json config.json
# edit config.json

docker compose -f deploy/docker-compose.yml up -d --build app
```

Remove all `.env` / `AUTH_KEY` setup instructions.

- [ ] **Step 6: Re-run Compose static validation**

Run: `docker compose -f deploy/docker-compose.yml config`
Expected: PASS with no env warnings for `AUTH_KEY`.

- [ ] **Step 7: Commit the deployment slice**

```bash
git add deploy/docker-compose.yml deploy/scripts/deploy.sh deploy/scripts/healthcheck.sh docs/deploy/manual-docker-nginx.md
git add -A deploy/env/app.env.example deploy/config/config.production.json.example
git commit -m "docs: simplify deploy config workflow"
```

## Chunk 4: Final Verification

### Task 4: Verify the full config-file-only workflow

**Files:**
- Modify only if verification reveals issues

- [ ] **Step 1: Run backend tests**

Run: `go test ./... -count=1`
Expected: PASS.

- [ ] **Step 2: Re-run frontend tests**

Run: `cd web && npm run test`
Expected: PASS.

- [ ] **Step 3: Re-run frontend build**

Run: `cd web && npm run build`
Expected: PASS.

- [ ] **Step 4: Re-run Compose static validation**

Run: `docker compose -f deploy/docker-compose.yml config`
Expected: PASS.

- [ ] **Step 5: Manual workflow checklist**

Confirm the documented runtime flow is coherent:
- clone repo
- `cp config.example.json config.json`
- edit `config.json` including `authKey`
- `docker compose -f deploy/docker-compose.yml up -d --build app`

- [ ] **Step 6: Commit any final verification fixes**

```bash
git add .gitignore config.example.json internal/config deploy docs/deploy
git commit -m "test: verify config json only workflow"
```

## Notes For The Implementer

- This change intentionally reverses the prior env-only `AUTH_KEY` design. Do not leave any hybrid loading behavior behind.
- Keep `config.json` as the runtime path, but stop tracking it in git.
- `config.example.json` must be the only tracked template.
- Remove the deployment env file entirely; the goal is a single config source, not a mixed strategy.
