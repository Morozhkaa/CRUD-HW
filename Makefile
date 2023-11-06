.PHONY: help
help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

run: ##  Run docker-compose
	docker-compose up --build
.PHONY: run

remove-volumes: ## Down docker-compose
	docker-compose down --volumes
.PHONY: remove-volumes

swag: ## generate swagger docs
	swag init --dir ./internal/adapters/http --generalInfo swagger.go --output ./api/swagger/public --parseDepth 1 --parseDependency
