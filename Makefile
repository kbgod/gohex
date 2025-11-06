.PHONY: docs

docs:
	swag fmt && swag init -g ./cmd/api/main.go

mocks:
	go generate ./internal/...