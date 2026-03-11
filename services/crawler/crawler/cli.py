from __future__ import annotations

import argparse
import json
import os
from pathlib import Path

from crawler.ingest import build_payload, send_payload
from crawler.spiders.offices import all_spiders


def main() -> None:
    parser = argparse.ArgumentParser(description="Run fixture-backed customs crawlers")
    parser.add_argument("--fixtures", default="fixtures", help="Fixture directory")
    parser.add_argument("--post", action="store_true", help="POST the normalized payload to the API")
    parser.add_argument(
        "--endpoint",
        default=os.getenv("INGEST_ENDPOINT", "http://localhost:8080/internal/ingest/auctions"),
        help="API endpoint for normalized ingestion",
    )
    parser.add_argument(
        "--token",
        default=os.getenv("INTERNAL_INGEST_TOKEN", ""),
        help="Ingestion token for the API",
    )
    args = parser.parse_args()

    fixtures_dir = Path(args.fixtures)
    rows = []
    for spider in all_spiders():
        rows.extend(spider.parse_fixture(fixtures_dir))

    payload = build_payload("fixtures", rows)

    if args.post:
        response = send_payload(endpoint=args.endpoint, ingest_token=args.token, payload=payload)
        print(json.dumps(response.json(), ensure_ascii=False, indent=2))
        return

    print(json.dumps(payload, ensure_ascii=False, indent=2))


if __name__ == "__main__":
    main()
