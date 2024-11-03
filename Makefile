# Build the application
.PHONY: build
build:
	go build  main.go
	
# Clean up build artifacts
.PHONY: clean
clean:
	@echo "Cleaning up..."
	rm -rf ./build

# Run a the service
.PHONY: run
run:
	./main

up:
	docker-compose up
down:
	docker-compose down

# Print help message
.PHONY: help
help:
	@echo "Makefile for Go Microservices"
	@echo "Usage:"
	@echo "  make build            - Build the application"
	@echo "  make clean            - Remove build artifacts"
	@echo "  make run        	   - run the service"
	@echo "  make up      		   - start the containers required for the service"
	@echo "  make down      	   - stop the containers required for the service"