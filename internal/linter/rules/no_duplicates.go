package rules

import (
	"github.com/depshubhq/depshub/pkg/types"
)

type RuleNoDuplicates struct {
	name  string
	level Level
}

func NewRuleNoDuplicates() RuleNoDuplicates {
	return RuleNoDuplicates{
		name:  "no-duplicates",
		level: LevelError,
	}
}

func (r RuleNoDuplicates) GetMessage() string {
	return `Disallow the same package to be listed multiple times`
}

func (r RuleNoDuplicates) GetName() string {
	return r.name
}

func (r RuleNoDuplicates) GetLevel() Level {
	return r.level
}

func (r RuleNoDuplicates) Check(manifests []types.Manifest) (mistakes []Mistake, err error) {
	for _, manifest := range manifests {
		deps := manifest.Dependencies

		for i := 0; i < len(deps)-1; i++ {
			for j := i + 1; j < len(deps); j++ {
				if deps[i].Name == deps[j].Name {
					mistakes = append(mistakes, Mistake{
						Rule:        r,
						Definitions: []types.Definition{deps[i].Definition},
					})
				}
			}
		}
	}

	return mistakes, nil
}
