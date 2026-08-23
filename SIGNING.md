# Registry signing status

`registry-source.json` is the validated catalogue source. The manual signing workflow validates,
signs, self-verifies, and publishes only `registry-signed.json`. No production trust root, Ed25519
private-key secret, or signed index exists yet, so Core cannot treat this catalogue as an
operational official registry.

Enabling signing is one explicit trust-root operation:

1. Approve rotation of the Core trust root.
2. Generate one Ed25519 key pair in a controlled environment.
3. Store the private key only as a GitHub Actions secret.
4. Pin the new public key and key ID in Core.
5. Run the manual workflow with an explicit expiry.
6. Test signature rejection, sequence rollback rejection, and key continuity through installed acceptance.

No development fallback, unsigned acceptance path, or committed private key is permitted.
