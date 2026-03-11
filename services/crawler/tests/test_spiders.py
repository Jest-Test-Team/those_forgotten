from pathlib import Path

from crawler.spiders.offices import KaohsiungSpider, TaipeiSpider, all_spiders


FIXTURES = Path(__file__).resolve().parents[1] / "fixtures"


def test_all_spiders_return_rows() -> None:
    rows = []
    for spider in all_spiders():
      rows.extend(spider.parse_fixture(FIXTURES))

    assert len(rows) == 4
    assert {row.office for row in rows} == {"基隆關", "臺北關", "臺中關", "高雄關"}


def test_warning_derivation() -> None:
    row = KaohsiungSpider().parse_fixture(FIXTURES)[0]
    assert "需確認菸酒相關資格" in row.warnings


def test_taipei_source_url_from_link() -> None:
    row = TaipeiSpider().parse_fixture(FIXTURES)[0]
    assert row.original_link == "https://example.com/tp-001"
