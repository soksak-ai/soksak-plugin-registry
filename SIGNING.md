# Registry signing status

`registry-source.json` is the validated catalogue source. This repository does not currently have
a signing workflow, an Ed25519 private-key secret, or a published `registry-signed.json`. The Core
therefore cannot treat this catalogue as an operational official registry.

Enabling signing is one explicit trust-root operation:

1. Approve rotation of the Core trust root.
2. Generate one Ed25519 key pair in a controlled environment.
3. Store the private key only as a GitHub Actions secret.
4. Add a workflow that validates the source, signs it, and publishes the signed document.
5. Pin the new public key and key ID in Core.
6. Test signature rejection, sequence rollback rejection, and key-rotation recovery.

No development fallback, unsigned acceptance path, or committed private key is permitted.
