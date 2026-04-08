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
mkdir -p /root/wen/go/jianli/data/uploads/avatars
```

This deployment uses files directly from the project directory:

- `/root/wen/go/jianli/config.json`
- `/root/wen/go/jianli/data`

## 3. Create the runtime config

After cloning the repo, generate your local runtime config:

```bash
cd /root/wen/go/jianli
cp config.example.json config.json
```

Then edit `config.json`.

Recommended production content:

```json
{
  "browserPath": "/usr/bin/chromium",
  "port": "8088",
  "dbPath": "./data/resume.db",
  "frontendOrigin": "https://wenemoji.com"
}
```

Notes:
- do not add `authKey` into `config.json`
- `AUTH_KEY` is environment-only
- keep `port` as `8088`

## 4. Environment variables

Create a local env file if you want:

```bash
cp deploy/env/app.env.example .env
```

Then set at least:

```env
AUTH_KEY=replace-with-strong-secret
APP_BIND=127.0.0.1:8088:8088
APP_HEALTHCHECK_URL=http://127.0.0.1:8088/api/resume
```

If you do not want a `.env` file, you can export `AUTH_KEY` inline when running Docker Compose.

## 5. First deployment

```bash
cd /root/wen/go/jianli
docker compose --env-file .env -f deploy/docker-compose.yml up -d --build app
```

Or without `.env`:

```bash
cd /root/wen/go/jianli
AUTH_KEY='replace-with-strong-secret' docker compose -f deploy/docker-compose.yml up -d --build app
```

## 6. Routine deployment

```bash
cd /root/wen/go/jianli
git pull
cp config.example.json config.json # if config.json does not exist yet on this server

docker compose --env-file .env -f deploy/docker-compose.yml up -d --build app
```

## 7. Health check

```bash
APP_HEALTHCHECK_URL=http://127.0.0.1:8088/api/resume bash deploy/scripts/healthcheck.sh
```

Expected: exit code `0`.

## 8. Existing Nginx container integration

Your existing Nginx container should reverse-proxy `wenemoji.com` to the host-local app port `127.0.0.1:8088`.

If that Nginx container can access the host network, add a server block like:

```nginx
server {
    listen 80;
    server_name wenemoji.com www.wenemoji.com;

    location / {
        proxy_pass http://host.docker.internal:8088;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    location /api/ {
        proxy_pass http://host.docker.internal:8088;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    location /uploads/ {
        proxy_pass http://host.docker.internal:8088;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## 9. Validation checklist

After deployment, verify:

- `https://wenemoji.com/` loads
- `https://wenemoji.com/api/resume` returns 200
- admin login works
- avatar upload works
- `/uploads/...` assets load
- PDF export works

## 10. Manual rollback

Rollback is manual:

```bash
cd /root/wen/go/jianli
# restore a known-good commit or branch

docker compose --env-file .env -f deploy/docker-compose.yml up -d --build app
```

## 11. Notes

- This repo no longer manages Jenkins or a project-internal Nginx container.
- The backend requires `AUTH_KEY` from environment and will not start if it is missing.
- PDF export depends on Chromium existing in the app image and `browserPath` matching that path.
- `config.example.json` is the only tracked config template; `config.json` is local-only.
