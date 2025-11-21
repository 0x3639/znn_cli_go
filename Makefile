.PHONY: help build install test lint security clean run fmt vet

# Binary name
BINARY_NAME=znn-cli
INSTALL_PATH=$(GOPATH)/bin

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
GOVET=$(GOCMD) vet

# Build flags
LDFLAGS=-ldflags "-s -w"

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Build the binary
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) -v

install: build ## Install the binary to GOPATH/bin
	cp $(BINARY_NAME) $(INSTALL_PATH)/

run: build ## Run the application
	./$(BINARY_NAME)

test: ## Run tests
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

test-coverage: test ## Run tests and show coverage
	$(GOCMD) tool cover -func=coverage.out

lint: ## Run golangci-lint
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Install with: brew install golangci-lint" && exit 1)
	golangci-lint run ./...

security: ## Run gosec security scanner
	@which gosec > /dev/null || (echo "gosec not installed. Install with: go install github.com/securego/gosec/v2/cmd/gosec@latest" && exit 1)
	gosec -conf .gosec.yml ./...

fmt: ## Format code
	$(GOFMT) -s -w .
	goimports -w .

vet: ## Run go vet
	$(GOVET) ./...

clean: ## Clean build artifacts
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html

deps: ## Download dependencies
	$(GOMOD) download
	$(GOMOD) tidy

update-deps: ## Update dependencies
	$(GOGET) -u ./...
	$(GOMOD) tidy

all: clean deps fmt vet lint security test build ## Run all checks and build
