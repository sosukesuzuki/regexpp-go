package parser_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/sosukesuzuki/regexpp-go/internal/parser"
)

const fixtures = "./fixtures"

func TestParsePattern(t *testing.T) {
	u := os.Getenv("UPDATE") == "true"

	fixtureDirs, err := ioutil.ReadDir(fixtures)
	if err != nil {
		t.Error("Failed to read fixtures dir")
	}

	for _, dir := range fixtureDirs {
		if !dir.IsDir() {
			continue
		}
		fixtureDirPath := filepath.Join(fixtures, dir.Name())
		bytes, err := ioutil.ReadFile(filepath.Join(fixtureDirPath, "input.txt"))
		if err != nil {
			t.Error("Failed to read input.txt file")
		}
		input := string(bytes)
		parser := parser.NewParser(input, true)
		pattern, err := parser.ParsePattern()
		if err != nil {
			t.Errorf(err.Error())
		}
		if u {
			j, err := json.MarshalIndent(pattern, "", "  ")
			if err != nil {
				t.Errorf(err.Error())
			}
			os.WriteFile(filepath.Join(fixtureDirPath, "output.json"), j, 0660)
		}
	}
}
