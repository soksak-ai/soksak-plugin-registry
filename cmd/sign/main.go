package main

import (
	"encoding/json"
	"fmt"
	"os"

	registry "github.com/soksak-ai/soksak-contract-registry"
	"github.com/soksak-ai/soksak-plugin-registry/internal/registrysigning"
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
	trustBody, err := os.ReadFile("registry-trust.json")
	if err != nil {
		fail(err)
	}
	trust, err := registrysigning.ParseTrustRoot(trustBody)
	if err != nil {
		fail(err)
	}
	private, err := registrysigning.PrivateKey(os.Getenv("SOKSAK_REGISTRY_ED25519_SEED"))
	if err != nil {
		fail(err)
	}
	document, err := registrysigning.Sign(source, trust, private, required("SOKSAK_REGISTRY_ISSUED_AT"), required("SOKSAK_REGISTRY_EXPIRES_AT"))
	if err != nil {
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
