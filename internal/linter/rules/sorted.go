package rules

import (
	"github.com/depshubhq/depshub/pkg/types"
)

type RuleSorted struct {
	name  string
	level Level
}

func NewRuleSorted() RuleSorted {
	return RuleSorted{
		name:  "sorted",
		level: LevelError,
	}
}

func (r RuleSorted) GetMessage() string {
	return "All the dependencies should be ordered alphabetically"
}

func (r RuleSorted) GetName() string {
	return r.name
}

func (r RuleSorted) GetLevel() Level {
	return r.level
}

func (r RuleSorted) Check(manifests []types.Manifest) (mistakes []Mistake, err error) {
	for _, manifest := range manifests {
		deps := manifest.Dependencies

		// Check if dependencies are ordered
		for i := 0; i < len(deps)-1; i++ {
			if deps[i].Name > deps[i+1].Name {
				// Make sure that we don't compare dependencies and devDependencies
				if deps[i].Dev != deps[i+1].Dev {
					continue
				}

				mistakes = append(mistakes, Mistake{
					Rule:        r,
					Definitions: []types.Definition{deps[i].Definition},
				})
			}
		}
	}

	return mistakes, nil
}
