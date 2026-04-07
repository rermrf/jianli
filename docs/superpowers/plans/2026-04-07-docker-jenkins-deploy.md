# Docker Jenkins Deploy Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Fully containerize the project and add Jenkins-driven CI/CD on a single Linux server, with automatic test/build on `master` and manual approval before deployment.

**Architecture:** Build the frontend and backend into one runtime `app` image, front it with an `nginx` container, and run Jenkins in a separate container on the same host with Docker socket access. Persist `config.json`, `data/`, and Jenkins home via host-mounted volumes so deployments are repeatable and state survives container replacement.

**Tech Stack:** Docker, Docker Compose, Jenkins, Nginx, Go, Node.js, Vite, Chromium, Bash

---

## File Structure

### Container and runtime files

- Create: `Dockerfile`
  Responsibility: Multi-stage build for frontend assets, Go binary, and final runtime image with Linux browser support for PDF export.
- Create: `.dockerignore`
  Responsibility: Exclude git metadata, local data, worktrees, build artifacts, and `node_modules` from Docker build context.
- Create: `deploy/docker-compose.yml`
  Responsibility: Define `jenkins`, `app`, and `nginx` services, networks, ports, volumes, and restart behavior.
- Create: `deploy/nginx/default.conf`
  Responsibility: Serve SPA static files and reverse proxy `/api` and `/uploads` to `app`.
- Create: `deploy/env/app.env.example`
  Responsibility: Show required production environment/config values referenced by Compose or deployment docs.
- Create: `deploy/config/config.production.json.example`
  Responsibility: Production-safe sample config showing Linux `browserPath`, persisted DB path, and production origin.

### CI/CD files

- Create: `Jenkinsfile`
  Responsibility: Checkout, test, build image, wait for manual approval, deploy with Compose, and health-check the app.
- Create: `deploy/scripts/healthcheck.sh`
  Responsibility: Verify the deployed service responds on `/api/resume`.
- Create: `deploy/scripts/deploy.sh`
  Responsibility: Wrap `docker compose` deployment commands so Jenkins steps stay readable and deterministic.

### Documentation files

- Create: `docs/deploy/docker-jenkins.md`
  Responsibility: Server preparation, Jenkins setup, Compose startup, config placement, and first-release checklist.

### Existing code to verify or minimally adjust

- Reference: `cmd/server/main.go`
  Ensure the app can still read `config.json` from the runtime working directory.
- Reference: `internal/config/config.go`
  Confirm config defaults do not break containerized production.
- Reference: `internal/pdf/export.go`
  Ensure the configured `browserPath` works with Linux runtime image.
- Reference: `web/package.json`
  Confirm the frontend build/test commands used by Docker and Jenkins are accurate.

## Chunk 1: Runtime Containerization

### Task 1: Build a production-ready app image

**Files:**
- Create: `Dockerfile`
- Create: `.dockerignore`
- Test with: local Docker build command

- [ ] **Step 1: Write the failing container smoke check**

Add a local verification note to the plan and immediately run a build command that is expected to fail because `Dockerfile` does not exist yet.

Run: `docker build -t jianli:test .`
Expected: FAIL with missing `Dockerfile`.

- [ ] **Step 2: Create `.dockerignore`**

Include at least:

```dockerignore
.git
.worktrees
web/node_modules
web/dist
bin
data
*.db
```

Also exclude temporary build outputs and local editor files.

- [ ] **Step 3: Create a multi-stage `Dockerfile`**

Stages:
- `web-build`: install frontend deps and run `npm run build`
- `go-build`: download modules and build `./cmd/server`
- `runtime`: install Linux browser package and copy runtime artifacts

Key runtime requirements:
- working directory `/app`
- copy built binary to `/app/bin/jianli-server`
- copy frontend dist to `/app/web/dist`
- create `/app/data/uploads/avatars`
- include Chromium or compatible browser for PDF export

Example runtime command target:

```dockerfile
CMD ["/app/bin/jianli-server"]
```

- [ ] **Step 4: Re-run the Docker build to verify it now succeeds**

Run: `docker build -t jianli:test .`
Expected: PASS and produce a local image.

- [ ] **Step 5: Smoke-run the image with mounted config/data**

Run a one-off container using a production-style config mount.
Expected: the app starts and binds on container port `8080`.

- [ ] **Step 6: Commit the app container slice**

```bash
git add Dockerfile .dockerignore
git commit -m "build: add app container image"
```

## Chunk 2: Compose And Reverse Proxy

### Task 2: Define Compose services and Nginx routing

**Files:**
- Create: `deploy/docker-compose.yml`
- Create: `deploy/nginx/default.conf`
- Create: `deploy/env/app.env.example`
- Create: `deploy/config/config.production.json.example`

- [ ] **Step 1: Write the failing Compose validation check**

Run: `docker compose -f deploy/docker-compose.yml config`
Expected: FAIL because the Compose file does not exist yet.

- [ ] **Step 2: Create production config examples**

`deploy/config/config.production.json.example` should show values like:

```json
{
  "authKey": "replace-with-strong-secret",
  "browserPath": "/usr/bin/chromium",
  "port": "8080",
  "dbPath": "/app/data/resume.db",
  "frontendOrigin": "https://resume.example.com"
}
```

If you choose to externalize some values via env, document that split clearly in `app.env.example`.

- [ ] **Step 3: Create `deploy/docker-compose.yml`**

Define:
- `app`
- `nginx`
- `jenkins`

Key requirements:
- `app` mounts production `config.json` and persistent `data/`
- `nginx` mounts the app image or static/proxy config and exposes `80:80`
- `jenkins` mounts `jenkins_home` and Docker socket
- internal network shared by all services
- restart policy for long-running services

