#!/usr/bin/env bash
set -euo pipefail

DEPLOY_ROOT="${1:-/srv/jianli/app-src}"
ENV_FILE="${2:-/srv/jianli/deploy.env}"
COMPOSE_FILE="$DEPLOY_ROOT/deploy/docker-compose.yml"

if [[ ! -f "$COMPOSE_FILE" ]]; then
  echo "docker compose file not found: $COMPOSE_FILE" >&2
  exit 1
fi

if [[ ! -f "$ENV_FILE" ]]; then
  echo "deployment env file not found: $ENV_FILE" >&2
  exit 1
fi

docker compose --env-file "$ENV_FILE" -f "$COMPOSE_FILE" up -d --build app
