#!/usr/bin/env bash
set -euo pipefail

HEALTHCHECK_BASE_URL="${HEALTHCHECK_BASE_URL:-http://127.0.0.1:8088}"
HEALTHCHECK_PUBLIC_URL="${HEALTHCHECK_PUBLIC_URL:-${HEALTHCHECK_BASE_URL%/}/}"
HEALTHCHECK_SPA_URL="${HEALTHCHECK_SPA_URL:-${HEALTHCHECK_BASE_URL%/}/drafts/healthcheck}"
HEALTHCHECK_API_URL="${HEALTHCHECK_API_URL:-${HEALTHCHECK_BASE_URL%/}/api/resume}"

curl --fail --silent --show-error "$HEALTHCHECK_PUBLIC_URL" >/dev/null
curl --fail --silent --show-error "$HEALTHCHECK_SPA_URL" >/dev/null
curl --fail --silent --show-error "$HEALTHCHECK_API_URL" >/dev/null
