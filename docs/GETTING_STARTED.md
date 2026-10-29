# Getting Started

This project is designed to be simple enough to run locally with a small set of tools.

## Prerequisites

- Go 1.26 or later
- Docker and Docker Compose
- Buf CLI
- golangci-lint

## First hour

1. Run `go work sync` from the repository root.
2. Run `make proto` to generate the gRPC stubs.
3. Run `make up` to bring up Postgres, NATS, and the observability stack.
4. Hit the health endpoint with `curl http://localhost:8080/healthz`.
5. Register a user via the auth service or gateway endpoint.
6. Create a note and inspect the outbox and audit flow.
7. Open Jaeger or Grafana to observe spans and metrics.

## Common workstation setup

```bash
go install github.com/bufbuild/buf/cmd/buf@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

After installation, make sure the Go bin directory is on your PATH.

## Local troubleshooting

- If the database is unhealthy, run `docker compose up postgres` and check the logs.
- If the gateway cannot reach the downstream services, verify the container names and addresses in the Compose file.
- If Buf generation fails, ensure the protobuf files are valid and the plugin output paths exist.
