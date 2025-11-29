.PHONY: help
help:
	@echo "Available targets:"
	@echo "  make test              - Run all tests"
	@echo "  make test-verbose      - Run tests with verbose output"
	@echo "  make test-coverage     - Run tests with coverage report"
	@echo "  make test-integration  - Run integration tests with Docker databases"
	@echo "  make lint              - Run linter"
	@echo "  make fmt               - Format code"
	@echo "  make build             - Build all packages"
	@echo "  make clean             - Clean build artifacts"
	@echo "  make deps              - Download dependencies"
	@echo "  make deps-verify       - Verify dependencies"
	@echo "  make vuln-check        - Check for vulnerabilities"

.PHONY: test
test:
	go test ./... -count=1

.PHONY: test-verbose
test-verbose:
	go test -v ./... -count=1

.PHONY: test-integration
test-integration:
	@echo "Starting integration tests with Docker databases..."
	@cd persistence/integration_test && powershell -ExecutionPolicy Bypass -File ./run-integration-tests.ps1

.PHONY: test-coverage
test-coverage:
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

.PHONY: lint
lint:
	golangci-lint run --timeout=5m

.PHONY: fmt
fmt:
	gofmt -s -w .
	goimports -w .

.PHONY: build
build:
	go build -v ./...

.PHONY: clean
clean:
	go clean -cache -testcache -modcache
	rm -f coverage.out coverage.html

.PHONY: deps
deps:
	go mod download
	go mod tidy

.PHONY: deps-verify
deps-verify:
	go mod verify

.PHONY: vuln-check
vuln-check:
	@which govulncheck > /dev/null || go install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./...

.PHONY: ci
ci: deps-verify lint test build vuln-check
	@echo "CI checks passed"
