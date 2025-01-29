package rules

import (
	"slices"

	"github.com/depshubhq/depshub/pkg/types"
)

type RuleNoMultipleVersions struct {
	name      string
	level     types.Level
	supported []types.ManagerType
}

func NewRuleNoMultipleVersions() *RuleNoMultipleVersions {
	return &RuleNoMultipleVersions{
		name:      "no-multiple-versions",
		level:     types.LevelError,
		supported: []types.ManagerType{types.Npm, types.Go, types.Cargo, types.Pip, types.Hex, types.Pyproject, types.Maven},
	}
}

func (r RuleNoMultipleVersions) GetMessage() string {
	return "Disallow the use of multiple versions of the same package"
}

func (r RuleNoMultipleVersions) GetName() string {
	return r.name
}

func (r RuleNoMultipleVersions) GetLevel() types.Level {
	return r.level
}

func (r *RuleNoMultipleVersions) SetLevel(level types.Level) {
	r.level = level
}

func (r *RuleNoMultipleVersions) SetValue(value any) error {
	return nil
}

func (r RuleNoMultipleVersions) IsSupported(t types.ManagerType) bool {
	return slices.Contains(r.supported, t)
}

func (r *RuleNoMultipleVersions) Reset() {
	*r = *NewRuleNoMultipleVersions()
}

func (r RuleNoMultipleVersions) Check(manifests []types.Manifest, info types.PackagesInfo, c types.Config) (mistakes []types.Mistake, err error) {
	type PackageInfo struct {
		Path    string
		Version string
		RawLine string
		Line    int
	}
	// Map to store the dependencies and their versions
	dependenciesMap := make(map[string][]PackageInfo)

	// Collect all dependencies
	for _, manifest := range manifests {
		if !r.IsSupported(manifest.Manager) {
			continue
		}

		for _, dep := range manifest.Dependencies {
			err := c.Apply(manifest.Path, dep.Name, &r)

			if err != nil {
				return nil, err
			}

			// Check if the dependency version is already in the map
			if len(dependenciesMap[dep.Name]) != 0 {
				for _, d := range dependenciesMap[dep.Name] {
					if d.Version != dep.Version {
						dependenciesMap[dep.Name] = append(dependenciesMap[dep.Name], PackageInfo{
							Path:    manifest.Path,
							Version: dep.Version,
							RawLine: dep.RawLine,
							Line:    dep.Line,
						})
					}
				}
			} else {
				dependenciesMap[dep.Name] = append(dependenciesMap[dep.Name], PackageInfo{
					Path:    manifest.Path,
					Version: dep.Version,
					RawLine: dep.RawLine,
					Line:    dep.Line,
				})
			}
		}
	}

	// Check for version conflicts
	for _, deps := range dependenciesMap {
		if len(deps) <= 1 {
			continue
		}

		// Create a single mistake with all definitions where versions differ
		var definitions []types.Definition

		// Add all occurrences with different versions
		for _, dep := range deps {
			definitions = append(definitions, types.Definition{
				Path:    dep.Path,
				RawLine: dep.RawLine,
				Line:    dep.Line,
			})
		}

		mistakes = append(mistakes, types.Mistake{
			Rule:        r,
			Definitions: definitions,
		})
	}

	return mistakes, nil
}
