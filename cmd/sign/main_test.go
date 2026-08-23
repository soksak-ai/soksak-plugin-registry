package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"testing"
	"time"

	registry "github.com/soksak-ai/soksak-contract-registry"
)

func TestSignRegistryRequiresThePinnedTrustRoot(t *testing.T) {
	public, private, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	source := registry.Registry{ID: "official", Sequence: 2, Plugins: []registry.PluginRelease{}, Sidecars: []registry.SidecarRelease{}, Kits: []registry.KitRelease{}, Contracts: []registry.ContractRelease{}, Specs: []registry.SpecRelease{}}
	trust := trustRoot{Algorithm: "ed25519", KeyID: "official-1", Value: base64.StdEncoding.EncodeToString(public)}
	document, err := signRegistry(source, trust, private, "2026-08-23T00:00:00Z", "2026-09-23T00:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	if document.Sequence != source.Sequence || document.KeyID != trust.KeyID || document.Signature == "" {
		t.Fatalf("signed document = %+v", document)
	}
	if _, err := registry.Verify(document, registry.Trust{RegistryID: "official", KeyID: trust.KeyID, PublicKey: public}, time.Date(2026, 8, 24, 0, 0, 0, 0, time.UTC), nil); err != nil {
		t.Fatal(err)
	}
}

func TestSignRegistryRejectsASecretForAnotherTrustRoot(t *testing.T) {
	public, _, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	_, private, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	source := registry.Registry{ID: "official", Sequence: 1, Plugins: []registry.PluginRelease{}, Sidecars: []registry.SidecarRelease{}, Kits: []registry.KitRelease{}, Contracts: []registry.ContractRelease{}, Specs: []registry.SpecRelease{}}
	trust := trustRoot{Algorithm: "ed25519", KeyID: "official-1", Value: base64.StdEncoding.EncodeToString(public)}
	if _, err := signRegistry(source, trust, private, "2026-08-23T00:00:00Z", "2026-09-23T00:00:00Z"); err == nil {
		t.Fatal("private key outside the pinned trust root was accepted")
	}
}
