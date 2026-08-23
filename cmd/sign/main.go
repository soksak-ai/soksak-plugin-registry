package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
	if err := requireSequenceAdvance("registry-signed.json", source.Sequence); err != nil {
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

func requireSequenceAdvance(path string, next uint64) error {
	body, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(bytes.NewReader(body))
	var current registry.SignedRegistry
	if err := decoder.Decode(&current); err != nil {
		return fmt.Errorf("read published registry: %w", err)
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return fmt.Errorf("published registry has trailing data")
	}
	if err := registry.Validate(current.Registry); err != nil {
		return fmt.Errorf("published registry is invalid: %w", err)
	}
	if next <= current.Sequence {
		return fmt.Errorf("registry sequence must advance beyond published sequence %d", current.Sequence)
	}
	return nil
}
func required(name string) string {
	value := os.Getenv(name)
	if value == "" {
		fail(fmt.Errorf("%s is required", name))
	}
	return value
}
func fail(err error) { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
