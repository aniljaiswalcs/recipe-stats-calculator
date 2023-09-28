export GO=$(shell which go)
export DOCKER=$(shell which docker)
export DOCKER_CONTAINER_LIST=$(shell docker ps -aq)
export DOCKER_IMAGE_LIST=$(shell docker image ls -aq)

APP_NAME=recipe-stats-calculator

default:cli_build

cli_build:
	@go build -o bin/recipe-stats-calculator main.go
	
cli_run:
	@go run main.go

.PHONY: test
test:
	@go test -coverprofile=coverage.out ./...


.PHONY: build
build:
	@$(DOCKER) build -t $(APP_NAME) -f ./Dockerfile .

run:
	@$(DOCKER) run $(APP_NAME)

docker_clean:
    
	@if [ -n "$(DOCKER_CONTAINER_LIST)" ]; then echo "Removing docker containers" && docker rm $(DOCKER_CONTAINER_LIST); else echo "No containers found"; fi;
	@if [ -n "$(DOCKER_IMAGE_LIST)" ]; then echo "Removing docker images" && docker rmi $(DOCKER_IMAGE_LIST); else echo "No image found"; fi;
	

.PHONY: clean
clean:
	rm -f ./bin/*
