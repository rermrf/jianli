# Docker + Jenkins Deployment Guide

## 1. Server prerequisites

Install on the Linux server:

- Docker Engine
- Docker Compose v2
- Git

The deployment expects these host directories:

- `/srv/jianli/app-src`
- `/srv/jianli/config`
- `/srv/jianli/data`
- `/srv/jenkins_home`

Create them:

```bash
sudo mkdir -p /srv/jianli/app-src /srv/jianli/config /srv/jianli/data /srv/jenkins_home
sudo chown -R $USER:$USER /srv/jianli /srv/jenkins_home
```

## 2. Production config

Copy the example config and env file:

```bash
cp deploy/config/config.production.json.example /srv/jianli/config/config.production.json
cp deploy/env/app.env.example /srv/jianli/deploy.env
```

Then edit:

- `/srv/jianli/config/config.production.json`
- `/srv/jianli/deploy.env`

Important values:

- `authKey`: replace with a strong secret
- `browserPath`: keep `/usr/bin/chromium` unless your server uses another binary
- `APP_CONFIG_PATH`: should point to `/srv/jianli/config/config.production.json`
- `APP_DATA_DIR`: should point to `/srv/jianli/data`

## 3. Start Jenkins

From the project repo root on the server:

```bash
docker compose --env-file /srv/jianli/deploy.env -f deploy/docker-compose.yml up -d jenkins
```

Jenkins will be exposed on `http://<server>:8081` by default.

## 4. Create Jenkins job

Create a Pipeline job pointing to this repository.

The pipeline uses:

- `Jenkinsfile` at repo root
- Docker socket from the Jenkins container
- host path `/srv/jianli`

When Jenkins runs:

1. checks out `master`
2. syncs the source into `/srv/jianli/app-src`
3. runs backend tests
4. runs frontend tests
5. builds `app` and `nginx` images
6. waits for manual approval
7. deploys with Docker Compose
8. health-checks `/api/resume`

## 5. First release

Once Jenkins is up and the job is configured:

- push to `master`
- wait for the pipeline to complete test/build stages
- click the manual `Deploy` approval in Jenkins

Jenkins will run:

```bash
/srv/jianli/app-src/deploy/scripts/deploy.sh /srv/jianli/app-src /srv/jianli/deploy.env
```

## 6. Validation checklist

After deployment, verify:

- `http://<server>/` loads the site
- `http://<server>/api/resume` returns 200
- admin login works
- avatar upload works
- `/uploads/...` assets load
- PDF export works

## 7. Manual rollback

This setup does not automate rollback. To roll back:

1. check out the previous good commit
2. rerun the Jenkins job
3. approve deployment

Or manually redeploy from a known-good source tree:

```bash
cd /srv/jianli/app-src
# restore previous commit however you prefer

docker compose --env-file /srv/jianli/deploy.env -f deploy/docker-compose.yml up -d --build app nginx
```

## 8. Notes

- Jenkins runs with Docker socket access. This is acceptable for this single-host setup but should be tightened later if the environment grows.
- PDF export depends on the browser path inside the runtime container. If Chromium is missing or the path is wrong, PDF export will fail even if the site otherwise works.
