.PHONY: build run test clean docker-up docker-down

# Build the application
build:
	go build -o bin/server cmd/server/main.go

# Run the application
run:
	go run cmd/server/main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Start Docker services
docker-up:
	docker-compose up --build -d

# Stop Docker services
docker-down:
	docker-compose down

# View logs
logs:
	docker-compose logs -f app

# Database migration (if needed)
migrate-up:
	migrate -path migrations -database "postgres://user:password@localhost:5432/salesdb?sslmode=disable" up

# Load sample data
load-data:
	curl -X POST "http://localhost:8080/api/v1/refresh?file_path=data/sales_data.csv"