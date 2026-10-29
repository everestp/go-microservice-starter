# Architecture

This repository intentionally uses a small but explicit system topology:

- gateway serves the public API and protects the system boundary
- auth-service owns identities and token issuance
- notes-service owns note data and access rules
- audit-service records the event trail from the messaging bus

The shared packages are intentionally limited to cross-cutting concerns and do not implement business behavior.

## Core design principles

- one database per service
- one service owns one business capability
- synchronous RPC for workflows that need immediate response
- asynchronous eventing for side effects and audit records
- explicit contracts through protobuf files
- environment-based config and fail-fast validation

## Reliability choices

The starter does not attempt to over-engineer a saga framework or a distributed transaction. It uses an outbox pattern so that the database and the event stream are consistent at the service boundary while keeping the design understandable for beginners.
