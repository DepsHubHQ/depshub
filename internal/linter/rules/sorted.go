package rules

import (
	"github.com/depshubhq/depshub/pkg/types"
)

type RuleSorted struct {
	name string
}

func NewRuleSorted() RuleSorted {
	return RuleSorted{name: "sorted"}
}

func (r RuleSorted) GetMessage() string {
	return "All the dependencies should be ordered alphabetically"
}

func (r RuleSorted) GetName() string {
	return r.name
}

func (r RuleSorted) Check(manifests []types.Manifest) (mistakes []Mistake, err error) {
	for _, manifest := range manifests {
		// Check each dependency section
		for _, deps := range [][]types.Dependency{
			manifest.Dependencies,
		} {
			// Check if dependencies are ordered
			for i := 0; i < len(deps)-1; i++ {
				if deps[i].Name > deps[i+1].Name {
					// Make sure that we don't compare dependencies and devDependencies
					if deps[i].Dev != deps[i+1].Dev {
						continue
					}

					mistakes = append(mistakes, Mistake{
						Rule:       r,
						Path:       manifest.Path,
						Definition: deps[i+1].Definition,
					})
				}
			}
		}
	}

	return mistakes, nil
}
