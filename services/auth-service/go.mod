module github.com/example/go-microservice-starter/services/auth-service

go 1.26.8

require (
	github.com/example/go-microservice-starter v0.0.0
	github.com/example/go-microservice-starter/pkg v0.0.0
	github.com/google/uuid v1.6.0
	golang.org/x/crypto v0.41.0
)

replace github.com/example/go-microservice-starter => ../..
replace github.com/example/go-microservice-starter/pkg => ../../pkg
