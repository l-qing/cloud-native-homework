.PHONY: default
default: help

##@ Development

.PHONY: run test
run: ## Run directly
	export VERSION=v1.2.3
	go run ./

test: ## Test use curl
	curl -v http://localhost:8081/health
	curl -v -H "Custom-Header-Key: custom-value" \
		-H 'Current-Header-Array: [ "value1", "value2", "value3" ]' \
		-H 'Current-Header-Map: {"key", "value"}' \
		http://localhost:8081/abc

build: ## Build docker image
	docker build --build-arg version=v2.1.0 . -t httpserver:latest

push: ## Push docker image
	docker push httpserver:latest

run-docker: ## Run docker image
	docker run -it --name=httpserver --rm -p 8081:8081 httpserver:latest

##@ General

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

