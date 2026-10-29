# Audit Service

This service consumes business events and writes an immutable audit trail.

## Purpose

- ingest NATS JetStream events
- dedupe via processed_events table
- persist audit records in a separate database
- expose an audit read API to the gateway

## Run locally

```bash
cd services/audit-service
go test ./...
```
