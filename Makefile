.PHONY: help build run test proto migrate-up migrate-down

help:
	@echo "Available commands:"
	@echo "  make build       - Build the application"
	@echo "  make run         - Run the application"
	@echo "  make test        - Run tests"
	@echo "  make proto       - Generate protobuf files"
	@echo "  make migrate-up  - Run database migrations up"
	@echo "  make migrate-down - Run database migrations down"

build:
	go build -o bin/api cmd/api/main.go

run:
	go run cmd/api/main.go

test:
	go test -v ./...

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/product.proto

migrate-up:
	@echo "Running migrations..."
	@for file in migrations/*up.sql; do \
		echo "Applying $$file"; \
		psql -h localhost -U postgres -d vave_db -f $$file; \
	done

migrate-down:
	@echo "Rolling back migrations..."
	@for file in $$(ls -r migrations/*down.sql); do \
		echo "Applying $$file"; \
		psql -h localhost -U postgres -d vave_db -f $$file; \
	done

deps:
	go mod download
	go mod tidy

clean:
	rm -rf bin/
