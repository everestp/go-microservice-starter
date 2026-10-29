# Runbook

This document is intentionally short and practical.

## Alerts

- error rate above baseline on gateway or downstream services
- p95 latency increasing above the established threshold
- outbox lag or unconsumed message backlog rising across services
- readiness failures for PostgreSQL or NATS

## Common failures

### Gateway returns 503

- verify downstream gRPC health checks
- check service logs for connection or timeout errors
- inspect circuit breaker and retry metrics

### Auth login fails

- verify the JWT key files are present and readable
- check JWKS endpoint availability
- ensure the username/email format matches database constraints

### Outbox backlog

- inspect NATS connectivity and stream health
- confirm the relay loop is running
- check for database locks or publish errors

## Response workflow

1. Identify affected service and request ID.
2. Check metrics and traces.
3. Check dependency health.
4. Review the recent outbox or queue backlog.
5. Apply quick mitigation before deeper debugging.
