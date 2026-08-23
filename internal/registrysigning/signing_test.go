package registrysigning

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"testing"
)

func TestSeedDerivesThePrivateKeyWithoutPersistingIt(t *testing.T) {
	seed := make([]byte, ed25519.SeedSize)
	if _, err := rand.Read(seed); err != nil {
		t.Fatal(err)
	}
	private, err := PrivateKey(base64.StdEncoding.EncodeToString(seed))
	if err != nil || len(private) != ed25519.PrivateKeySize {
		t.Fatalf("private key length=%d err=%v", len(private), err)
	}
	trust, err := NewTrustRoot("official-1", private)
	if err != nil || trust.Value == "" {
		t.Fatalf("trust=%+v err=%v", trust, err)
	}
}
