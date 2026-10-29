# Notes Service

This service owns note storage and the authorization rules for note access.

## Purpose

- create notes
- list notes by owner
- update with optimistic concurrency
- delete notes by ownership
- emit outbox events for downstream audit processing

## Run locally

```bash
cd services/notes-service
go test ./...
```
