#!/usr/bin/env bash
set -euo pipefail

HEALTHCHECK_URL="${APP_HEALTHCHECK_URL:-http://127.0.0.1/api/resume}"

curl --fail --silent "$HEALTHCHECK_URL" >/dev/null
