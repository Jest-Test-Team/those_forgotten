from crawler.ingest import build_payload
from crawler.models import NormalizedAuction


def test_build_payload_is_stable() -> None:
    payload = build_payload(
        "fixtures",
        [
            NormalizedAuction(
                announcement_no="TP-001",
                office="臺北關",
                title="沒入數位相機與鏡頭一批",
                category="3C",
                closing_at="2026-03-16T14:00:00+08:00",
                original_link="https://example.com/tp-001",
                warnings=["現狀交付"],
            )
        ],
    )

    assert payload["source"] == "fixtures"
    assert len(payload["checksum"]) == 64
    assert payload["rows"][0]["announcement_no"] == "TP-001"
