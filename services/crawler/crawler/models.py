from dataclasses import asdict, dataclass


@dataclass(slots=True)
class NormalizedAuction:
    announcement_no: str
    office: str
    title: str
    category: str
    closing_at: str
    original_link: str
    warnings: list[str]

    def to_dict(self) -> dict[str, object]:
        return asdict(self)
