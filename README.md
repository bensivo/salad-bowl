# Salad Bowl


## Development
Development dependencies:
- golang 1.20+
- docker
- node 16+
- taskfile (brew install go-task/tap/go-task)

Setting up your development environment:
```
# Download all code dependencies
go mod download

# Tool for auto-generated mocks for unit testing
go install github.com/vektra/mockery/v2@v2.30.1


# Download dependencies for e2e test suite
cd e2e && npm i
```