package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"

	registry "github.com/soksak-ai/soksak-contract-registry"
)

func main() {
	if len(os.Args) < 2 {
		fail(fmt.Errorf("usage: register <release.json path or HTTPS URL> [...]"))
	}
	body, err := os.ReadFile("registry-source.json")
	if err != nil {
		fail(err)
	}
	source, err := registry.Parse(body)
	if err != nil {
		fail(err)
	}
	documents := make([][]byte, 0, len(os.Args)-1)
	for _, location := range os.Args[1:] {
		document, err := readReleaseDocument(location)
		if err != nil {
			fail(err)
		}
		documents = append(documents, document)
	}
	merged, err := mergeReleaseDocuments(source, documents)
	if err != nil {
		fail(err)
	}
	encoded, err := json.MarshalIndent(merged, "", "  ")
	if err != nil {
		fail(err)
	}
	if err := os.WriteFile("registry-source.json.next", append(encoded, '\n'), 0o644); err != nil {
		fail(err)
	}
	if err := os.Rename("registry-source.json.next", "registry-source.json"); err != nil {
		fail(err)
	}
}

func readReleaseDocument(location string) ([]byte, error) {
	if !strings.HasPrefix(location, "https://") {
		return os.ReadFile(location)
	}
	response, err := http.Get(location)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s: %s", location, response.Status)
	}
	return io.ReadAll(io.LimitReader(response.Body, 4<<20))
}

func mergeReleaseDocuments(source registry.Registry, documents [][]byte) (registry.Registry, error) {
	baseline := source
	baseline.Sequence = 0
	for _, document := range documents {
		kind, err := releaseKind(document)
		if err != nil {
			return registry.Registry{}, err
		}
		switch kind {
		case "plugin":
			var release registry.PluginRelease
			if err := decodeExact(document, &release); err != nil {
				return registry.Registry{}, err
			}
			source.Plugins = replacePlugin(source.Plugins, release)
		case "sidecar":
			var release registry.SidecarRelease
			if err := decodeExact(document, &release); err != nil {
				return registry.Registry{}, err
			}
			source.Sidecars = replaceSidecar(source.Sidecars, release)
		case "kit":
			var release registry.KitRelease
			if err := decodeExact(document, &release); err != nil {
				return registry.Registry{}, err
			}
			source.Kits = replaceKit(source.Kits, release)
		case "contract":
			var release registry.ContractRelease
			if err := decodeExact(document, &release); err != nil {
				return registry.Registry{}, err
			}
			source.Contracts = replaceContract(source.Contracts, release)
		case "spec":
			var release registry.SpecRelease
			if err := decodeExact(document, &release); err != nil {
				return registry.Registry{}, err
			}
			source.Specs = replaceSpec(source.Specs, release)
		}
	}
	updated := source
	updated.Sequence = 0
	before, err := json.Marshal(baseline)
	if err != nil {
		return registry.Registry{}, err
	}
	after, err := json.Marshal(updated)
	if err != nil {
		return registry.Registry{}, err
	}
	if !bytes.Equal(before, after) {
		if source.Sequence == ^uint64(0) {
			return registry.Registry{}, fmt.Errorf("registry sequence exhausted")
		}
		source.Sequence++
	}
	if err := registry.Validate(source); err != nil {
		return registry.Registry{}, err
	}
	return source, nil
}

func releaseKind(document []byte) (string, error) {
	var value map[string]json.RawMessage
	if err := decodeExact(document, &value); err != nil {
		return "", err
	}
	kind := ""
	for _, candidate := range []string{"plugin", "sidecar", "kit", "contract", "spec"} {
		if _, ok := value[candidate]; ok {
			if kind != "" {
				return "", fmt.Errorf("release document declares multiple kinds")
			}
			kind = candidate
		}
	}
	if kind == "" {
		return "", fmt.Errorf("release document has no supported kind")
	}
	return kind, nil
}

func decodeExact(document []byte, value any) error {
	decoder := json.NewDecoder(bytes.NewReader(document))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(value); err != nil {
		return err
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return fmt.Errorf("release document has trailing data")
	}
	return nil
}

func replacePlugin(values []registry.PluginRelease, release registry.PluginRelease) []registry.PluginRelease {
	result := values[:0]
	for _, value := range values {
		if value.Plugin != release.Plugin {
			result = append(result, value)
		}
	}
	result = append(result, release)
	sort.Slice(result, func(i, j int) bool {
		return result[i].Plugin.ID+"@"+result[i].Plugin.Version < result[j].Plugin.ID+"@"+result[j].Plugin.Version
	})
	return result
}
func replaceSidecar(values []registry.SidecarRelease, release registry.SidecarRelease) []registry.SidecarRelease {
	result := values[:0]
	for _, value := range values {
		if value.Sidecar != release.Sidecar {
			result = append(result, value)
		}
	}
	result = append(result, release)
	sort.Slice(result, func(i, j int) bool {
		return result[i].Sidecar.ID+"@"+result[i].Sidecar.Version < result[j].Sidecar.ID+"@"+result[j].Sidecar.Version
	})
	return result
}
func replaceKit(values []registry.KitRelease, release registry.KitRelease) []registry.KitRelease {
	result := values[:0]
	for _, value := range values {
		if value.Kit != release.Kit {
			result = append(result, value)
		}
	}
	result = append(result, release)
	sort.Slice(result, func(i, j int) bool {
		return result[i].Kit.ID+"@"+result[i].Kit.Version < result[j].Kit.ID+"@"+result[j].Kit.Version
	})
	return result
}
func replaceContract(values []registry.ContractRelease, release registry.ContractRelease) []registry.ContractRelease {
	result := values[:0]
	for _, value := range values {
		if value.Contract != release.Contract {
			result = append(result, value)
		}
	}
	result = append(result, release)
	sort.Slice(result, func(i, j int) bool {
		return result[i].Contract.ID+"@"+result[i].Contract.Version < result[j].Contract.ID+"@"+result[j].Contract.Version
	})
	return result
}
func replaceSpec(values []registry.SpecRelease, release registry.SpecRelease) []registry.SpecRelease {
	result := values[:0]
	for _, value := range values {
		if value.Spec != release.Spec {
			result = append(result, value)
		}
	}
	result = append(result, release)
	sort.Slice(result, func(i, j int) bool {
		return result[i].Spec.ID+"@"+result[i].Spec.Version < result[j].Spec.ID+"@"+result[j].Spec.Version
	})
	return result
}

func fail(err error) { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
