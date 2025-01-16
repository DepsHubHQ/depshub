package rules

import (
	"slices"

	"github.com/depshubhq/depshub/pkg/types"
)

type RuleNoDuplicates struct {
	name      string
	level     types.Level
	supported []types.ManagerType
}

func NewRuleNoDuplicates() *RuleNoDuplicates {
	return &RuleNoDuplicates{
		name:      "no-duplicates",
		level:     types.LevelError,
		supported: []types.ManagerType{types.Npm, types.Go, types.Cargo, types.Pip},
	}
}

func (r RuleNoDuplicates) GetMessage() string {
	return `Disallow the same package to be listed multiple times`
}

func (r RuleNoDuplicates) GetName() string {
	return r.name
}

func (r RuleNoDuplicates) GetLevel() types.Level {
	return r.level
}

func (r *RuleNoDuplicates) SetLevel(level types.Level) {
	r.level = level
}

func (r *RuleNoDuplicates) SetValue(value any) error {
	return nil
}

func (r RuleNoDuplicates) IsSupported(t types.ManagerType) bool {
	return slices.Contains(r.supported, t)
}

func (r *RuleNoDuplicates) Reset() {
	*r = *NewRuleNoDuplicates()
}

func (r RuleNoDuplicates) Check(manifests []types.Manifest, info types.PackagesInfo, c types.Config) (mistakes []types.Mistake, err error) {
	for _, manifest := range manifests {
		if !r.IsSupported(manifest.Manager) {
			continue
		}

		deps := manifest.Dependencies

		for i := 0; i < len(deps)-1; i++ {
			for j := i + 1; j < len(deps); j++ {
				err := c.Apply(manifest.Path, deps[j].Name, &r)

				if err != nil {
					return nil, err
				}

				if deps[i].Name == deps[j].Name {
					mistakes = append(mistakes, types.Mistake{
						Rule:        r,
						Definitions: []types.Definition{deps[i].Definition},
					})
				}
			}
		}
	}

	return mistakes, nil
}
