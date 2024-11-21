.PHONY: default run build test docs clean

# Variables
APP_NAME=gocrud

# Tasks
default: run

run:
	@go run cmd/main.go
run-with-docs:
	@swag init -d cmd
	@go run cmd/main.go
build:
	@go build -o $(APP_NAME)
test:
	@go test ./ ...
docs:
	@swag init -d cmd
clean:
	@rm -f $(APP_NAME)
	@rm -rf ./docs
