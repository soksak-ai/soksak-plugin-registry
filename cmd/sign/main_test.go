package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	registry "github.com/soksak-ai/soksak-contract-registry"
	"github.com/soksak-ai/soksak-plugin-registry/internal/registrysigning"
)

func TestSignRegistryRequiresASequenceAdvanceOverThePublishedIndex(t *testing.T) {
	path := filepath.Join(t.TempDir(), "registry-signed.json")
	current := registry.SignedRegistry{Registry: registry.Registry{ID: "official", Sequence: 3, Plugins: []registry.PluginRelease{}, Sidecars: []registry.SidecarRelease{}, Kits: []registry.KitRelease{}, Contracts: []registry.ContractRelease{}, Specs: []registry.SpecRelease{}}, IssuedAt: "2026-08-23T00:00:00Z", ExpiresAt: "2026-09-23T00:00:00Z", Algorithm: "ed25519", KeyID: "fixture", Signature: "fixture"}
	body, err := json.Marshal(current)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, body, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := requireSequenceAdvance(path, 3); err == nil {
		t.Fatal("same-sequence resign was accepted")
	}
	if err := requireSequenceAdvance(path, 4); err != nil {
		t.Fatal(err)
	}
}

func TestSignRegistryRequiresThePinnedTrustRoot(t *testing.T) {
	public, private, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	source := registry.Registry{ID: "official", Sequence: 2, Plugins: []registry.PluginRelease{}, Sidecars: []registry.SidecarRelease{}, Kits: []registry.KitRelease{}, Contracts: []registry.ContractRelease{}, Specs: []registry.SpecRelease{}}
	trust := registrysigning.TrustRoot{Algorithm: "ed25519", KeyID: "official-1", Value: base64.StdEncoding.EncodeToString(public)}
	document, err := registrysigning.Sign(source, trust, private, "2026-08-23T00:00:00Z", "2026-09-23T00:00:00Z")
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
	trust := registrysigning.TrustRoot{Algorithm: "ed25519", KeyID: "official-1", Value: base64.StdEncoding.EncodeToString(public)}
	if _, err := registrysigning.Sign(source, trust, private, "2026-08-23T00:00:00Z", "2026-09-23T00:00:00Z"); err == nil {
		t.Fatal("private key outside the pinned trust root was accepted")
	}
}
