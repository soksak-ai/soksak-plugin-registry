package registrysigning

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"testing"
)

func TestSeedDerivesThePrivateKeyWithoutPersistingIt(t *testing.T) {
	seed := make([]byte, ed25519.SeedSize)
	for index := range seed {
		seed[index] = byte(index)
	}
	private, err := PrivateKey(base64.StdEncoding.EncodeToString(seed))
	if err != nil || len(private) != ed25519.PrivateKeySize {
		t.Fatalf("private key length=%d err=%v", len(private), err)
	}
	trust, err := NewTrustRoot(private)
	if err != nil || trust.Value == "" {
		t.Fatalf("trust=%+v err=%v", trust, err)
	}
	if trust.KeyID != "56475aa75463474c0285df5dbf2bcab7" {
		t.Fatalf("key id = %s", trust.KeyID)
	}
}

func TestTrustRootRejectsAKeyIDOutsideItsPublicKey(t *testing.T) {
	trust := TrustRoot{Algorithm: "ed25519", KeyID: "wrong", Value: "A6EHv/POEL4dcN0Y50vAmWfk1jCbpQ1fHdyGZBJVMbg="}
	body, err := json.Marshal(trust)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ParseTrustRoot(body); err == nil {
		t.Fatal("unrelated key id was accepted")
	}
}
