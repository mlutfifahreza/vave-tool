.PHONY: help build run test proto migrate-up migrate-down seed-products seed-products-py swagger obs-up obs-down obs-logs obs-restart test-tempo

help:
	@echo "Available commands:"
	@echo "  make build       - Build the application"
	@echo "  make run         - Run the application"
	@echo "  make test        - Run tests"
	@echo "  make proto       - Generate protobuf files"
	@echo "  make swagger     - Generate Swagger documentation"
	@echo "  make migrate-up  - Run database migrations up"
	@echo "  make migrate-down - Run database migrations down"
	@echo "  make seed-products - Insert 10 dummy products (dev only)"
	@echo "  make seed-products-py COUNT=X - Insert X products using Python batch insert"
	@echo ""
	@echo "Observability commands:"
	@echo "  make obs-up      - Start observability stack (Prometheus, Loki, Tempo, Grafana)"
	@echo "  make obs-down    - Stop observability stack"
	@echo "  make obs-restart - Restart observability stack"
	@echo "  make obs-logs    - View observability container logs"
	@echo "  make test-tempo  - Test Tempo integration and get a trace ID"

run:
	go run cmd/api/main.go

test:
	go test -v ./...

proto:
	protoc --go_out=. --go_opt=module=github.com/vave-tool \
		--go-grpc_out=. --go-grpc_opt=module=github.com/vave-tool \
		proto/product.proto

swagger:
	swag init -g cmd/api/main.go -o docs

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

obs-up:
	@echo "Starting observability stack..."
	docker-compose up -d
	@echo ""
	@echo "Observability stack is starting..."
	@echo "Grafana:     http://localhost:3000 (admin/admin)"
	@echo "Prometheus:  http://localhost:9090"
	@echo "Tempo:       http://localhost:3200"
	@echo "Loki:        http://localhost:3100"
	@echo ""
	@echo "Waiting for services to be ready..."
	@sleep 5
	@echo "✓ Stack is ready!"

obs-down:
	@echo "Stopping observability stack..."
	docker-compose down
	@echo "✓ Stack stopped"

obs-restart:
	@echo "Restarting observability stack..."
	docker-compose restart
	@echo "✓ Stack restarted"

obs-logs:
	docker-compose logs -f

test-tempo:
	@./script/test_tempo.sh

build:
	@echo "Building application..."
	go build -o bin/api cmd/api/main.go
	@echo "✓ Build complete: bin/api"
