package registry

import (
	"os"
	"testing"

	contract "github.com/soksak-ai/soksak-contract-registry"
)

func TestRegistrySourceMatchesThePublicContract(t *testing.T) {
	body, err := os.ReadFile("registry-source.json")
	if err != nil {
		t.Fatal(err)
	}
	registry, err := contract.Parse(body)
	if err != nil {
		t.Fatal(err)
	}
	if registry.ID != "official" || registry.Sequence < 1 {
		t.Fatalf("registry = %+v", registry)
	}
}
