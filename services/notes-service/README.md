# Notes Service

The notes service owns the note domain. It handles the lifecycle of user notes and keeps the business rules close to the data it modifies.

## What this service does

This service is responsible for:

- creating notes
- listing notes
- fetching a single note by ID
- updating note contents
- deleting notes
- protecting update operations with version checks

## Why version checks matter

This service uses optimistic concurrency. Each note has a version number, and updates only succeed when the caller provides the current version. This prevents stale writes from overwriting newer data.

## Main responsibilities by layer

### Domain layer

The note domain model stores the note record and relevant error states like missing note or version mismatch.

### Service layer

The service layer implements the core CRUD behavior:

- Create
- Get
- List
- Update
- Delete

It also sorts and paginates response lists in a simple local pattern.

### Repository layer

The repository stores notes in memory during the starter phase. In a real deployment, this would be backed by a notes table in PostgreSQL.

### Transport layer

The gRPC transport layer exposes the notes contract so the gateway or other services can call the note APIs in a strongly typed way.

## Important files

- cmd/notes/main.go: service entrypoint
- internal/config/config.go: runtime defaults and environment setup
- internal/domain/note.go: note model and domain errors
- internal/service/notes.go: note business logic
- internal/service/memory.go: in-memory persistence for local development
- internal/transport/grpc/server.go: gRPC adapter

## Run locally

```bash
cd services/notes-service
go test ./...
```

## Future production upgrades

This starter version is intentionally small. A production version would add:

- PostgreSQL-backed storage
- note ownership and authorization checks
- pagination and filtering on larger datasets
- event publishing to NATS when notes change
- better audit and domain validation
