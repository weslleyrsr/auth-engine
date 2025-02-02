# Variables
COMPOSE_FILE=docker-compose.yml

# Targets
.PHONY: up down restart logs clean build test-account

# Start the app
up:
	@echo "Starting the app..."
	docker-compose -f $(COMPOSE_FILE) up -d

# Stop the app
down:
	@echo "Stopping the app..."
	docker-compose -f $(COMPOSE_FILE) down

# Restart the app
restart: down up

# View logs
logs:
	@echo "Viewing logs..."
	docker-compose -f $(COMPOSE_FILE) logs -f

# Clean up volumes and networks
clean:
	@echo "Cleaning up..."
	docker-compose -f $(COMPOSE_FILE) down --volumes --remove-orphans

# Build the app
build:
	@echo "Building the app..."
	docker-compose -f $(COMPOSE_FILE) up --build -d

# Test account
test-account:
	@echo "Testing Account..."
	cd account && go test ./...