# Gateway

The gateway is the public API entrypoint for the system. It stands in front of the internal microservices and shapes the external HTTP experience.

## What this service does

This service is responsible for:

- receiving client HTTP requests
- routing requests to the appropriate internal service
- handling public API versioning
- validating or forwarding identity information
- collecting request metadata such as trace IDs and request IDs
- exposing a single public network boundary for the platform

## Why it matters

The gateway centralizes access to the platform and helps keep internal services clean. Clients only talk to the gateway, while the backend services stay focused on their business domains.

## Main responsibilities by layer

### Config layer

This loads environment values such as ports and upstream service addresses.

### Entry point

The startup file bootstraps the service and logs the listening addresses.

### Future evolution

This starter is intentionally minimal. A production gateway would normally add:

- JWT validation
- rate limiting
- request validation
- OpenAPI documentation
- middleware for auth, tracing, and logging
- upstream gRPC client calls

## Important files

- cmd/gateway/main.go: service entrypoint
- internal/config/config.go: config values and defaults

## Run locally

```bash
cd services/gateway
go test ./...
```

## Deployment note

The gateway is the place where the platform meets the outside world, so it is usually the first service to harden in real production environments.
