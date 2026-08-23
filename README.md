# soksak-plugin-registry

The official Soksak release catalogue. The public format and signing rules are owned by
`soksak-spec` and enforced by `soksak-contract-registry@0.0.1`.

The catalogue publishes direct plugin, sidecar, kit, contract, and spec releases. Owner repositories
build and verify their own release archives and conformance reports. Runtime requirements remain in
owner manifests; user sidecar role bindings remain in `environment.json`. This repository never reads or builds
an owner source tree.

## Registering

A PR registers owner-produced immutable `release.json` documents. The command replaces the same
exact release identity, sorts each release kind, validates the complete catalogue, and only then
atomically updates `registry-source.json`:

```sh
go run ./cmd/register https://github.com/soksak-ai/example/releases/download/v0.0.1/release.json
go test ./...
go vet ./...
go run ./cmd/verify
```

`registry-signed.json` is the operational signed index. Core pins the matching public trust root,
and [SIGNING.md](SIGNING.md) defines expiry renewal and sequence continuity. Unsigned catalogue
bytes are never an installation source.
