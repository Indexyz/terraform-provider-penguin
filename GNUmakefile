default: fmt lint install generate

GO_ENV := GOCACHE=$(CURDIR)/.cache/go-build GOMODCACHE=$(CURDIR)/.cache/go-mod GOPATH=$(CURDIR)/.cache/go
LINT_ENV := $(GO_ENV) GOLANGCI_LINT_CACHE=$(CURDIR)/.cache/golangci-lint

build:
	$(GO_ENV) go build -v ./...

install: build
	$(GO_ENV) go install -v ./...

lint:
	$(LINT_ENV) golangci-lint run

generate:
	cd tools; go generate -tags=generate ./...

fmt:
	find . -name '*.go' -not -path './.cache/*' -not -path './.serena/*' -print0 | xargs -0 -r gofmt -s -w -e

test:
	$(GO_ENV) go test -v -cover -timeout=120s -parallel=10 ./...

testacc:
	$(GO_ENV) TF_ACC=1 go test -v -cover -timeout 120m ./...

.PHONY: fmt lint test testacc build install generate
