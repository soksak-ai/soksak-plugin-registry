# Registry contribution and publication

## Pull request

A plugin owner adds or replaces exactly one `plugins/<plugin-id>.json` file:

```json
{
  "id": "soksak-plugin-example",
  "version": "1.2.3",
  "url": "https://github.com/example/soksak-plugin-example/releases/download/v1.2.3/release.json",
  "size": 1234,
  "sha256": "..."
}
```

The file name, entry identity, and fetched plugin release identity must agree. The release and every
transitive plugin or Sidecar dependency are verified by URL, byte size, SHA-256, kind, ID, version,
manifest, artifact matrix, and conformance evidence. Source checkouts, branches, `latest`, package
registry fallback, and repository topology discovery are forbidden.

## Publication

Main publication derives time from the exact commit timestamp. If the same commit is run again it
reuses the same sequence and recreates identical authenticated bytes. A new commit advances the
highest immutable `registry-N` release to `N+1`. The first plugins-only Registry continues the old
signed sequence 10 at sequence 11.

The workflow downloads one exact `soksak-spec` tarball, verifies its digest, builds and authenticates
the Registry, then uses that package's immutable publisher. The signing seed exists only as a GitHub
secret. The public key is embedded in Core, outside the downloaded Registry.

The repository contains no generated Registry file and no independent signing or parsing code.
