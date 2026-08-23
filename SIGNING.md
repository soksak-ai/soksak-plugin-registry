# Registry signing status

`registry-source.json` is the validated catalogue source. The manual signing workflow validates,
signs, self-verifies, and publishes only `registry-signed.json`. The production private seed exists
only as a GitHub Actions secret; `registry-trust.json` and Core contain its public trust root.

Trust rotation is one explicit operation:

1. Approve rotation of the Core trust root.
2. Generate one Ed25519 key pair in a controlled environment.
3. Store the private key only as a GitHub Actions secret.
4. Pin the new public key and key ID in Core.
5. Run the manual workflow with an explicit expiry.
6. Test signature rejection, sequence rollback rejection, and key continuity through installed acceptance.

No development fallback, unsigned acceptance path, or committed private key is permitted.

## Expiry renewal

`issuedAt` and `expiresAt` are signed bytes. Publishing new timestamps at the same sequence would
be equivocation, so the signer refuses any sequence that does not advance beyond the published
index. Before expiry:

1. Run `go run ./cmd/renew`.
2. Validate and commit the sequence-only `registry-source.json` change.
3. Run the manual signing workflow with a later explicit expiry.

An ordinary catalogue registration already advances sequence and needs no separate renewal.
