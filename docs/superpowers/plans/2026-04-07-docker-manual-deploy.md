# Docker Manual Deploy Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Simplify deployment to a single `app` container that reuses the existing server-side Nginx container, removing Jenkins and the project-internal Nginx service.

**Architecture:** Keep the multi-stage build for the `app` runtime image, but remove the repo-managed `nginx` and `jenkins` service targets from Compose and docs. Deployment becomes a manual `git pull` plus `docker compose up -d --build` workflow, with the server's existing Nginx container reverse-proxying to the app container on a host-only port.

**Tech Stack:** Docker, Docker Compose, Go, Node.js, Vite, Chromium, Markdown docs

---

## File Structure

### Container files

- Modify: `Dockerfile`
  Responsibility: Remove the repo-managed `nginx` target and keep only the `app` runtime image path.
- Modify: `.dockerignore`
  Responsibility: Keep existing ignore rules unless the simplification reveals anything unnecessary; likely no change needed.

### Deployment files

- Modify: `deploy/docker-compose.yml`
  Responsibility: Reduce the stack to a single `app` service and bind it only to a host-local port for the external Nginx container to reach.
- Modify: `deploy/env/app.env.example`
  Responsibility: Remove Jenkins-related variables and keep only the app deployment inputs, including `AUTH_KEY`.
- Modify: `deploy/config/config.production.json.example`
  Responsibility: Ensure it matches the simplified single-container runtime assumptions.
- Modify: `deploy/scripts/deploy.sh`
  Responsibility: Simplify rollout to only build and restart the `app` service.
- Modify: `deploy/scripts/healthcheck.sh`
  Responsibility: Keep healthcheck aligned with the app-only stack.

### Remove obsolete CI/CD files

- Delete: `Jenkinsfile`
  Responsibility: Remove repo-managed Jenkins pipeline now that CI/CD is explicitly out of scope.
- Delete: `deploy/jenkins/Dockerfile`
  Responsibility: Remove the Jenkins container image definition.
- Optional delete if empty after cleanup: `deploy/jenkins/`

### Documentation files

- Modify: `docs/deploy/docker-jenkins.md`
  Responsibility: Replace Jenkins instructions with manual deployment steps and explain how to wire the existing server Nginx container to the app container.
- Create: `docs/deploy/manual-docker-nginx.md`
  Responsibility: Provide the new canonical deployment guide if splitting docs is clearer than mutating the old Jenkins-specific guide.

## Chunk 1: Simplify Container Build And Compose

### Task 1: Remove internal Nginx/Jenkins targets and keep only the app runtime image

**Files:**
- Modify: `Dockerfile`
- Modify: `deploy/docker-compose.yml`
- Modify: `deploy/env/app.env.example`

- [ ] **Step 1: Write the failing/static validation check**

Run: `docker compose --env-file deploy/env/app.env.example -f deploy/docker-compose.yml config`
Expected before changes: PASS syntactically, but output still contains unwanted `nginx` and `jenkins` services.

- [ ] **Step 2: Simplify `Dockerfile`**

Remove the final `nginx` target and leave:
- `web-build`
- `go-build`
- `app`

The file should still produce a runtime image with:
- `/app/bin/jianli-server`
- `/app/web/dist`
- Chromium installed for PDF export

- [ ] **Step 3: Simplify `deploy/docker-compose.yml`**

Reduce to one service:
- `app`

Recommended characteristics:
- host bind like `127.0.0.1:8080:8080`
- `AUTH_KEY` passed via environment
- mounts for `config.production.json` and `data/`
- no internal nginx service
- no jenkins service

- [ ] **Step 4: Simplify `deploy/env/app.env.example`**

Keep only app-relevant variables, for example:

```env
APP_CONFIG_PATH=/srv/jianli/config/config.production.json
APP_DATA_DIR=/srv/jianli/data
APP_BIND=127.0.0.1:8080:8080
AUTH_KEY=replace-with-strong-secret
APP_HEALTHCHECK_URL=http://127.0.0.1:8080/api/resume
```

Remove Jenkins-specific values.

- [ ] **Step 5: Re-run Compose validation**

Run: `docker compose --env-file deploy/env/app.env.example -f deploy/docker-compose.yml config`
Expected: PASS and rendered output shows only the `app` service.

- [ ] **Step 6: Commit the runtime simplification slice**

```bash
git add Dockerfile deploy/docker-compose.yml deploy/env/app.env.example
git commit -m "ops: simplify deployment to single app container"
```

## Chunk 2: Remove Jenkins-Specific Assets

### Task 2: Delete obsolete Jenkins deployment files

**Files:**
- Delete: `Jenkinsfile`
- Delete: `deploy/jenkins/Dockerfile`
- Optionally delete: `deploy/jenkins/`

