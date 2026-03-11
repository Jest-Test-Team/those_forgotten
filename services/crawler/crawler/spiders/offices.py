from __future__ import annotations

from bs4 import BeautifulSoup

from crawler.models import NormalizedAuction
from crawler.spiders.base import OfficeSpider


class TableSpider(OfficeSpider):
    def extract_rows(self, soup: BeautifulSoup) -> list[NormalizedAuction]:
        rows = []
        for row in soup.select("table tbody tr"):
            columns = row.find_all("td")
            if len(columns) < 5:
                continue

            link = columns[1].find("a")
            rows.append(
                NormalizedAuction(
                    announcement_no=columns[0].get_text(strip=True),
                    office=self.office,
                    title=columns[1].get_text(strip=True),
                    category=columns[2].get_text(strip=True),
                    closing_at=columns[3].get_text(strip=True),
                    original_link=(link["href"] if link and link.has_attr("href") else self.source_url),
                    warnings=_derive_warnings(columns[4].get_text(strip=True)),
                )
            )
        return rows


def _derive_warnings(text: str) -> list[str]:
    warnings = ["現狀交付", "不負瑕疵擔保"]
    keywords = {
        "菸酒": "需確認菸酒相關資格",
        "進口車": "需確認驗車與牌照流程",
        "醫療": "需確認輸入許可與法規限制",
    }
    for key, warning in keywords.items():
        if key in text:
            warnings.append(warning)
    return warnings


class KeelungSpider(TableSpider):
    office = "基隆關"
    fixture_name = "keelung.html"
    source_url = "https://web.customs.gov.tw/keelung/htmlList/216"


class TaipeiSpider(TableSpider):
    office = "臺北關"
    fixture_name = "taipei.html"
    source_url = "https://web.customs.gov.tw/taipei/htmlList/86208152b05545bdad39749ea730870d"


class TaichungSpider(TableSpider):
    office = "臺中關"
    fixture_name = "taichung.html"
    source_url = "https://web.customs.gov.tw/taichung"


class KaohsiungSpider(TableSpider):
    office = "高雄關"
    fixture_name = "kaohsiung.html"
    source_url = "https://web.customs.gov.tw/kaohsiung/htmlList/127d64b348734567b41bdc83590b9d12"


def all_spiders() -> list[OfficeSpider]:
    return [KeelungSpider(), TaipeiSpider(), TaichungSpider(), KaohsiungSpider()]
