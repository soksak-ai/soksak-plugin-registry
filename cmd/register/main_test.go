package main

import (
	"encoding/json"
	"testing"

	registry "github.com/soksak-ai/soksak-contract-registry"
)

func TestMergeReleaseDocumentsReplacesAndSortsDirectReleases(t *testing.T) {
	source := registry.Registry{
		ID: "official", Sequence: 1,
		Plugins: []registry.PluginRelease{}, Sidecars: []registry.SidecarRelease{},
		Kits: []registry.KitRelease{}, Contracts: []registry.ContractRelease{}, Specs: []registry.SpecRelease{},
	}
	documents := [][]byte{releaseDocument(t, "plugin", "z-plugin"), releaseDocument(t, "plugin", "a-plugin")}
	merged, err := mergeReleaseDocuments(source, documents)
	if err != nil {
		t.Fatal(err)
	}
	if len(merged.Plugins) != 2 || merged.Plugins[0].Plugin.ID != "a-plugin" || merged.Plugins[1].Plugin.ID != "z-plugin" {
		t.Fatalf("plugins = %+v", merged.Plugins)
	}
	updated, err := mergeReleaseDocuments(merged, [][]byte{releaseDocument(t, "plugin", "a-plugin")})
	if err != nil {
		t.Fatal(err)
	}
	if len(updated.Plugins) != 2 {
		t.Fatalf("replacement duplicated release: %+v", updated.Plugins)
	}
}

func TestMergeReleaseDocumentsRejectsUnknownReleaseKinds(t *testing.T) {
	source := registry.Registry{ID: "official", Sequence: 1, Plugins: []registry.PluginRelease{}, Sidecars: []registry.SidecarRelease{}, Kits: []registry.KitRelease{}, Contracts: []registry.ContractRelease{}, Specs: []registry.SpecRelease{}}
	if _, err := mergeReleaseDocuments(source, [][]byte{[]byte(`{"unknown":{}}`)}); err == nil {
		t.Fatal("unknown release kind was accepted")
	}
}

func releaseDocument(t *testing.T, kind, id string) []byte {
	t.Helper()
	value := map[string]any{
		kind:        map[string]any{"id": id, "version": "0.0.1"},
		"source":    map[string]any{"repository": "https://github.com/soksak-ai/" + id, "commit": "0123456789abcdef0123456789abcdef01234567"},
		"artifacts": []map[string]any{{"target": "any", "url": "https://github.com/soksak-ai/" + id + "/releases/download/v0.0.1/" + id + "-0.0.1-any.tgz", "size": 1, "sha256": "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef", "format": "tgz", "manifest": kind + ".json"}},
		"reports":   []map[string]any{{"url": "https://github.com/soksak-ai/" + id + "/releases/download/v0.0.1/conformance-release.json", "sha256": "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"}},
	}
	body, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return body
}
