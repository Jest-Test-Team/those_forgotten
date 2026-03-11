from __future__ import annotations

import hashlib
import json

import httpx

from crawler.models import NormalizedAuction


def build_payload(source: str, rows: list[NormalizedAuction]) -> dict[str, object]:
    serialized_rows = [row.to_dict() for row in rows]
    checksum = hashlib.sha256(
        json.dumps(serialized_rows, ensure_ascii=False, sort_keys=True).encode("utf-8")
    ).hexdigest()
    return {
        "source": source,
        "checksum": checksum,
        "rows": serialized_rows,
    }


def send_payload(
    *,
    endpoint: str,
    ingest_token: str,
    payload: dict[str, object],
    timeout: float = 10.0,
) -> httpx.Response:
    return httpx.post(
        endpoint,
        json=payload,
        headers={"X-Ingest-Token": ingest_token},
        timeout=timeout,
    )
