# terraform-provider-penguin: project overview

## Purpose
- This repo currently matches HashiCorp’s `terraform-provider-scaffolding` template for Terraform Plugin Framework.
- The Go provider implementation is still named **scaffolding** in code/docs (`internal/provider/provider.go`, `main.go`, `docs/index.md`) and contains example resource/data source/function/action implementations.
- There is also a `penguin/` directory containing API documentation and TypeScript types for a separate “Penguin” service (health endpoints + Tencent Cloud CVM API), which looks like the intended real target for a future provider, but it is not yet wired into the Terraform provider code.

## Tech stack
- Go `1.24` module (`go.mod`), HashiCorp Terraform Plugin Framework.
- Tooling: `golangci-lint`, `gofmt`, `goreleaser` config, `terraform-plugin-docs` + `copywrite` invoked via `go generate`.

## Entry point
- Provider binary entrypoint: `main.go` (calls `providerserver.Serve(...)`).
- Provider implementation: `internal/provider/*`.
