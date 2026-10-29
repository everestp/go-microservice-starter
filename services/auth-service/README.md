# Auth Service

The auth service is the identity and session service for the platform. It owns user accounts, password verification, and token issuance.

## What this service does

This service is responsible for:

- registering a new user
- checking whether an email is already in use
- hashing passwords before storing them
- verifying a submitted password against the stored hash
- issuing access tokens and refresh tokens
- rotating refresh tokens for security
- revoking refresh token families on logout or refresh
- returning the current user profile for authenticated requests

## Main responsibilities by layer

### Domain layer

The domain package stores the core model types:

- user record
- refresh token record
- token pair response
- validation errors such as already exists, invalid credentials, bad refresh token, and revoked refresh token

### Service layer

The service layer holds the business logic:

- Register
- Login
- Refresh
- Logout
- Me

This is where the rules live for user creation and refresh token rotation.

### Repository layer

The repository layer persists user and refresh token data. In the starter, the repository is represented as an in-memory repository for local testing and local development.

In a real deployment, this would be replaced by a PostgreSQL-backed implementation.

### Transport layer

The transport layer adapts the domain service to gRPC.

It converts the service calls into protobuf messages and maps domain errors to gRPC status codes such as:

- AlreadyExists
- NotFound
- Unauthenticated
- Internal

## Important files

- cmd/auth/main.go: starts the service
- internal/config/config.go: loads runtime config and defaults
- internal/domain/user.go: core auth domain models and error definitions
- internal/service/auth.go: actual business logic
- internal/service/crypto.go: hashing and verification helpers
- internal/service/service_test.go: unit tests for auth behavior
- internal/repository/postgres/memory.go: in-memory persistence implementation for the scaffold
- internal/transport/grpc/server.go: gRPC adapter

## Typical request flow

When a client registers:

1. the gateway calls the auth gRPC service
2. auth-service checks if the user already exists
3. it hashes the password
4. it creates the user record
5. it issues an access token and refresh token
6. it stores the refresh token hash and returns the pair to the caller

## Run locally

```bash
cd services/auth-service
go test ./...
```

## Notes

This is a good starter implementation for understanding service ownership and token flows, but it is not yet a full production identity system. A production version would also include:

- real PostgreSQL persistence
- proper JWT signing and JWKS distribution
- full mTLS and service identity
- audit and event publishing
- stricter refresh-token security checks
