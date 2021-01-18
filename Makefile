.ONESHELL:
GOFILES := $(shell find . -name '*.go' -not -path "./vendor/*" | egrep -v "^\./\.go" | grep -v _test.go)
DATE = $(shell date +'%s')

.PHONY: build lint

test:
	go test  ./... -cover -coverprofile=coverage.txt -covermode=atomic

build:
	@CGO_ENABLED=0 GOOS=linux go build -a -o build/bin/orchestrate-hashicorp-vault-plugin

lint-tools: ## Install linting tools
	@GO111MODULE=on go get github.com/client9/misspell/cmd/misspell@v0.3.4
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.27.0

lint:
	@misspell -w $(GOFILES)
	@golangci-lint run --fix

lint-ci: ## Check linting
	@misspell -error $(GOFILES)
	@golangci-lint run

prod: build
	@docker-compose -f docker-compose.yml up --build vault-init vault
dev: build
	@docker-compose -f docker-compose.yml up --build vault-dev-init vault-dev
down:
	@docker-compose -f docker-compose.yml down --volumes --timeout 0

