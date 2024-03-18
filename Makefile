# Dir where build binaries are generated. The dir should be gitignored
BUILD_OUT_DIR := "bin/"

API_OUT       := "bin/api"
API_MAIN_FILE := "cmd/api/main.go"

MIGRATION_OUT       := "bin/migration"
MIGRATION_MAIN_FILE := "cmd/migration/main.go"

ABSOLUTE_PATH := $(shell pwd)


# go binary. Change this to experiment with different versions of go.
GO       = go
MODULE   = $(shell $(GO) list -m)
SERVICE  = $(shell basename $(MODULE))

# Fetch OS info
GOVERSION=$(shell go version)
UNAME_OS=$(shell go env GOOS)
UNAME_ARCH=$(shell go env GOARCH)

MOCKGEN_VERSION := v1.6.0

VERBOSE = 0
Q 		= $(if $(filter 1,$VERBOSE),,@)
M 		= $(shell printf "\033[34;1m▶\033[0m")


BIN 	 = $(CURDIR)/bin
PKGS     = $(or $(PKG),$(shell $(GO) list ./...))

$(BIN)/%: | $(BIN) ; $(info $(M) building package: $(PACKAGE)…)
	tmp=$$(mktemp -d); \
	   env GOBIN=$(BIN) go get $(PACKAGE) \
		|| ret=$$?; \
	   rm -rf $$tmp ; exit $$ret

$(BIN)/golint: PACKAGE=golang.org/x/lint/golint

GOLINT = $(BIN)/golint

#DOCKER Variables
DOCKER_STOP = docker stop
DOCKER = docker
DOCKER_RMI = docker rmi
DOCKER_RM = docker rm

.PHONY: lint
lint: | $(GOLINT) ; $(info $(M) running golint…) @ ## Run golint
	$Q $(GOLINT) -set_exit_status $(PKGS)

.PHONY: fmt
fmt: ; $(info $(M) running gofmt…) @ ## Run gofmt on all source files
	$Q $(GO) fmt $(PKGS)

.PHONY: all
all: build

.PHONY: deps
deps:
	@echo "\n + Fetching dependencies \n"
	@go install golang.org/x/lint/golint@latest
	@go install github.com/bykof/go-plantuml@latest
	@go install github.com/golang/mock/mockgen@latest
	@go install gotest.tools/gotestsum@latest

.PHONY: pre-build
pre-build: deps mock-gen go-build-migration go-build-api

.PHONY: go-build-api ## Build the binary file for API server
go-build-api:
	@CGO_ENABLED=0 go build -v -o $(API_OUT) $(API_MAIN_FILE)

.PHONY: go-build-migration ## Build the binary file for database migrations
go-build-migration:
	@CGO_ENABLED=0 go build -v -o $(MIGRATION_OUT) $(MIGRATION_MAIN_FILE)

.PHONY: go-run-api ## Run the API server
go-run-api: go-build-api
	@go run $(API_MAIN_FILE)

.PHONY: up-migration ## Run the database migrations
up-migration: go-build-migration
	@go run $(MIGRATION_MAIN_FILE) up

.PHONY: down-migration ## Run the database migrations
down-migration: go-build-migration
	@go run $(MIGRATION_MAIN_FILE) down

.PHONY: build
build: build-info  docker-build

.PHONY: build-info
build-info:
	@echo "\nBuild Info:\n"
	@echo "\t\033[33mOS\033[0m: $(UNAME_OS)"
	@echo "\t\033[33mArch\033[0m: $(UNAME_ARCH)"
	@echo "\t\033[33mGo Version\033[0m: $(GOVERSION)\n"

.PHONY: docker-build
docker-build: docker-build-api docker-build-migration

.PHONY: docker-build-api
docker-build-api:
	@docker build . -f docker/Dockerfile.api -t pizmo/api:latest

.PHONY: docker-build-migration
docker-build-migration:
	@docker build . -f build/Dockerfile.migration -t pizmo/migration:latest

.PHONY: dev-docker-up ## Bring up docker-compose for local dev-setup
dev-docker-up:
	docker-compose -f docker/docker-compose.yaml up -d --build

.PHONY: mock-gen ## Generates mocks
mock-gen: ## generates mocks
	@which mockgen || go install github.com/golang/mock/mockgen@$(MOCKGEN_VERSION)
	mockgen -source=$(ABSOLUTE_PATH)/internal/account/core.go -destination=$(ABSOLUTE_PATH)/internal/account/mock/mock_core.go -package=mock
	mockgen -source=$(ABSOLUTE_PATH)/internal/account/repo.go -destination=$(ABSOLUTE_PATH)/internal/account/mock/mock_repo.go -package=mock
	mockgen -source=$(ABSOLUTE_PATH)/internal/transaction/core.go -destination=$(ABSOLUTE_PATH)/internal/transaction/mock/mock_core.go -package=mock
	mockgen -source=$(ABSOLUTE_PATH)/internal/transaction/repo.go -destination=$(ABSOLUTE_PATH)/internal/transaction/mock/mock_repo.go -package=mock

.PHONY: test
test: ## Run tests
	@echo "\n + Running tests\n"
	@go test -v ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	@echo "\n + Running tests with coverage\n"
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out

.PHONY: clean
clean: ## Remove previous builds
	@echo " + Removing cloned and generated files\n"
	@rm -rf $(API_OUT) $(MIGRATION_OUT)

check-swagger:
	which swagger || (GO111MODULE=off go get -u github.com/go-swagger/go-swagger/cmd/swagger)

swagger: check-swagger
	GO111MODULE=on go mod vendor  && GO111MODULE=off swagger generate spec -o ./swagger.yaml --scan-models

serve-swagger: check-swagger
	swagger serve -F=swagger swagger.yaml

