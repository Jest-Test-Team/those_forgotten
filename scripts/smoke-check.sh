#!/usr/bin/env bash

set -euo pipefail

WEB_BASE_URL="${WEB_BASE_URL:-}"
API_BASE_URL="${API_BASE_URL:-}"

if [[ -z "${API_BASE_URL}" ]]; then
  echo "API_BASE_URL is required" >&2
  exit 1
fi

curl --fail --silent --show-error "${API_BASE_URL}/healthz" >/dev/null
curl --fail --silent --show-error "${API_BASE_URL}/readyz" >/dev/null

if [[ -n "${WEB_BASE_URL}" ]]; then
  curl --fail --silent --show-error "${WEB_BASE_URL}/" >/dev/null
  curl --fail --silent --show-error "${WEB_BASE_URL}/sitemap.xml" >/dev/null
fi

echo "smoke check passed"
