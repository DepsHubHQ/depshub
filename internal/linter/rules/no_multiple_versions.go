package rules

import (
	"github.com/depshubhq/depshub/pkg/types"
)

type RuleNoMultipleVersions struct {
	name  string
	level Level
}

func NewRuleNoMultipleVersions() RuleNoMultipleVersions {
	return RuleNoMultipleVersions{
		name:  "no-multiple-versions",
		level: LevelError,
	}
}

func (r RuleNoMultipleVersions) GetMessage() string {
	return "Disallow the use of multiple versions of the same package"
}

func (r RuleNoMultipleVersions) GetName() string {
	return r.name
}

func (r RuleNoMultipleVersions) GetLevel() Level {
	return r.level
}

func (r RuleNoMultipleVersions) Check(manifests []types.Manifest) (mistakes []Mistake, err error) {
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
		for _, dep := range manifest.Dependencies {
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

		mistakes = append(mistakes, Mistake{
			Rule:        r,
			Definitions: definitions,
		})
	}

	return mistakes, nil
}
