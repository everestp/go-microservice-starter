# Architecture Decision Records

This project follows a lightweight ADR format for each major system decision.

## ADR-001: Microservices vs. monolith

### Context

The business domain is small but the system should demonstrate service boundaries with independent deployability.

### Decision

The repository uses separate Go modules and service folders while keeping the business domain intentionally narrow.

### Consequences

Pros: clear team boundaries, per-service deployment, explicit contracts.
Cons: more operational work, more moving parts, more observability complexity.

## ADR-002: gRPC inside, HTTP at the edge

### Context

The public API is versioned and the internal network should be efficient and contract-first.

### Decision

The gateway exposes HTTP and calls downstream services with gRPC.

### Consequences

This keeps the public edge easy to version while using a strongly typed internal contract.

## ADR-003: NATS JetStream for async events

### Context

The system needs durable event processing, retries, and replay support.

### Decision

Use NATS JetStream as the event bus with durable consumers and dead-letter semantics.

### Consequences

This is a good fit for a small distributed starter and keeps the operational footprint lighter than Kafka.

## ADR-004: Database per service

### Context

A shared database would create hidden coupling and service dependency leaks.

### Decision

Each service owns its own PostgreSQL database and repository layer.

### Consequences

Services can scale and deploy independently but cross-service reporting must be done via events or API aggregation.

## ADR-005: Transactional outbox

### Context

A business write and an external event must not be published in the same transaction.

### Decision

Write the business row and outbox row in one transaction, then publish asynchronously.

### Consequences

This gives a practical at-least-once delivery model at the cost of some operational complexity and replay handling.

## ADR-006: JWT asymmetric keys + JWKS

### Context

The system needs constrained token verification and rotation.

### Decision

The auth service signs access tokens with an asymmetric key pair and exposes the public keys via JWKS.

### Consequences

The gateway can verify tokens without sharing private keys while supporting key rotation.

## ADR-007: Service-to-service auth

### Context

Internal gRPC traffic must not blindly trust the network.

### Decision

The system uses mutual TLS and identity propagation in metadata for service verification.

### Consequences

This is safer than raw network trust but requires certificate management and validation hooks.

## ADR-008: Pagination strategy

### Context

Collections must be predictable and safe under load.

### Decision

Use keyset pagination with a whitelist of sort fields and explicit max page size limits.

### Consequences

This reduces DB load and avoids offset performance cliffs at scale.

## ADR-009: Monorepo with go.work

### Context

The project contains multiple services and shared packages.

### Decision

Use a root-level Go workspace and per-service modules.

### Consequences

This keeps local development coherent while preserving each service's own module boundary.
