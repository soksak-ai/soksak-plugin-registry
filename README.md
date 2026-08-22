# soksak-plugin-registry

The official Soksak release catalogue. The public format and signing rules are owned by
`soksak-spec` and enforced by `soksak-contract-registry@0.0.1`.

The catalogue publishes direct plugin, sidecar, kit, contract, and spec releases. Owner repositories
build and verify their own release archives and conformance reports. Runtime requirements remain in
owner manifests; user provider selection remains in settings. This repository never reads or builds
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

CI signs the complete source and publishes `registry-signed.json`. The Ed25519 private key is a CI
secret and is never stored in this repository.
