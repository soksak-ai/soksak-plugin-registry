package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/soksak-ai/soksak-plugin-registry/internal/registrysigning"
)

func main() {
	private, err := registrysigning.PrivateKey(os.Getenv("SOKSAK_REGISTRY_ED25519_SEED"))
	if err != nil {
		fail(err)
	}
	trust, err := registrysigning.NewTrustRoot(os.Getenv("SOKSAK_REGISTRY_KEY_ID"), private)
	if err != nil {
		fail(err)
	}
	if err := json.NewEncoder(os.Stdout).Encode(trust); err != nil {
		fail(err)
	}
}
func fail(err error) { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
