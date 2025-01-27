package pyproject

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/depshubhq/depshub/pkg/types"
	"github.com/pelletier/go-toml"
)

type Pyproject struct{}

func (Pyproject) GetType() types.ManagerType {
	return types.Pyproject
}

func (Pyproject) Managed(path string) bool {
	return filepath.Base(path) == "pyproject.toml"
}

func (Pyproject) Dependencies(path string) ([]types.Dependency, error) {
	var dependencies []types.Dependency

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	tree, err := toml.LoadReader(file)
	if err != nil {
		return nil, err
	}

	// Parse [project.dependencies]
	if deps, ok := tree.Get("project.dependencies").(*toml.Tree); ok {
		for name, value := range deps.ToMap() {
			version, ok := value.(string)
			if !ok {
				continue // Skip non-string values like nested tables
			}

			dependencies = append(dependencies, types.Dependency{
				Manager: types.Pyproject,
				Name:    name,
				Version: cleanVersion(version),
				Dev:     false,
				Definition: types.Definition{
					Path:    path,
					RawLine: name + " = \"" + version + "\"",
					Line:    deps.GetPosition(name).Line,
				},
			})
		}
	}

	// github.com/pelletier/go-toml V1 emits struct fields order alphabetically by default.
	// Source: https://github.com/pelletier/go-toml?tab=readme-ov-file#default-struct-fields-order
	// We need to sort the dependencies by line number to keep the original order.

	dependencies = sortDependencies(dependencies)

	return dependencies, nil
}

func sortDependencies(dependencies []types.Dependency) []types.Dependency {
	sort.Slice(dependencies, func(i, j int) bool {
		return dependencies[i].Line < dependencies[j].Line
	})
	return dependencies
}

// Returns the version without any prefix or suffix
func cleanVersion(version string) string {
	version = strings.TrimSpace(version)

	// Check for a range of versions (e.g. ">=1.0.0, <2.0.0")
	if strings.Contains(version, ",") {
		split := strings.Split(version, ",")
		v1 := split[0]
		v2 := split[1]

		if strings.Contains(v1, ">") {
			version = v1
		}

		if strings.Contains(v2, ">") {
			version = v2
		}
	}

	// Remove any comments that might be at the end
	if idx := strings.Index(version, "#"); idx != -1 {
		version = version[:idx]
	}

	version = strings.Trim(version, "^~*><= ")

	return version
}

func (Pyproject) LockfilePath(path string) (string, error) {
	// We don't have a lockfile for pyproject.toml
	return "", nil
}
