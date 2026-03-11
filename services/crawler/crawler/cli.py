from __future__ import annotations

import argparse
import json
from pathlib import Path

from crawler.spiders.offices import all_spiders


def main() -> None:
    parser = argparse.ArgumentParser(description="Run fixture-backed customs crawlers")
    parser.add_argument("--fixtures", default="fixtures", help="Fixture directory")
    args = parser.parse_args()

    fixtures_dir = Path(args.fixtures)
    payload = []
    for spider in all_spiders():
        payload.extend([row.to_dict() for row in spider.parse_fixture(fixtures_dir)])

    print(json.dumps(payload, ensure_ascii=False, indent=2))


if __name__ == "__main__":
    main()
