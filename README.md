# Soksak plugin registry

The official Registry is a data repository. Each `plugins/<id>.json` file names one current
immutable plugin release as the intent `{id, version}`. Plugin owners change only their own file.
Sidecars are runtime dependencies inside the owner release and are never registered as independent
products.

The release location is derived, never written: the release document of `{id, version}` is
`https://github.com/soksak-ai/<id>/releases/download/v<version>/release.json`. Pull requests
validate the changed entry, read the release document at the derived location, and verify the
complete transitive release chain with the exact `soksak-spec` package. A main push repeats the
complete verification and never publishes. An explicit publication workflow dispatch against a
validated commit deterministically projects all entries into the signed index `registry.json`,
whose plugin references `{id, version, size, sha256}` carry the size and SHA-256 of the fetched
release document and nothing else: the index does not copy `runtimeDependencies`, because Core
walks the closure from each release document. The index is authenticated and published as the only
asset of an owner-enforced immutable `registry-<sequence>` GitHub Release. Generated Registry
bytes, signing keys, parsers, and publisher implementations are not stored here. Core embeds the
public trust root.

See [docs/REGISTRY.md](docs/REGISTRY.md).

## Command and dependency ownership

The package depends on the exact immutable `@soksak/soksak-spec` release asset declared in
`package.json`, so every `make` invocation that installs requires `REGISTRY` on the make command
line. A value from the environment is refused. The Makefile reads the requirement from
`package.json` and refuses `REGISTRY required: this package depends on @soksak/soksak-spec` when it
is absent. The Spec package is fetched from its owner-enforced immutable release URL and verified
against the lockfile integrity; public npm availability is not part of the validator identity.

The build input is identified by the `pnpm-lock.yaml` integrity, not by `REGISTRY`. pnpm fetches
from `REGISTRY` only a package whose integrity its content-addressable store does not already hold,
so a second install of the same lockfile on the same machine reads the store and never contacts
`REGISTRY`.

```sh
make verify REGISTRY=http://host:port/
```

`.node-version`, `package.json#packageManager`, the exact release URL, and `pnpm-lock.yaml` own the
environment and immutable spec validation package. `make build` accepts explicit sequence/time/output inputs, and
`make authenticate` accepts explicit input/output plus the signing seed from the caller. GitHub
Actions owns tag discovery, secrets, credentials, and publication only. The publication workflow is
manual-only; a push workflow cannot receive those credentials or mutate Registry tags, Releases,
or assets. Registry is the last release train consumer: after the changed spec is immutable, this
repository updates its exact version, passes the full clean-install chain, builds authenticated
bytes, and publishes without rebuilding any plugin, Sidecar, Kit, or spec source.
