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

lint:
	golangci-lint run ./...
