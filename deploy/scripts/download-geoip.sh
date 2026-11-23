#!/usr/bin/env bash
set -euo pipefail

# Pin databases to the same upstream revision as the Go ip2region binding.
IP2REGION_REF="593e02867796803c99fcddc8aaedda1f323e56b7"
OUTPUT_DIR="${1:-./data}"
mkdir -p "$OUTPUT_DIR"

for version in 4 6; do
  filename="ip2region_v${version}.xdb"
  case "$version" in
    4) expected_sha256="59d4ce4ffd9dfcab94ff08dd0fa842e3cc7f19684f16a12396e80a0226a5d3a7" ;;
    6) expected_sha256="5859527ac7a44b9207d08e303b3d10634d01d0a5ced85dc3f17928c40829877c" ;;
  esac
  temporary_file="$OUTPUT_DIR/$filename.tmp"
  trap 'rm -f "$temporary_file"' EXIT
  curl --fail --silent --show-error --location --retry 3 --retry-all-errors \
    --connect-timeout 30 --max-time 1200 \
    -H 'Accept: application/vnd.github.raw+json' \
    "https://api.github.com/repos/lionsoul2014/ip2region/contents/data/$filename?ref=$IP2REGION_REF" \
    --output "$temporary_file"
  printf '%s  %s\n' "$expected_sha256" "$temporary_file" | sha256sum --check --status
  mv "$temporary_file" "$OUTPUT_DIR/$filename"
done
