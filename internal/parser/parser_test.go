package parser_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/sosukesuzuki/regexpp-go/internal/parser"
)

const fixtures = "./fixtures"

func TestParsePattern(t *testing.T) {
	u := os.Getenv("UPDATE") == "true"
	target := os.Getenv("TARGET")

	fixtureDirs, err := ioutil.ReadDir(fixtures)
	if err != nil {
		t.Error("Failed to read fixtures dir")
	}

	for _, dir := range fixtureDirs {
		if !dir.IsDir() {
			continue
		}
		if target != "" && target != dir.Name() {
			continue;
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
			t.Errorf("%s: (%s)", fixtureDirPath, err.Error())
		}
		outputPath := filepath.Join(fixtureDirPath, "output.json")
		if u {
			j, err := json.MarshalIndent(pattern, "", "  ")
			if err != nil {
				t.Errorf(err.Error())
			}
			os.WriteFile(outputPath, j, 0660)
		} else {
			bytes1, err := ioutil.ReadFile(outputPath)
			if err != nil {
				t.Error("Failed to read output.json file")
			}
			bytes2, err := json.MarshalIndent(pattern, "", "  ")
			if err != nil {
				t.Error("Failed to marshal pattern")
			}

			var p1, p2 interface{}

			json.Unmarshal(bytes1, &p1)
			json.Unmarshal(bytes2, &p2)

			// TODO: more readable diff view
			if !reflect.DeepEqual(p1, p2) {
				t.Error("Diff")
			}
		}
	}
}
