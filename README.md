# soksak-plugin-registry

The official Soksak plugin install catalogue. The public format and signing rules are owned by
soksak-contract-registry@0.0.1.

Users select plugins. Plugin releases reference exact sidecar and kit dependencies. Sidecars and
kits are dependency nodes rather than independent catalogue products. Owner repositories build and
verify their own release archives, composition manifests and conformance reports. This repository
never reads or builds an owner source tree.

## Registering

A PR edits `registry-source.json` with immutable release references, exact dependency identities,
target archive digests, conformance report digests and a plugin install profile. Run:

```sh
go test ./...
go vet ./...
go run ./cmd/verify
```

CI signs the complete source and publishes `registry-signed.json`. The Ed25519 private key is a CI
secret and is never stored in this repository.
