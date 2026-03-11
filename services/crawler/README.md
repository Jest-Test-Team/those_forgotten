# Crawler Service

Fixture-backed crawler adapters for each customs office.

## Run

```bash
python -m crawler.cli --fixtures fixtures
```

Post normalized fixture output to the API:

```bash
python -m crawler.cli \
  --fixtures fixtures \
  --post \
  --endpoint http://localhost:8080/internal/ingest/auctions \
  --token demo-token
```

The container image also supports env-driven execution via `services/crawler/scripts/run-crawler.sh`.
