package main

import (
	"encoding/json"
	"testing"

	registry "github.com/soksak-ai/soksak-contract-registry"
)

func TestRenewAdvancesOnlyTheRegistrySequence(t *testing.T) {
	source := registry.Registry{ID: "official", Sequence: 3, Plugins: []registry.PluginRelease{}, Sidecars: []registry.SidecarRelease{}, Kits: []registry.KitRelease{}, Contracts: []registry.ContractRelease{}, Specs: []registry.SpecRelease{}}
	renewed, err := renew(source)
	if err != nil {
		t.Fatal(err)
	}
	if renewed.Sequence != 4 {
		t.Fatalf("sequence = %d", renewed.Sequence)
	}
	renewed.Sequence = source.Sequence
	before, _ := json.Marshal(source)
	after, _ := json.Marshal(renewed)
	if string(after) != string(before) {
		t.Fatal("renew changed catalogue content")
	}
}

func TestRenewRejectsSequenceExhaustion(t *testing.T) {
	source := registry.Registry{ID: "official", Sequence: ^uint64(0), Plugins: []registry.PluginRelease{}, Sidecars: []registry.SidecarRelease{}, Kits: []registry.KitRelease{}, Contracts: []registry.ContractRelease{}, Specs: []registry.SpecRelease{}}
	if _, err := renew(source); err == nil {
		t.Fatal("sequence exhaustion was accepted")
	}
}