- [ ] **Step 1: Verify these files are no longer referenced**

Search references before deleting.

Run: `rg -n "Jenkinsfile|deploy/jenkins|jenkins" docs deploy .`
Expected: identify references that need to be updated in docs or scripts.

- [ ] **Step 2: Remove the obsolete Jenkins files**

Delete the pipeline file and Jenkins image definition.

- [ ] **Step 3: Re-run reference search**

Run: `rg -n "Jenkinsfile|deploy/jenkins|jenkins" docs deploy .`
Expected: any remaining references are intentional documentation mentions of the old approach or are removed/updated.

- [ ] **Step 4: Commit the cleanup slice**

```bash
git add -A Jenkinsfile deploy/jenkins
git commit -m "chore: remove jenkins deployment assets"
```

## Chunk 3: Update Scripts And Documentation

### Task 3: Rewrite deployment scripts and docs for manual rollout behind existing Nginx

**Files:**
- Modify: `deploy/scripts/deploy.sh`
- Modify: `deploy/scripts/healthcheck.sh`
- Modify or replace: `docs/deploy/docker-jenkins.md`
- Create optionally: `docs/deploy/manual-docker-nginx.md`

- [ ] **Step 1: Write the failing/script validation step**

Run shell syntax checks before editing so you have a baseline.

Examples:
- `bash -n deploy/scripts/deploy.sh`
- `bash -n deploy/scripts/healthcheck.sh`

Expected before changes: scripts are syntactically valid but still mention the old multi-service deployment assumptions.

- [ ] **Step 2: Simplify `deploy/scripts/deploy.sh`**

Responsibilities after rewrite:
- accept deploy root and env file arguments
- run `docker compose build app`
- run `docker compose up -d app`
- no nginx or Jenkins orchestration

- [ ] **Step 3: Keep `deploy/scripts/healthcheck.sh` aligned**

The healthcheck should continue curling the configured app endpoint, but update any defaults or comments so they reflect the app-only stack.

- [ ] **Step 4: Rewrite deployment docs**

Replace Jenkins-centric instructions with:
- server directory prep
- production config placement
- `deploy.env` contents
- manual deployment commands:
  - `git pull`
  - `docker compose --env-file /srv/jianli/deploy.env up -d --build`
- how to wire the existing Nginx container to `127.0.0.1:8080`
- validation checklist
- manual rollback procedure

The guide should explicitly state:
- no repo-managed Jenkins
- no repo-managed nginx
- existing server Nginx container remains the public entrypoint

- [ ] **Step 5: Add the exact Nginx reverse-proxy snippet for `wenemoji.com`**

Document one concrete server block or location snippet pointing the existing Nginx container to the app container's host-bound port.

- [ ] **Step 6: Re-run shell syntax checks**

Run:
- `bash -n deploy/scripts/deploy.sh`
- `bash -n deploy/scripts/healthcheck.sh`

Expected: PASS.

- [ ] **Step 7: Commit the manual deployment slice**

```bash
git add deploy/scripts/deploy.sh deploy/scripts/healthcheck.sh docs/deploy
git commit -m "docs: switch to manual docker deployment"
```

## Chunk 4: Final Verification

### Task 4: Verify the simplified deployment assets end to end

**Files:**
- Modify only if verification reveals issues in Docker or docs assets

- [ ] **Step 1: Run backend tests**

Run: `go test ./... -count=1`
Expected: PASS.

- [ ] **Step 2: Run frontend tests**

Run: `cd web && npm run test`
Expected: PASS.

- [ ] **Step 3: Run frontend build**

Run: `cd web && npm run build`
Expected: PASS.

- [ ] **Step 4: Run Compose static validation**

Run: `docker compose --env-file deploy/env/app.env.example -f deploy/docker-compose.yml config`
Expected: PASS and output contains only the `app` service.

- [ ] **Step 5: If Docker daemon is available, smoke-check the app service locally**

Run: `docker compose --env-file deploy/env/app.env.example -f deploy/docker-compose.yml up -d --build app`
Then verify:
- `curl http://127.0.0.1:8080/api/resume`

If Docker daemon is unavailable in the execution environment, record that explicitly instead of claiming success.

- [ ] **Step 6: Commit any final verification fixes**

```bash
git add Dockerfile deploy docs/deploy
git commit -m "test: verify manual docker deployment"
```

## Notes For The Implementer

- Keep the deployment intentionally small. Do not preserve Jenkins scaffolding "just in case".
- Reuse the existing server Nginx; this repo should not keep an internal nginx service after the simplification.
- Preserve `AUTH_KEY` env-only behavior already merged into `master`.
- Be explicit in docs that the existing server Nginx container must reverse-proxy to this app service's host port.
