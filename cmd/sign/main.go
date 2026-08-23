package main

import (
	"bytes"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	registry "github.com/soksak-ai/soksak-contract-registry"
)

type trustRoot struct {
	Algorithm string `json:"algorithm"`
	KeyID     string `json:"keyId"`
	Value     string `json:"value"`
}

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
	trust, err := parseTrustRoot(trustBody)
	if err != nil {
		fail(err)
	}
	privateRaw, err := base64.StdEncoding.DecodeString(os.Getenv("SOKSAK_REGISTRY_ED25519_PRIVATE_KEY"))
	if err != nil || len(privateRaw) != ed25519.PrivateKeySize {
		fail(fmt.Errorf("SOKSAK_REGISTRY_ED25519_PRIVATE_KEY must be a base64 Ed25519 private key"))
	}
	document, err := signRegistry(source, trust, ed25519.PrivateKey(privateRaw), required("SOKSAK_REGISTRY_ISSUED_AT"), required("SOKSAK_REGISTRY_EXPIRES_AT"))
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

func parseTrustRoot(body []byte) (trustRoot, error) {
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.DisallowUnknownFields()
	var trust trustRoot
	if err := decoder.Decode(&trust); err != nil {
		return trustRoot{}, err
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return trustRoot{}, fmt.Errorf("trust root has trailing data")
	}
	public, err := base64.StdEncoding.DecodeString(trust.Value)
	if trust.Algorithm != "ed25519" || trust.KeyID == "" || err != nil || len(public) != ed25519.PublicKeySize {
		return trustRoot{}, fmt.Errorf("invalid registry trust root")
	}
	return trust, nil
}

func signRegistry(source registry.Registry, trust trustRoot, private ed25519.PrivateKey, issuedAt, expiresAt string) (registry.SignedRegistry, error) {
	publicRaw, err := base64.StdEncoding.DecodeString(trust.Value)
	if err != nil || len(publicRaw) != ed25519.PublicKeySize || len(private) != ed25519.PrivateKeySize {
		return registry.SignedRegistry{}, fmt.Errorf("invalid registry signing key")
	}
	if !bytes.Equal(private.Public().(ed25519.PublicKey), publicRaw) {
		return registry.SignedRegistry{}, fmt.Errorf("private key does not match registry trust root")
	}
	document := registry.SignedRegistry{Registry: source, IssuedAt: issuedAt, ExpiresAt: expiresAt, Algorithm: trust.Algorithm, KeyID: trust.KeyID}
	if err := registry.Sign(&document, private); err != nil {
		return registry.SignedRegistry{}, err
	}
	issued, err := time.Parse(time.RFC3339, issuedAt)
	if err != nil {
		return registry.SignedRegistry{}, fmt.Errorf("invalid issuedAt")
	}
	if _, err := registry.Verify(document, registry.Trust{RegistryID: source.ID, KeyID: trust.KeyID, PublicKey: ed25519.PublicKey(publicRaw)}, issued, nil); err != nil {
		return registry.SignedRegistry{}, fmt.Errorf("signed registry self-verification failed: %w", err)
	}
	return document, nil
}
func required(name string) string {
	value := os.Getenv(name)
	if value == "" {
		fail(fmt.Errorf("%s is required", name))
	}
	return value
}
func fail(err error) { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
