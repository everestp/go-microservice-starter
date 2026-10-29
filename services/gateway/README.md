# Gateway

The gateway is the public edge. It handles HTTP, API versioning, JWT verification, and downstream calls.

## Purpose

- validate JWTs from the auth service
- rate limit public requests
- map the REST surface to gRPC calls
- aggregate user-facing responses
- forward request IDs and tracing metadata

## Run locally

```bash
cd services/gateway
go test ./...
```
