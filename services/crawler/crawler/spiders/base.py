from __future__ import annotations

from abc import ABC, abstractmethod
from pathlib import Path

from bs4 import BeautifulSoup

from crawler.models import NormalizedAuction


class OfficeSpider(ABC):
    office: str
    fixture_name: str
    source_url: str

    def parse_fixture(self, fixtures_dir: Path) -> list[NormalizedAuction]:
        html = (fixtures_dir / self.fixture_name).read_text(encoding="utf-8")
        return self.parse_html(html)

    def parse_html(self, html: str) -> list[NormalizedAuction]:
        soup = BeautifulSoup(html, "html.parser")
        return self.extract_rows(soup)

    @abstractmethod
    def extract_rows(self, soup: BeautifulSoup) -> list[NormalizedAuction]:
        raise NotImplementedError
