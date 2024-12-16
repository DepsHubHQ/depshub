package rules

import (
	"slices"

	"github.com/depshubhq/depshub/pkg/types"
)

type RuleSorted struct {
	name      string
	level     Level
	supported []types.ManagerType
}

func NewRuleSorted() *RuleSorted {
	return &RuleSorted{
		name:      "sorted",
		level:     LevelError,
		supported: []types.ManagerType{types.Npm, types.Go, types.Cargo},
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

func (r *RuleSorted) SetLevel(level Level) {
	r.level = level
}

func (r *RuleSorted) SetValue(value any) error {
	return nil
}

func (r RuleSorted) IsSupported(t types.ManagerType) bool {
	return slices.Contains(r.supported, t)
}

func (r RuleSorted) Check(manifests []types.Manifest, info types.PackagesInfo) (mistakes []Mistake, err error) {
	for _, manifest := range manifests {
		if !r.IsSupported(manifest.Manager) {
			continue
		}

		deps := manifest.Dependencies

		// Check if dependencies are ordered
		for i := 0; i < len(deps)-1; i++ {
			if deps[i].Name > deps[i+1].Name {
				// Make sure that we don't compare dependencies and devDependencies
				if deps[i].Dev != deps[i+1].Dev {
					continue
				}

				mistakes = append(mistakes, Mistake{
					Rule:        NewRuleSorted(),
					Definitions: []types.Definition{deps[i].Definition},
				})
			}
		}
	}

	return mistakes, nil
}