- [ ] **Step 4: Create `deploy/nginx/default.conf`**

Must support:
- `try_files $uri /index.html;`
- `/api/` proxy to `app:8080`
- `/uploads/` proxy to `app:8080`

- [ ] **Step 5: Validate Compose configuration**

Run: `docker compose -f deploy/docker-compose.yml config`
Expected: PASS.

- [ ] **Step 6: Bring up only app + nginx locally for smoke verification**

Run: `docker compose -f deploy/docker-compose.yml up -d app nginx`
Then verify:
- `curl http://127.0.0.1/api/resume`
- front page serves through nginx

Expected: PASS.

- [ ] **Step 7: Commit the deployment topology slice**

```bash
git add deploy/docker-compose.yml deploy/nginx/default.conf deploy/env/app.env.example deploy/config/config.production.json.example
git commit -m "ops: add compose deployment stack"
```

## Chunk 3: Jenkins Pipeline And Deployment Scripts

### Task 3: Add Jenkins pipeline and deployment automation

**Files:**
- Create: `Jenkinsfile`
- Create: `deploy/scripts/deploy.sh`
- Create: `deploy/scripts/healthcheck.sh`

- [ ] **Step 1: Write the failing pipeline validation check**

Run a syntax-oriented smoke check by viewing the Jenkinsfile path or running a linter if available.
Expected: FAIL because `Jenkinsfile` does not exist yet.

- [ ] **Step 2: Create `deploy/scripts/healthcheck.sh`**

Script should:
- curl `http://127.0.0.1/api/resume`
- fail non-zero if the endpoint is unreachable or non-200

Example:

```bash
#!/usr/bin/env bash
set -euo pipefail
curl --fail --silent http://127.0.0.1/api/resume >/dev/null
```

- [ ] **Step 3: Create `deploy/scripts/deploy.sh`**

Responsibilities:
- switch into `deploy/`
- run `docker compose build app`
- run `docker compose up -d app nginx`
- optionally restart nginx if config changes are mounted from host

- [ ] **Step 4: Create `Jenkinsfile`**

Required stages:
- `Checkout`
- `Test Backend`
- `Test Frontend`
- `Build Image`
- `Approve Deploy` (manual `input` step)
- `Deploy`
- `Health Check`

Use exact commands already proven in repo:
- `go test ./... -count=1`
- `cd web && npm ci && npm run test`
- `docker build ...` or `docker compose build app`
- `deploy/scripts/deploy.sh`
- `deploy/scripts/healthcheck.sh`

- [ ] **Step 5: Validate pipeline script readability and executable permissions**

Run:
- `bash deploy/scripts/healthcheck.sh` against a running local stack
- `bash deploy/scripts/deploy.sh` against the local Docker environment (or dry-run equivalent if needed)

Expected: PASS.

- [ ] **Step 6: Commit the CI/CD slice**

```bash
git add Jenkinsfile deploy/scripts/deploy.sh deploy/scripts/healthcheck.sh
git commit -m "ci: add jenkins deployment pipeline"
```

## Chunk 4: Deployment Documentation

### Task 4: Document server setup and first release procedure

**Files:**
- Create: `docs/deploy/docker-jenkins.md`

- [ ] **Step 1: Write the documentation with exact server steps**

Include:
- required packages (`docker`, `docker compose`)
- host directories to create
- where to place production `config.json`
- how to start Jenkins container the first time
- how to configure Jenkins job or multibranch pipeline
- how to mount Docker socket safely enough for this setup
- first deployment checklist
- how to verify PDF export in production

- [ ] **Step 2: Include rollback guidance even if manual**

At minimum:
- note how to rerun Compose with a previous image tag
- or how to restore a previous commit and redeploy

- [ ] **Step 3: Commit the deployment docs**

```bash
git add docs/deploy/docker-jenkins.md
git commit -m "docs: add docker jenkins deployment guide"
```

## Chunk 5: Final Verification

### Task 5: Verify the complete deployment setup end to end

**Files:**
- Modify only if verification reveals issues in deployment files

- [ ] **Step 1: Re-run backend tests**

Run: `go test ./... -count=1`
Expected: PASS.

- [ ] **Step 2: Re-run frontend tests**

Run: `cd web && npm run test`
Expected: PASS.

- [ ] **Step 3: Re-run frontend build**

Run: `cd web && npm run build`
Expected: PASS.

- [ ] **Step 4: Rebuild the production image from scratch**

Run: `docker build -t jianli:final .`
Expected: PASS.

- [ ] **Step 5: Validate Compose one last time**

Run: `docker compose -f deploy/docker-compose.yml config`
Expected: PASS.

- [ ] **Step 6: Bring up the local stack and verify the critical paths**

Run the stack and confirm:
- `/` loads through nginx
- `/api/resume` returns 200
- `/uploads/...` serves uploaded assets if present
- PDF export works with the configured Linux browser path

- [ ] **Step 7: Commit any final verification fixes**

```bash
git add Dockerfile .dockerignore deploy Jenkinsfile docs/deploy
git commit -m "test: verify docker jenkins deployment"
```

## Notes For The Implementer

- Follow @test-driven-development where practical for deployment assets too: run the command first, watch it fail, then add the minimal file needed to make it pass.
- Do not bake production secrets into the image. Persist config and data via mounted files/directories.
- Do not assume Windows browser paths in containers; production must use Linux browser paths.
- Keep the first version simple: one app image, one nginx service, one jenkins service. Do not introduce registries, staging, or extra environments yet.
