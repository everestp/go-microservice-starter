# Auth Service

This service owns user identity, authentication, and token issuance.

## Purpose

- register new users
- authenticate with password hashing
- issue access and refresh tokens
- expose a JWKS endpoint for gateway verification
- maintain refresh token rotation and reuse detection

## Folder map

| Path                                                         | What it is                 |
| ------------------------------------------------------------ | -------------------------- |
| [cmd/auth/main.go](cmd/auth/main.go)                         | service entry point        |
| [internal/config](internal/config)                           | env validation             |
| [internal/domain](internal/domain)                           | business types             |
| [internal/service](internal/service)                         | auth logic                 |
| [internal/transport/grpc](internal/transport/grpc)           | gRPC server                |
| [internal/repository/postgres](internal/repository/postgres) | SQL access                 |
| [migrations](migrations)                                     | database schema migrations |

## Run locally

```bash
cd services/auth-service
go test ./...
```
