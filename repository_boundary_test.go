package registry

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"
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
	if strings.Contains(text, "CI signs the complete source") {
		t.Fatal("README claims a signing workflow that does not exist")
	}
	for _, required := range []string{"not yet signed", "SIGNING.md", "not operational"} {
		if !strings.Contains(text, required) {
			t.Errorf("README omits signing status: %s", required)
		}
	}
	if _, err := os.Stat("README.ko.md"); err != nil {
		t.Errorf("README has no Korean translation: %v", err)
	}
}
