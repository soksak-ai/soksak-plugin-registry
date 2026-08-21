package main

import (
	"fmt"
	"os"

	registry "github.com/soksak-ai/soksak-contract-registry"
)

func main() {
	body, err := os.ReadFile("registry-source.json")
	if err != nil {
		fail(err)
	}
	if _, err := registry.Parse(body); err != nil {
		fail(err)
	}
}
func fail(err error) { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
