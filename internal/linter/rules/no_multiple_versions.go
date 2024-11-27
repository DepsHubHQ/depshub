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
	// Map to store the dependencies and their versions
	dependenciesMap := make(map[string]string)

	for _, manifest := range manifests {
		for _, deps := range [][]types.Dependency{
			manifest.Dependencies,
		} {
			for i := 0; i < len(deps)-1; i++ {
				dep := deps[i]

				if version, ok := dependenciesMap[dep.Name]; ok && version != dep.Version {
					mistakes = append(mistakes, Mistake{
						Rule:       r,
						Path:       manifest.Path,
						Definition: &deps[i+1].Definition,
					})
				}

				dependenciesMap[dep.Name] = dep.Version
			}
		}
	}

	return mistakes, nil
}
