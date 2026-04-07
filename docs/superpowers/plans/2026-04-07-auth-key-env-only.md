# AUTH_KEY Env-Only Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make `AUTH_KEY` the only source of the backend auth secret, removing `authKey` from `config.json` while keeping all other config fields file-based.

**Architecture:** Keep the existing `Config` struct and consumer code intact, but hydrate `Config.AuthKey` from `os.Getenv("AUTH_KEY")` after reading `config.json`. Update tests to set and clear environment state explicitly, then wire the deployment stack and docs to pass `AUTH_KEY` via container environment instead of committed config.

**Tech Stack:** Go, Docker Compose, Jenkins, Bash, Vitest-free backend tests, Markdown docs

---

## File Structure

### Backend config files

- Modify: `internal/config/config.go`
  Responsibility: Read `AUTH_KEY` from environment and stop treating `authKey` as a JSON-backed field requirement.
- Modify: `internal/config/config_test.go`
  Responsibility: Prove config loading now depends on `AUTH_KEY` from environment, not `authKey` in `config.json`.

### Deployment files

- Modify: `deploy/docker-compose.yml`
  Responsibility: Inject `AUTH_KEY` into the `app` container environment.
- Modify: `deploy/env/app.env.example`
  Responsibility: Document the required `AUTH_KEY` variable for deployments.
- Modify: `deploy/config/config.production.json.example`
  Responsibility: Remove `authKey` from the production config example so the repo no longer suggests committing the secret.
- Modify: `docs/deploy/docker-jenkins.md`
  Responsibility: Update server setup and deployment instructions so operators set `AUTH_KEY` in `deploy.env` instead of the production config file.

### Local config example

- Modify: `config.json`
  Responsibility: Remove the committed `authKey` field so local example config matches the new design.

## Chunk 1: Backend Config Loading

### Task 1: Switch auth secret loading from config file to environment

**Files:**
- Modify: `internal/config/config.go`
- Modify: `internal/config/config_test.go`

- [ ] **Step 1: Write the failing backend config tests**

Update or add tests in `internal/config/config_test.go` to assert:
- config loads successfully when `AUTH_KEY` is set and `config.json` has no `authKey`
- loading fails when `AUTH_KEY` is missing
- BOM-prefixed config files still load when `AUTH_KEY` is set

Example target test:

```go
func TestLoadUsesAuthKeyFromEnvironment(t *testing.T) {
	os.Setenv("AUTH_KEY", "resume-key")
	defer os.Unsetenv("AUTH_KEY")

	if err := os.WriteFile("config.json", []byte(`{"port":"9090"}`), 0o644); err != nil {
		t.Fatalf("WriteFile(config.json) returned error: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if cfg.AuthKey != "resume-key" {
		t.Fatalf("expected auth key from environment, got %q", cfg.AuthKey)
	}
}
```

- [ ] **Step 2: Run the focused config tests to verify they fail**

Run: `go test ./internal/config -count=1`
Expected: FAIL because `Load()` still requires `authKey` in `config.json`.

- [ ] **Step 3: Implement environment-only auth key loading**

Change `internal/config/config.go` so:
- `config.json` is still read and decoded
- `cfg.AuthKey = os.Getenv("AUTH_KEY")`
- missing env returns `errors.New("AUTH_KEY is required")`

Do not move any other fields to environment variables.

Implementation target:

```go
cfg.AuthKey = os.Getenv("AUTH_KEY")
if cfg.AuthKey == "" {
	return Config{}, errors.New("AUTH_KEY is required")
}
```

- [ ] **Step 4: Re-run the focused config tests**

Run: `go test ./internal/config -count=1`
Expected: PASS.

- [ ] **Step 5: Run the full Go suite to catch regressions**

Run: `go test ./... -count=1`
Expected: PASS.

- [ ] **Step 6: Commit the backend config slice**

```bash
git add internal/config/config.go internal/config/config_test.go
git commit -m "refactor: load auth key from environment"
```

## Chunk 2: Deployment Asset Updates

### Task 2: Update Docker and Jenkins deployment inputs to pass AUTH_KEY

**Files:**
- Modify: `deploy/docker-compose.yml`
- Modify: `deploy/env/app.env.example`
- Modify: `deploy/config/config.production.json.example`
- Modify: `docs/deploy/docker-jenkins.md`
- Modify: `config.json`

- [ ] **Step 1: Write the failing/static validation check for Compose expectations**

Update the deployment files first in the plan, then run Compose config validation expecting the `app` service to expose `AUTH_KEY`.

Validation command:

Run: `docker compose --env-file deploy/env/app.env.example -f deploy/docker-compose.yml config`
Expected before code change: PASS syntactically, but `AUTH_KEY` is absent from the rendered `app.environment` block.

- [ ] **Step 2: Update `deploy/env/app.env.example`**

Add:

```env
AUTH_KEY=replace-with-strong-secret
```

Keep all existing deployment variables.

- [ ] **Step 3: Update `deploy/docker-compose.yml`**

Inject the variable into the `app` service, for example:

```yaml
services:
  app:
    environment:
      AUTH_KEY: ${AUTH_KEY}
```

Do not add unrelated environment variables.

- [ ] **Step 4: Remove `authKey` from production and local config examples**

Update:
- `deploy/config/config.production.json.example`
- `config.json`

Both files should no longer include `authKey`.

- [ ] **Step 5: Update deployment documentation**

Edit `docs/deploy/docker-jenkins.md` so it tells operators to set `AUTH_KEY` in `/srv/jianli/deploy.env` and explicitly states that `config.production.json` no longer contains the auth secret.

- [ ] **Step 6: Re-run Compose validation**

Run: `docker compose --env-file deploy/env/app.env.example -f deploy/docker-compose.yml config`
Expected: PASS and rendered output includes `AUTH_KEY` under the `app` service environment.

- [ ] **Step 7: Commit the deployment slice**

```bash
git add deploy/docker-compose.yml deploy/env/app.env.example deploy/config/config.production.json.example docs/deploy/docker-jenkins.md config.json
git commit -m "ops: pass auth key through environment"
```

## Chunk 3: Final Verification

### Task 3: Verify the whole auth-key migration end to end

**Files:**
- Modify only if verification reveals issues in the files above

- [ ] **Step 1: Re-run backend tests**

Run: `go test ./... -count=1`
Expected: PASS.

- [ ] **Step 2: Re-run frontend tests for safety**

Run: `cd web && npm run test`
Expected: PASS. Frontend behavior should not change, but this guards against unintended breakage.

- [ ] **Step 3: Re-run frontend build**

Run: `cd web && npm run build`
Expected: PASS.

- [ ] **Step 4: Re-run Compose static validation**

Run: `docker compose --env-file deploy/env/app.env.example -f deploy/docker-compose.yml config`
Expected: PASS.

- [ ] **Step 5: Manual smoke checklist**

Confirm manually when running locally or on server:
- backend fails fast if `AUTH_KEY` is missing
- backend starts if `AUTH_KEY` is set and `config.json` has no `authKey`
- admin login still works using the same key value
- Docker deployment docs and env examples are self-consistent

- [ ] **Step 6: Commit any final verification fixes**

```bash
git add internal/config deploy docs/deploy config.json
git commit -m "test: verify auth key env-only configuration"
```

## Notes For The Implementer

- Follow @test-driven-development strictly for the config change: update the failing tests first, then change `Load()`.
- Do not leave backward compatibility for `config.json.authKey`; the approved design is env-only.
- Keep `Config.AuthKey` in the struct to minimize blast radius on existing call sites.
- Avoid touching unrelated deployment assets beyond the `AUTH_KEY` flow.
