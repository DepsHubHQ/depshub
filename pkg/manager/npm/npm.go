package npm

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/depshubhq/depshub/pkg/types"
)

type Npm struct{}

type PackageJSON struct {
	Dependencies    OrderedMap `json:"dependencies"`
	DevDependencies OrderedMap `json:"devDependencies"`
}

func (Npm) Managed(path string) bool {
	fileName := filepath.Base(path)
	return fileName == "package.json"
}

func (Npm) Dependencies(path string) ([]types.Dependency, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var packageJSON PackageJSON
	if err := json.Unmarshal(file, &packageJSON); err != nil {
		return nil, err
	}

	var dependencies []types.Dependency

	// Add regular dependencies in order
	for _, name := range packageJSON.Dependencies.Order {
		dependencies = append(dependencies, types.Dependency{
			Name:    name,
			Version: packageJSON.Dependencies.Values[name],
			Dev:     false,
			Definition: types.Definition{
				RawLine: packageJSON.Dependencies.RawLines[name],
				Line:    packageJSON.Dependencies.LineNums[name],
			},
		})
	}

	// Add dev dependencies in order
	for _, name := range packageJSON.DevDependencies.Order {
		dependencies = append(dependencies, types.Dependency{
			Name:    name,
			Version: packageJSON.DevDependencies.Values[name],
			Dev:     true,
			Definition: types.Definition{
				RawLine: packageJSON.Dependencies.RawLines[name],
				Line:    packageJSON.Dependencies.LineNums[name],
			},
		})
	}

	return dependencies, nil
}
