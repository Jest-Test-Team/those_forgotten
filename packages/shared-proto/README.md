# Shared Proto Contracts

This folder stores the internal Protobuf contracts used by the Rust sidecar services and the Go control plane.

Current services defined in `customs/platform/v1/platform.proto`:

- `MatcherService`
- `PolicyService`
- `PushService`
- `FeedService`

Recommended usage:

- public browser/API traffic stays on REST/JSON
- internal service-to-service traffic uses Protobuf/gRPC
- Kafka topics carry durable async events between crawler, matcher, push, and billing flows
