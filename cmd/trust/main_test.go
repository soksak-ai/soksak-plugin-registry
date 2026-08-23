package main

import (
	"encoding/base64"
	"encoding/json"
	"os/exec"
	"strings"
	"testing"

	"github.com/soksak-ai/soksak-plugin-registry/internal/registrysigning"
)

func TestTrustCommandReadsTheSeedFromStdinAndWritesOnlyPublicTrust(t *testing.T) {
	seed := make([]byte, 32)
	for index := range seed {
		seed[index] = byte(index)
	}
	command := exec.Command("go", "run", ".")
	command.Stdin = strings.NewReader(base64.StdEncoding.EncodeToString(seed) + "\n")
	output, err := command.Output()
	if err != nil {
		t.Fatal(err)
	}
	var trust registrysigning.TrustRoot
	if err := json.Unmarshal(output, &trust); err != nil {
		t.Fatal(err)
	}
	if trust.KeyID != "56475aa75463474c0285df5dbf2bcab7" || strings.Contains(string(output), base64.StdEncoding.EncodeToString(seed)) {
		t.Fatalf("trust output = %s", output)
	}
}
