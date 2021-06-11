VERSION=$(shell cat VERSION)
export BASE_BINARY_NAME=terraform-provider-looker_v$(VERSION)

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: vendor
vendor: go.sum ## vendor dependencies
	@GO111MODULE=on go mod vendor
	@GO111MODULE=on go mod tidy

.PHONY: lint
lint: ## run linter
	@golangci-lint run ./...

.PHONY: test
test: ## run tests
	@go test -v -cover -race $(shell go list ./... | grep -v vendor)

.PHONY: test-acceptance
test-acceptance: ## runs all tests, including the acceptance tests
	@TF_ACC=1 $(go_test) go test  -v -cover $(shell go list ./... | grep -v vendor)

.PHONY: build
build: ## build binary
	@go build -o build/$(BASE_BINARY_NAME) .

.PHONY: docs
docs: ## generate docs
	@go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

.PHONY: check-docs
check-docs: docs ## check that docs have been generated
	@git diff --exit-code -- docs

.PHONY: check-mod
check-mod: ## check go.mod is up-to-date
	@go mod tidy
	@git diff --exit-code -- go.mod go.sum
