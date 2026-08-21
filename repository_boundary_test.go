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
