package main

import (
	"encoding/json"
	"fmt"
	"os"

	registry "github.com/soksak-ai/soksak-contract-registry"
)

func main() {
	body, err := os.ReadFile("registry-source.json")
	if err != nil {
		fail(err)
	}
	source, err := registry.Parse(body)
	if err != nil {
		fail(err)
	}
	renewed, err := renew(source)
	if err != nil {
		fail(err)
	}
	encoded, err := json.MarshalIndent(renewed, "", "  ")
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

func renew(source registry.Registry) (registry.Registry, error) {
	if source.Sequence == ^uint64(0) {
		return registry.Registry{}, fmt.Errorf("registry sequence exhausted")
	}
	source.Sequence++
	if err := registry.Validate(source); err != nil {
		return registry.Registry{}, err
	}
	return source, nil
}

func fail(err error) { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
