#!/bin/sh
set -eu

fixtures="${FIXTURES:-fixtures}"
set -- python -m crawler.cli --fixtures "$fixtures"

if [ "${POST_TO_API:-false}" = "true" ]; then
  endpoint="${INGEST_ENDPOINT:-http://api:8080/internal/ingest/auctions}"
  token="${INTERNAL_INGEST_TOKEN:-}"
  set -- "$@" --post --endpoint "$endpoint" --token "$token"
fi

exec "$@"
