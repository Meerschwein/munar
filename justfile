
default: lint fmt

lint:
    go vet ./...
    staticcheck ./...

fmt:
    gofumpt -w --extra ./