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
  "authKey": "replace-with-strong-secret",
  "browserPath": "/usr/bin/chromium",
  "port": "8088",
  "dbPath": "./data/resume.db",
  "frontendOrigin": "https://wenemoji.com",
  "ip2RegionPath": "./data/ip2region_v4.xdb",
  "ip2RegionV6Path": "./data/ip2region_v6.xdb",
  "trustedProxies": ["127.0.0.1", "::1", "172.16.0.0/12"]
}
```

Notes:
- all runtime business config now lives in `config.json`
- `config.example.json` is the only tracked template
- keep `port` as `8088`
- `trustedProxies` must include the IP or CIDR of the Nginx container seen by the app; the default includes Docker's usual `172.16.0.0/12` bridge range. For a custom Docker subnet, replace it with that subnet. Use `[]` for direct access without a reverse proxy.
- The image downloads pinned IPv4 and IPv6 ip2region databases at build time and stores them in `/app/geoip`. Files specified in `config.json` take priority; the bundled files are used when the configured files are unavailable. Mounting `data` does not hide the bundled databases.
- For a local Go run, download both databases first with `bash deploy/scripts/download-geoip.sh ./data`. The server loads them once at startup; restart after replacing them.
- Opening the visitor list automatically fills and saves previously unknown regions when the stored IP can be resolved. Local/private addresses display `本地网络`; old rows that contain only a proxy IP cannot recover the original client IP.
- Visitor records retain the full IPv4/IPv6 address and request User-Agent. The page displays country/province/city, operator, device type, operating system, and available system/browser versions. Windows NT 10.0 cannot distinguish Windows 10 from 11 and is displayed as `Windows 10/11`. Existing system fields remain usable; versions that were never recorded cannot be reconstructed.
- Recent-visit deduplication uses both IP and User-Agent so different systems or browsers behind one public IP retain separate records. New metadata columns are migrated automatically on server startup.
- Daily analytics use Beijing calendar dates for the chart, summary cards and visitor list. The chart defaults to distinct visitors per day (IP-based), can switch to recorded visit counts, fills zero-visit days and provides a daily detail table. `7天` and `30天` include today; `全部` starts at the earliest recorded visit.

## 4. First deployment

```bash
cd /root/wen/go/jianli
docker compose -f deploy/docker-compose.yml up -d --build app
```

## 5. Routine deployment

```bash
cd /root/wen/go/jianli
git pull
cp config.example.json config.json # if config.json does not exist yet on this server

docker compose -f deploy/docker-compose.yml up -d --build app
```

## 6. Health check

```bash
bash deploy/scripts/healthcheck.sh
```

The script now verifies three paths on the app process:

- `/` returns the public homepage shell
- `/drafts/healthcheck` returns the SPA fallback shell
- `/api/resume` returns the public resume payload

Default target:

```bash
http://127.0.0.1:8088
```

Optional overrides:

- `HEALTHCHECK_BASE_URL`
- `HEALTHCHECK_PUBLIC_URL`
- `HEALTHCHECK_SPA_URL`
- `HEALTHCHECK_API_URL`

Expected: exit code `0`.

## 7. Existing Nginx container integration

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

## 8. Validation checklist

After deployment, verify:

- `https://wenemoji.com/` loads
- `https://wenemoji.com/drafts/healthcheck` returns the SPA shell instead of a 404
- `https://wenemoji.com/api/resume` returns 200
- admin login works
- avatar upload works
- `/uploads/...` assets load
- PDF export works

## 9. Manual rollback

Rollback is manual:

```bash
cd /root/wen/go/jianli
# restore a known-good commit or branch

docker compose -f deploy/docker-compose.yml up -d --build app
```

## 10. Notes

- This repo no longer manages Jenkins or a project-internal Nginx container.
- The app now reads all runtime settings from `config.json`.
- `config.json` is local-only and should be created from `config.example.json`.
- PDF export depends on Chromium existing in the app image and `browserPath` matching that path.
