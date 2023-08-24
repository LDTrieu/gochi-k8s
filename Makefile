
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


