package registrysigning

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"time"

	registry "github.com/soksak-ai/soksak-contract-registry"
)

type TrustRoot struct {
	Algorithm string `json:"algorithm"`
	KeyID     string `json:"keyId"`
	Value     string `json:"value"`
}

func PrivateKey(seedBase64 string) (ed25519.PrivateKey, error) {
	seed, err := base64.StdEncoding.DecodeString(seedBase64)
	if err != nil || len(seed) != ed25519.SeedSize {
		return nil, fmt.Errorf("registry signing seed must be base64 Ed25519 seed bytes")
	}
	return ed25519.NewKeyFromSeed(seed), nil
}

func NewTrustRoot(private ed25519.PrivateKey) (TrustRoot, error) {
	if len(private) != ed25519.PrivateKeySize {
		return TrustRoot{}, fmt.Errorf("registry private key is required")
	}
	public := private.Public().(ed25519.PublicKey)
	return TrustRoot{Algorithm: "ed25519", KeyID: keyID(public), Value: base64.StdEncoding.EncodeToString(public)}, nil
}

func keyID(public ed25519.PublicKey) string {
	digest := sha256.Sum256(public)
	return fmt.Sprintf("%x", digest[:16])
}

func ParseTrustRoot(body []byte) (TrustRoot, error) {
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.DisallowUnknownFields()
	var trust TrustRoot
	if err := decoder.Decode(&trust); err != nil {
		return TrustRoot{}, err
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return TrustRoot{}, fmt.Errorf("trust root has trailing data")
	}
	public, err := base64.StdEncoding.DecodeString(trust.Value)
	if trust.Algorithm != "ed25519" || err != nil || len(public) != ed25519.PublicKeySize || trust.KeyID != keyID(public) {
		return TrustRoot{}, fmt.Errorf("invalid registry trust root")
	}
	return trust, nil
}

func Sign(source registry.Registry, trust TrustRoot, private ed25519.PrivateKey, issuedAt, expiresAt string) (registry.SignedRegistry, error) {
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
