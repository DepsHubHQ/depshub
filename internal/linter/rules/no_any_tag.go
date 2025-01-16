package rules

import (
	"slices"

	"github.com/depshubhq/depshub/pkg/types"
)

type RuleNoAnyTag struct {
	name      string
	level     types.Level
	supported []types.ManagerType
}

func NewRuleNoAnyTag() *RuleNoAnyTag {
	return &RuleNoAnyTag{
		name:      "no-any-tag",
		level:     types.LevelWarning,
		supported: []types.ManagerType{types.Npm, types.Go, types.Cargo, types.Pip},
	}
}

func (r RuleNoAnyTag) GetMessage() string {
	return `Disallow the use of the "any" version tag`
}

func (r RuleNoAnyTag) GetName() string {
	return r.name
}

func (r RuleNoAnyTag) GetLevel() types.Level {
	return r.level
}

func (r *RuleNoAnyTag) SetLevel(level types.Level) {
	r.level = level
}

func (r *RuleNoAnyTag) SetValue(value any) error {
	return nil
}

func (r RuleNoAnyTag) IsSupported(t types.ManagerType) bool {
	return slices.Contains(r.supported, t)
}

func (r *RuleNoAnyTag) Reset() {
	*r = *NewRuleNoAnyTag()
}

func (r RuleNoAnyTag) Check(manifests []types.Manifest, info types.PackagesInfo, c types.Config) (mistakes []types.Mistake, err error) {
	for _, manifest := range manifests {
		if !r.IsSupported(manifest.Manager) {
			continue
		}

		for _, dep := range manifest.Dependencies {
			err := c.Apply(manifest.Path, dep.Name, &r)

			if err != nil {
				return nil, err
			}

			if dep.Version == "*" || dep.Version == "latest" || dep.Version == "" {
				mistakes = append(mistakes, types.Mistake{
					Rule:        r,
					Definitions: []types.Definition{dep.Definition},
				})
			}
		}
	}

	return mistakes, nil
}
