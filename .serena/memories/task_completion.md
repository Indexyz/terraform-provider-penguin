# When you finish a task

- Run `make fmt` and `make lint`.
- Run `make test`.
- If changes affect provider behavior or schema, run docs generation (`make generate` or `cd tools && go generate -tags=generate ./...`) and confirm `docs/` updates.
- If changes affect real Terraform interactions, run `make testacc` (only when appropriate; creates real resources and may cost money).

Notes:
- This repo still contains template identifiers (`scaffolding` provider name/address). If you are implementing the real provider, expect to update `main.go`, `internal/provider/provider.go`, and `tools/tools.go`’s `-provider-name` to match the published provider name.
