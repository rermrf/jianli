#!/usr/bin/env bash
set -euo pipefail

DEPLOY_ROOT="${1:-/srv/jianli/app-src}"
COMPOSE_FILE="$DEPLOY_ROOT/deploy/docker-compose.yml"

if [[ ! -f "$COMPOSE_FILE" ]]; then
  echo "docker compose file not found: $COMPOSE_FILE" >&2
  exit 1
fi

docker compose -f "$COMPOSE_FILE" up -d --build app
