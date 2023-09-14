
LOGGER_BINARY=logServiceApp
LISTENER_BINARY=listener
BROKER_BINARY=brokerServiceApp  
AUTH_BINARY=authServiceApp
MAIL_BINARY=mailServiceApp

# AUTH_VERSION=1.0.0
# BROKER_VERSION=1.0.0
# LOGGER_VERSION=1.0.0
# LISTENER_VERSION=1.0.0 management


RABBITMQ_BINARY ?= rabbitmq-binary
## up: starts all containers 
up:
	@echo "Starting docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## down: stop docker compose
down:
	@echo "Stopping docker images..."
	docker-compose down
	@echo "Docker stopped!"

## build core
build_core:  
	@echo "Building core docker images..."
	docker-compose -f docker-compose-core.yml
	@echo "Staring core docker images..."
	docker-compose -f docker-compose-core.yml up -d

build_service: build_broker build_auth build_mail build_logger build_listener
	@echo "Building service docker images..."
	docker-compose -f docker-compose.yml build --no-cache
	@echo "Staring service docker images..."
	docker-compose -f docker-compose.yml up -d
	
## up_build stop docker-compose, builds all services and starts docker compose
up_build: build_auth build_broker build_listener  build_logger build_mail
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

## build_auth: builds the auth binary as a linux executable
build_auth:
	@echo "Building authentication binary..."
	cd authentication-service && env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${AUTH_BINARY} ./cmd/api
	@echo "Authentication binary built!"

## build_broker: builds the broker binary as a linux executable
build_broker:
	@echo "Building broker binary..."
	cd broker-service && env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd/api
	@echo "Broker binary built!"

## build_listener: builds the listener binary as a linux executable
build_listener:
	@echo "Building listener binary..."
	cd listener-service && env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${LISTENER_BINARY} .
	@echo "Listener binary built!"

## build_logger: builds the logger binary as a linux executable
build_logger:
	@echo "Building logger binary..."
	cd logger-service && env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${LOGGER_BINARY} ./cmd/web
	@echo "Logger binary built!"


## build_mail: builds the mail binary as a linux executable
build_mail:
	@echo "Building mail binary..."
	cd mail-service && env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${MAIL_BINARY} ./cmd/api
	@echo "Mail binary built!"
	
## logger: stops logger-service, removes docker image, builds service, and starts it
logger: build_logger
	@echo "Building logger-service docker image..."
	- docker-compose stop logger-service
	- docker-compose rm -f logger-service
	docker-compose up --build -d logger-service
	docker-compose start logger-service
	@echo "broker-service rebuilt and started!"

## listener: stops listener-service, removes docker image, builds service, and starts it
listener: build_listener
	@echo "Building listener-service docker image..."
	- docker-compose stop listener-service
	- docker-compose rm -f listener-service
	docker-compose up --build -d listener-service
	docker-compose start listener-service
	@echo "listener-service rebuilt and started!"

run_rabbitmqmanagement:
	docker build --tag ${RABBITMQ_BINARY} .
	
run_auth:
	@$(MAKE) -C authentication-service/ run
	@echo "Running authentication-service..."
	
run_broker:
	@$(MAKE) -C broker-service/ run
	@echo "Running broker-service..."
	
run_frontend:
	@$(MAKE) -C front-end/ run
	@echo "Running frontend-service..."
	
run_logger:
	@$(MAKE) -C logger-service/ run
	@echo "Running logger-service..."

run_mail:
	@$(MAKE) -C mail-service/ run
	@echo "Running mail-service..."

prune:
	docker system prune -all