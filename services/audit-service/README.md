# Audit Service

The audit service is the accountability layer of the platform. It records important system events so they can be reviewed later.

## What this service does

This service is responsible for:

- receiving event data from other services
- storing audit records in an append-only style
- keeping a timeline of important actions
- exposing the event list for inspection or support workflows

## Why it exists

A microservice platform needs strong observability and traceability. Even when the primary business logic is in other services, the audit service ensures that you can answer questions like:

- who did what
- when did it happen
- which service produced the event
- what metadata was associated with it

## Main responsibilities by layer

### Domain layer

The domain layer defines the audit event model and the event metadata structure.

### Service layer

The service layer provides:

- Record: write a new audit event
- List: read back the event history
- helper methods for marshaling payload data and validating timestamps

### Repository layer

The repository stores events in memory for the scaffold. In production, this would be backed by a dedicated audit table or append-only storage.

### Transport layer

The gRPC adapter lets the gateway or other services request event history through a typed API.

## Important files

- cmd/audit/main.go: service entrypoint
- internal/config/config.go: env configuration
- internal/domain/event.go: event model
- internal/service/audit.go: service logic
- internal/service/memory.go: in-memory event persistence
- internal/transport/grpc/server.go: gRPC exposure

## Run locally

```bash
cd services/audit-service
go test ./...
```

## Production direction

The starter version is intentionally simple. A full production version should include:

- durable database storage
- deduplication and idempotency checks
- NATS consumer integration
- strict retention policies
- stronger event schema validation
