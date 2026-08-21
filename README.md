# soksak-plugin-registry

The official Soksak release catalogue. The public format and signing rules are owned by
`soksak-spec` and enforced by `soksak-contract-registry@0.0.1`.

The catalogue publishes direct plugin, sidecar, kit, contract, and spec releases. Owner repositories
build and verify their own release archives and conformance reports. Runtime requirements remain in
owner manifests; user provider selection remains in settings. This repository never reads or builds
an owner source tree.

## Registering

A PR edits `registry-source.json` with immutable release references, target archive sizes and
digests, and conformance report digests. Run:

```sh
go test ./...
go vet ./...
go run ./cmd/verify
```

CI signs the complete source and publishes `registry-signed.json`. The Ed25519 private key is a CI
secret and is never stored in this repository.
