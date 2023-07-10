# Salad Bowl

WIP - Full functionality has not yet been completed. This project is a work in progress.

Live at: https://saladbowl.bensivo.com

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


The included taskfile has shortcuts for all the basic development tasks:
- Run unit tests: `task test`
- Start local env: `task up`
- Run e2e tests: `task e2e`
- Stop local env: `task down`
