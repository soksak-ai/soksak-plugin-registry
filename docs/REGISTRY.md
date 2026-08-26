# Registry contribution and publication

## Pull request

A plugin owner adds or replaces exactly one `plugins/<plugin-id>.json` file:

```json
{
  "id": "soksak-plugin-example",
  "version": "1.2.3"
}
```

The entry is a dependency intent: id and version, nothing else. An entry that carries `url`,
`size`, or `sha256` is refused; the location is derived and the digest is read at publication.
The file name and the entry id must agree, and the release document at the derived location
`https://github.com/soksak-ai/<id>/releases/download/v<version>/release.json` must declare the same
id and version. The release and every transitive plugin or Sidecar dependency are verified by byte
size, SHA-256, kind, ID, version, manifest, artifact matrix, and conformance evidence, each at its
derived location. Source checkouts, branches, `latest`, package registry fallback, and repository
topology discovery are forbidden.

## Publication

Main publication derives time from the exact commit timestamp. If the same commit is run again it
reuses the same sequence and recreates identical authenticated bytes. A new commit advances the
highest immutable `registry-N` release to `N+1`. The first plugins-only Registry continues the old
signed sequence 10 at sequence 11.

`make build` reads every entry, fetches the release document at the derived location, and writes
the index entry `{id, version, size, sha256}` from the fetched bytes. The index copies nothing else
from the release document: Core walks `runtimeDependencies` from each release document. `make
authenticate` signs the index. The signing seed exists only as a GitHub secret. The public key is
embedded in Core, outside the downloaded Registry.

The `soksak-spec` package is declared by exact version in `package.json` and pinned by integrity in
`pnpm-lock.yaml`; the package registry is the `REGISTRY` make argument. The repository contains no
generated Registry file and no independent signing, parsing, or resolving code.
