# Robot Framework Live Tests

This folder contains Robot Framework suites for live smoke and regression testing against the hosted API and web frontend.

## Covered Surfaces

- API root, health, readiness, Swagger UI, Swagger spec
- public API endpoints
- protected API endpoints without auth
- public web routes
- protected web routes regression checks

## Install

```bash
python3 -m venv test/.venv
source test/.venv/bin/activate
pip install -r test/requirements.txt
```

## Run Against Current Hosted URLs

```bash
source test/.venv/bin/activate
robot -d test/results test
```

## Override Target URLs

```bash
export ROBOT_API_BASE_URL=https://comfortable-adelind-dennis-team-1d8552ab.koyeb.app
export ROBOT_WEB_BASE_URL=https://those-forgotten.vercel.app
robot -d test/results test
```

## Notes

- The suites intentionally treat `/member` and `/admin` as regression checks. They should not return `500`.
- If those routes are still broken in the hosted frontend, the protected web suite will fail and surface the deployment issue.
