BINARY = pr_review_service

build:
	@echo "Building go binary"
	go build -o $(BINARY) ./cmd/service

run:
	@echo "Running Go service"
	go run ./cmd/service/main.go

docker-build:
	@echo "Building Docker image"
	docker-compose build

up:
	@echo "Starting docker-compose"
	docker-compose up

down:
	@echo "Stopping docker-compose"
	docker-compose down

logs:
	docker-compose logs -f app

test:
	go test ./... -v

lint:
	golangci-lint run
