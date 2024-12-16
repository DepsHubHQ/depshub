package cargo

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/depshubhq/depshub/pkg/types"
)

// DependencyValue represents either a string version or a map of dependency attributes
type DependencyValue struct {
	Version string
	Path    string
}

// UnmarshalTOML implements the interface for handling both string and table TOML values
func (d *DependencyValue) UnmarshalTOML(v interface{}) error {
	switch val := v.(type) {
	case string:
		d.Version = val
	case map[string]interface{}:
		if version, ok := val["version"].(string); ok {
			d.Version = version
		}
		if path, ok := val["path"].(string); ok {
			d.Path = path
		}
	}
	return nil
}

type Cargo struct{}

type CargoTOML struct {
	Dependencies      map[string]DependencyValue `toml:"dependencies"`
	DevDependencies   map[string]DependencyValue `toml:"dev-dependencies"`
	BuildDependencies map[string]DependencyValue `toml:"build-dependencies"`
}

func (Cargo) GetType() types.ManagerType {
	return types.Cargo
}

func (Cargo) Managed(path string) bool {
	path = strings.ToLower(path)
	return filepath.Base(path) == "cargo.toml"
}

func (Cargo) Dependencies(path string) ([]types.Dependency, error) {
	var dependencies []types.Dependency

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cargoTOML CargoTOML
	if err := toml.Unmarshal(file, &cargoTOML); err != nil {
		return nil, err
	}

	// Add regular dependencies
	for name, version := range cargoTOML.Dependencies {
		line, rawLine := findLineInfo(file, name)
		dependencies = append(dependencies, types.Dependency{
			Manager: types.Cargo,
			Name:    name,
			//  TODO We should use the version from the lockfile instead
			Version: cleanVersion(version.Version),
			Dev:     false,
			Definition: types.Definition{
				Path:    path,
				RawLine: rawLine,
				Line:    line,
			},
		})
	}

	// Add dev dependencies
	for name, version := range cargoTOML.DevDependencies {
		line, rawLine := findLineInfo(file, name)
		dependencies = append(dependencies, types.Dependency{
			Manager: types.Cargo,
			Name:    name,
			//  TODO We should use the version from the lockfile instead
			Version: cleanVersion(version.Version),
			Dev:     true,
			Definition: types.Definition{
				Path:    path,
				RawLine: rawLine,
				Line:    line,
			},
		})
	}

	for name, version := range cargoTOML.BuildDependencies {
		line, rawLine := findLineInfo(file, name)
		dependencies = append(dependencies, types.Dependency{
			Manager: types.Cargo,
			Name:    name,
			//  TODO We should use the version from the lockfile instead
			Version: cleanVersion(version.Version),
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

func (Cargo) LockfilePath(path string) (string, error) {
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

func findLineInfo(data []byte, key string) (line int, rawLine string) {
	lines := bytes.Split(data, []byte{'\n'})

	for i, line := range lines {
		trimmed := bytes.TrimSpace(line)

		// Look for our key while in the correct section
		if bytes.Contains(trimmed, []byte(key+" =")) {
			return i, string(trimmed)
		}
	}

	return 0, ""
}
