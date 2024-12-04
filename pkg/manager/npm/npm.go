package npm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/depshubhq/depshub/pkg/types"
)

type Npm struct{}

type PackageJSON struct {
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

func (Npm) GetType() types.ManagerType {
	return types.Npm
}

func (Npm) Managed(path string) bool {
	return filepath.Base(path) == "package.json"
}

func (Npm) Dependencies(path string) ([]types.Dependency, error) {
	var dependencies []types.Dependency

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var packageJSON PackageJSON
	if err := json.Unmarshal(file, &packageJSON); err != nil {
		return nil, err
	}

	// Add regular dependencies
	for name, version := range packageJSON.Dependencies {
		line, rawLine := findLineInfo(file, "dependencies", name)
		dependencies = append(dependencies, types.Dependency{
			Name: name,
			//  TODO We should use the version from the lockfile instead
			Version: cleanVersion(version),
			Dev:     false,
			Definition: types.Definition{
				Path:    path,
				RawLine: rawLine,
				Line:    line,
			},
		})
	}

	// Add dev dependencies
	for name, version := range packageJSON.DevDependencies {
		line, rawLine := findLineInfo(file, "devDependencies", name)
		dependencies = append(dependencies, types.Dependency{
			Name: name,
			//  TODO We should use the version from the lockfile instead
			Version: cleanVersion(version),
			Dev:     true,
			Definition: types.Definition{
				Path:    path,
				RawLine: rawLine,
				Line:    line,
			},
		})
	}

	// Some of the rules require the original order of dependencies
	// Sort dependencies by line number
	sort.Slice(dependencies, func(i, j int) bool {
		return dependencies[i].Line < dependencies[j].Line
	})

	return dependencies, nil
}

func (Npm) LockfilePath(path string) (string, error) {
	lockfilePath := filepath.Join(filepath.Dir(path), "package-lock.json")

	if _, err := os.Stat(lockfilePath); os.IsNotExist(err) {
		return "", fmt.Errorf("lockfile not found")
	}

	return lockfilePath, nil
}

// Returns the version without any prefix or suffix
func cleanVersion(version string) string {
	return strings.Trim(version, "v^~*><= ")
}

func findLineInfo(data []byte, section string, key string) (line int, rawLine string) {
	lines := bytes.Split(data, []byte{'\n'})
	inSection := false
	quotedKey := `"` + key + `"`

	for i, line := range lines {
		trimmed := bytes.TrimSpace(line)

		// Check if we're entering the right section
		if bytes.Contains(trimmed, []byte(`"`+section+`"`)) {
			inSection = true
			continue
		}

		// Check if we're leaving the section
		if inSection && bytes.Contains(trimmed, []byte("}")) {
			inSection = false
			continue
		}

		// Look for our key while in the correct section
		if inSection && bytes.Contains(trimmed, []byte(quotedKey)) {
			return i + 1, string(trimmed)
		}
	}

	return 0, ""
}
