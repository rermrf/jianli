# Docker Manual Deployment Guide

## 1. Server prerequisites

Install on the Linux server:

- Docker Engine
- Docker Compose v2
- Git

This deployment assumes you already have a long-running Nginx container that owns ports `80/443` and TLS.

## 2. Required host directories

Create:

```bash
sudo mkdir -p /srv/jianli/app-src /srv/jianli/config /srv/jianli/data
sudo chown -R $USER:$USER /srv/jianli
```

These are used for:

- `/srv/jianli/app-src` -> checked-out project source
- `/srv/jianli/config` -> production config file
- `/srv/jianli/data` -> SQLite DB and uploads

## 3. Production config

Copy the example files:

```bash
cp deploy/config/config.production.json.example /srv/jianli/config/config.production.json
cp deploy/env/app.env.example /srv/jianli/deploy.env
```

Then edit them.

`/srv/jianli/config/config.production.json`
- keep `browserPath` as a Linux browser path such as `/usr/bin/chromium`
- keep `frontendOrigin` as `https://wenemoji.com`
- do not add `authKey` here

`/srv/jianli/deploy.env`
- set `AUTH_KEY` to a strong secret
- keep `APP_BIND=127.0.0.1:8080:8080` so the app is reachable only from the host and your existing Nginx container

## 4. First deployment

```bash
cd /srv/jianli
git clone <your-repo-url> app-src
cd app-src
git checkout master

docker compose --env-file /srv/jianli/deploy.env -f deploy/docker-compose.yml up -d --build app
```

## 5. Routine deployment

```bash
cd /srv/jianli/app-src
git pull
docker compose --env-file /srv/jianli/deploy.env -f deploy/docker-compose.yml up -d --build app
```

## 6. Health check

```bash
APP_HEALTHCHECK_URL=http://127.0.0.1:8080/api/resume bash deploy/scripts/healthcheck.sh
```

Expected: exit code `0`.

## 7. Existing Nginx container integration

Your existing Nginx container should reverse-proxy `wenemoji.com` to the host-local app port `127.0.0.1:8080`.

If that Nginx container can access the host network, add a server block similar to:

```nginx
server {
    listen 80;
    server_name wenemoji.com www.wenemoji.com;

    location / {
        proxy_pass http://host.docker.internal:8080;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    location /api/ {
        proxy_pass http://host.docker.internal:8080;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    location /uploads/ {
        proxy_pass http://host.docker.internal:8080;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

If `host.docker.internal` is unavailable in your Linux Docker setup, use one of:
- the Docker bridge gateway IP
- a shared custom network between your existing Nginx container and the app container
- host networking rules that fit your server layout

## 8. Validation checklist

After deployment, verify:

- `https://wenemoji.com/` loads
- `https://wenemoji.com/api/resume` returns 200
- admin login works
- avatar upload works
- `/uploads/...` assets load
- PDF export works

## 9. Manual rollback

Rollback is manual:

```bash
cd /srv/jianli/app-src
# restore a known-good commit or branch

docker compose --env-file /srv/jianli/deploy.env -f deploy/docker-compose.yml up -d --build app
```

## 10. Notes

- This repo no longer manages Jenkins or a project-internal Nginx container.
- The backend now requires `AUTH_KEY` from environment and will not start if that variable is missing.
- PDF export depends on Chromium existing in the app image and `browserPath` matching that path.
