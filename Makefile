
LOGGER_BINARY=logServiceApp
LOGGER_VERSION=1.0.0

up:
	@echo "Starting docker images..."
	docker-compose up -d
	@echo "Docker images started!"

down:
	ls

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


## logger: stops logger-service, removes docker image, builds service, and starts it
logger: build_logger
	@echo "Building logger-service docker image..."
	- docker-compose stop logger-service
	- docker-compose rm -f logger-service
	docker-compose up --build -d logger-service
	docker-compose start logger-service
	@echo "broker-service rebuilt and started!"
