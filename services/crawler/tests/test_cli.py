from __future__ import annotations

import json
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
import os
from pathlib import Path
import subprocess
import threading


ROOT = Path(__file__).resolve().parents[1]
FIXTURES = ROOT / "fixtures"


def test_cli_outputs_fixture_payload() -> None:
    result = subprocess.run(
        [
            str(ROOT / ".venv" / "bin" / "python"),
            "-m",
            "crawler.cli",
            "--fixtures",
            str(FIXTURES),
        ],
        cwd=ROOT,
        check=True,
        capture_output=True,
        text=True,
    )

    payload = json.loads(result.stdout)
    assert payload["source"] == "fixtures"
    assert len(payload["rows"]) == 4


def test_cli_posts_payload_to_api() -> None:
    captured: list[dict[str, object]] = []

    class Handler(BaseHTTPRequestHandler):
        def do_POST(self) -> None:  # noqa: N802
            length = int(self.headers["Content-Length"])
            body = self.rfile.read(length)
            captured.append(
                {
                    "path": self.path,
                    "headers": {
                        "X-Ingest-Token": self.headers.get("X-Ingest-Token", ""),
                    },
                    "body": json.loads(body.decode("utf-8")),
                }
            )
            response = json.dumps({"accepted": True}).encode("utf-8")
            self.send_response(202)
            self.send_header("Content-Type", "application/json")
            self.send_header("Content-Length", str(len(response)))
            self.end_headers()
            self.wfile.write(response)

        def log_message(self, format: str, *args: object) -> None:  # noqa: A003
            return

    server = ThreadingHTTPServer(("127.0.0.1", 0), Handler)
    thread = threading.Thread(target=server.serve_forever, daemon=True)
    thread.start()
    try:
        result = subprocess.run(
            [
                str(ROOT / ".venv" / "bin" / "python"),
                "-m",
                "crawler.cli",
                "--fixtures",
                str(FIXTURES),
                "--post",
                "--endpoint",
                f"http://127.0.0.1:{server.server_address[1]}/internal/ingest/auctions",
                "--token",
                "demo-token",
            ],
            cwd=ROOT,
            check=True,
            capture_output=True,
            text=True,
        )
    finally:
        server.shutdown()
        server.server_close()
        thread.join(timeout=2)

    payload = json.loads(result.stdout)
    assert payload == {"accepted": True}
    assert captured[0]["path"] == "/internal/ingest/auctions"
    assert captured[0]["headers"]["X-Ingest-Token"] == "demo-token"
    assert captured[0]["body"]["source"] == "fixtures"


def test_cli_prefers_ingest_url_env_alias() -> None:
    result = subprocess.run(
        [
            str(ROOT / ".venv" / "bin" / "python"),
            "-c",
            "from crawler.cli import default_ingest_endpoint; print(default_ingest_endpoint())",
        ],
        cwd=ROOT,
        check=True,
        capture_output=True,
        text=True,
        env={
            **os.environ,
            "INGEST_URL": "https://api.example.com/internal/ingest/auctions",
            "INGEST_ENDPOINT": "http://should-not-win.example.com",
        },
    )

    assert result.stdout.strip() == "https://api.example.com/internal/ingest/auctions"
