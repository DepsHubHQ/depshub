package rules

import (
	"github.com/depshubhq/depshub/pkg/types"
)

type RuleNoDeprecated struct {
	name  string
	level Level
}

func NewRuleNoDeprecated() RuleNoDeprecated {
	return RuleNoDeprecated{
		name:  "no-deprecated",
		level: LevelError,
	}
}

func (r RuleNoDeprecated) GetMessage() string {
	return "Disallow the use of deprecated package versions"
}

func (r RuleNoDeprecated) GetName() string {
	return r.name
}

func (r RuleNoDeprecated) GetLevel() Level {
	return r.level
}

func (r RuleNoDeprecated) Check(manifests []types.Manifest, info PackagesInfo) ([]Mistake, error) {
	mistakes := []Mistake{}

	for _, manifest := range manifests {
		for _, dep := range manifest.Dependencies {
			if pkg, ok := info[dep.Name]; ok {
				for _, version := range pkg.Versions {
					if version.Version == dep.Version && version.Deprecated != "" {
						mistakes = append(mistakes, Mistake{
							Rule: r,
							Definitions: []types.Definition{
								dep.Definition,
							},
						})
					}
				}
			}
		}
	}

	return mistakes, nil
}
