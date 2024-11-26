package rules

import (
	"github.com/depshubhq/depshub/pkg/types"
)

type Rule001Order struct {
	name string
}

func NewRule001Order() Rule001Order {
	return Rule001Order{name: "001_order"}
}

func (r Rule001Order) GetMessage() string {
	return "All the dependencies should be ordered."
}

func (r Rule001Order) GetName() string {
	return r.name
}

func (r Rule001Order) Check(manifests []types.Manifest) (mistakes []Mistake, err error) {
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
						Rule: r,
						Path: deps[i].Name,
						Line: 12,
					})
				}
			}
		}
	}

	return mistakes, nil
}
