package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/soksak-ai/soksak-plugin-registry/internal/registrysigning"
)

func main() {
	seed, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		fail(err)
	}
	private, err := registrysigning.PrivateKey(strings.TrimSpace(seed))
	if err != nil {
		fail(err)
	}
	trust, err := registrysigning.NewTrustRoot(private)
	if err != nil {
		fail(err)
	}
	if err := json.NewEncoder(os.Stdout).Encode(trust); err != nil {
		fail(err)
	}
}
func fail(err error) { fmt.Fprintln(os.Stderr, err); os.Exit(1) }
