.DEFAULT_GOAL := help

.PHONY: test-data
test-data: export DATA ?= ColorGroup
test-data: ## Run data test
	@go test -timeout 30s -v -run ^Test_Data$

.PHONY: test
test: export DATA ?= ColorGroup
test: ## Run unit tests
	@go test -race -cover ./...

.PHONY: coverage
coverage: ## Run unit tests with coverage
	@go test -v -race -cover -coverpkg=./... -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -func=coverage.out

.PHONY: html
html: ## Show html coverage result
	@go tool cover -html=./coverage.out

.PHONY: html-gen
html-gen: ## Export html coverage result
	@go tool cover -html=./coverage.out -o ./coverage.html

.PHONY: html-wsl
html-wsl: html-gen ## Open html coverage result in wsl
	@explorer.exe `wslpath -w ./coverage.html` || true

.PHONY: help
help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
