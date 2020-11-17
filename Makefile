.ONESHELL:
GOFILES := $(shell find . -name '*.go' -not -path "./vendor/*" | egrep -v "^\./\.go" | grep -v _test.go)
DATE = $(shell date +'%s')

test:
	go test  ./... -cover -coverprofile=coverage.txt -covermode=atomic
build:
	@CGO_ENABLED=1 GOOS=linux go build -a -v -o build/bin/orchestrate-hashicorp-vault-plugin
lint-tools: ## Install linting tools
	@GO111MODULE=on go get github.com/client9/misspell/cmd/misspell@v0.3.4
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.27.0
lint:
	@misspell -w $(GOFILES)
	@golangci-lint run --fix
lint-ci: ## Check linting
	@misspell -error $(GOFILES)
	@golangci-lint run
up:
	docker-compose -f docker/docker-compose.yml up --build --remove-orphans
down:
	sudo bash ./docker/clear.sh
	@docker-compose -f docker/docker-compose.yml down --volumes --timeout 0
