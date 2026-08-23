package registry

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/soksak-ai/soksak-plugin-registry/internal/registrysigning"
)

func TestRegistryDoesNotExecuteOwnerSourceTrees(t *testing.T) {
	root, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	tokens := []string{"../" + "core", "soksak-plugin-" + "doctor", "/plugin." + "json", "/main." + "js"}
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if path != root && info.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) != ".sh" && filepath.Ext(path) != ".go" {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		line := 0
		for scanner.Scan() {
			line++
			for _, token := range tokens {
				if strings.Contains(scanner.Text(), token) {
					t.Errorf("%s:%d executes or reads owner source through %s", path, line, token)
				}
			}
		}
		return scanner.Err()
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestREADMEStatesTheActualSigningStatus(t *testing.T) {
	body, err := os.ReadFile("README.md")
	if err != nil {
		t.Fatal(err)
	}
	text := string(body)
	for _, required := range []string{"not yet signed", "SIGNING.md", "not operational"} {
		if !strings.Contains(text, required) {
			t.Errorf("README omits signing status: %s", required)
		}
	}
	if _, err := os.Stat("README.ko.md"); err != nil {
		t.Errorf("README has no Korean translation: %v", err)
	}
}

func TestSigningWorkflowIsManualAndPublishesOnlyTheSignedIndex(t *testing.T) {
	body, err := os.ReadFile(".github/workflows/sign.yml")
	if err != nil {
		t.Fatal(err)
	}
	text := string(body)
	for _, required := range []string{
		"workflow_dispatch:",
		"go-version-file: go.mod",
		"go run ./cmd/verify",
		"go run ./cmd/sign",
		"SOKSAK_REGISTRY_ED25519_SEED",
		"registry-signed.json",
		"git diff --quiet -- registry-source.json registry-trust.json",
		"git add -- registry-signed.json",
		"permission-contents: write",
	} {
		if !strings.Contains(text, required) {
			t.Errorf("signing workflow omits %s", required)
		}
	}
	for _, forbidden := range []string{"schedule:", "push:", "SOKSAK_REGISTRY_ED25519_SEED=", "git add ."} {
		if strings.Contains(text, forbidden) {
			t.Errorf("signing workflow contains %s", forbidden)
		}
	}
	signer, err := os.ReadFile("cmd/sign/main.go")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(signer), "requireSequenceAdvance") {
		t.Fatal("signer does not reject same-sequence publication")
	}
}

func TestCommittedTrustRootIsValid(t *testing.T) {
	body, err := os.ReadFile("registry-trust.json")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := registrysigning.ParseTrustRoot(body); err != nil {
		t.Fatal(err)
	}
}
