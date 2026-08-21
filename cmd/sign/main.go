package main

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	registry "github.com/soksak-ai/soksak-contract-registry"
)

func main() {
	body, err := os.ReadFile("registry-source.json")
	if err != nil {
		fail(err)
	}
	source, err := registry.Parse(body)
	if err != nil {
		fail(err)
	}
	privateRaw, err := base64.StdEncoding.DecodeString(os.Getenv("SOKSAK_REGISTRY_ED25519_PRIVATE_KEY"))
	if err != nil || len(privateRaw) != ed25519.PrivateKeySize {
		fail(fmt.Errorf("SOKSAK_REGISTRY_ED25519_PRIVATE_KEY must be a base64 Ed25519 private key"))
	}
	document := registry.SignedRegistry{Registry: source, IssuedAt: required("SOKSAK_REGISTRY_ISSUED_AT"), ExpiresAt: required("SOKSAK_REGISTRY_EXPIRES_AT"), Algorithm: "ed25519", KeyID: required("SOKSAK_REGISTRY_KEY_ID")}
	if err := registry.Sign(&document, ed25519.PrivateKey(privateRaw)); err != nil {
		fail(err)
	}
	encoded, err := json.MarshalIndent(document, "", "  ")
	if err != nil {
		fail(err)
	}
	if err := os.WriteFile("registry-signed.json.next", append(encoded, byte(10)), 0o600); err != nil {
		fail(err)
	}
	if err := os.Rename("registry-signed.json.next", "registry-signed.json"); err != nil {
		fail(err)
	}
}
func required(name string) string {
	value := os.Getenv(name)
	if value == "" {
		fail(fmt.Errorf("%s is required", name))
	}
	return value
}
func fail(err error) { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
