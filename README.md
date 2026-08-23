# Soksak plugin registry

The official Registry is a data repository. Each `plugins/<id>.json` file names one current
immutable plugin release with the common `{id, version, url, size, sha256}` reference. Plugin
owners change only their own file. Sidecars are runtime dependencies inside the owner release and
are never registered as independent products.

Pull requests validate the changed entry and the complete transitive release chain with the
digest-pinned public `soksak-spec` package. A merge to main deterministically projects all entries,
authenticates `registry.json`, and publishes it as the only asset of an owner-enforced immutable
`registry-<sequence>` GitHub Release. Generated Registry bytes, signing keys, parsers, and publisher
implementations are not stored here. Core embeds the public trust root.

See [docs/REGISTRY.md](docs/REGISTRY.md).
