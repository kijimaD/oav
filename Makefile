.DEFAULT_GOAL := help
DOCKER_TAG := latest

.PHONY: build
build: ## Build image for deploy
	docker build -t kijimad/oav:${DOCKER_TAG} \
	--target deploy ./

.PHONY: build-local
build-local: ## Build image for local development
	docker-compose build --no-cache

.PHONY: up
up: ## Do docker compose up
	docker-compose up -d

.PHONY: down
down: ## Do docker compose down
	docker-compose down

.PHONY: logs
logs: ## Tail docker compose logs
	docker-compose logs -f

.PHONY: ps
ps: ## Check container status
	docker-compose ps

.PHONY: lint
lint: ## Run lint
	docker run --rm -v ${PWD}:/app -w /app golangci/golangci-lint:v1.50.1 golangci-lint run -v

.PHONY: test
test: ## Run test
	go test -race -shuffle=on -v ./...

.PHONY: help
help: ## Show help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: o/lint
o/lint: spectral yamllint openapi-validator ## すべてのLintを実行する

.PHONY: spectral
spectral: ## spectral lintを実行する
	docker run --rm -v ${PWD}:/work -w /work/docs stoplight/spectral lint openapi.yml

.PHONY: yamllint
yamllint: ## yaml lintを実行する
	docker run --rm -v ${PWD}:/data cytopia/yamllint -s openapi.yml

.PHONY: openapi-validator
openapi-validator: ## openapi-validatorを実行する
	docker run --rm -v ${PWD}:/work jamescooke/openapi-validator --verbose --report_statistics /work/openapi.yml
