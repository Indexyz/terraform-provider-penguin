# Suggested commands

## Build / install
- `go install ./...` (installs provider binary into `$GOBIN`/`$GOPATH/bin`)
- `make build`
- `make install`

## Format / lint
- `make fmt` (runs `gofmt -s -w -e .`)
- `make lint` (runs `golangci-lint run`)

## Tests
- `make test` (unit tests)
- `make testacc` (acceptance tests; sets `TF_ACC=1`, long timeout; creates real resources)

## Docs / generation
- `make generate`
  - Note: `tools/tools.go` has `//go:build generate`. If `make generate` appears to do nothing, run:
  - `cd tools && go generate -tags=generate ./...`

## Useful repo exploration
- `rg -n "TODO:" .`
- `go test ./...`
- `golangci-lint run`
