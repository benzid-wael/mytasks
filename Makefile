export GOTESTSUM_FORMAT := $(or ${GOTESTSUM_FORMAT},${GOTESTSUM_FORMAT},short-verbose)

.PHONY: test-unit
test-unit: ## run unit tests, to change the output format use: GOTESTSUM_FORMAT=(dots|short|standard-quiet|short-verbose|standard-verbose) make test-unit
	gotestsum $(TESTFLAGS) -- $${TESTDIRS:-$(shell go list ./... | grep -vE '/vendor/|/e2e/')}

.PHONY: test
test: test-unit ## run tests

.PHONY: test-coverage
test-coverage: ## run test coverage
	gotestsum -- -coverprofile=coverage.txt $(shell go list ./... | grep -vE '/vendor/|/e2e/')

.PHONY: fmt
fmt:
	go list -f {{.Dir}} ./... | xargs gofmt -w -s -d

.PHONY: lint
lint: ## run all the lint tools
	# gometalinter --config gometalinter.json ./...
	golangci-lint run

.PHONY: help
help: ## print this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / {gsub("\\\\n",sprintf("\n%22c",""), $$2);printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
