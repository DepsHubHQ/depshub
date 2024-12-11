package rules

import (
	"slices"

	"github.com/depshubhq/depshub/pkg/types"
)

type RuleNoAnyTag struct {
	name      string
	level     Level
	supported []types.ManagerType
}

func NewRuleNoAnyTag() *RuleNoAnyTag {
	return &RuleNoAnyTag{
		name:      "no-any-tag",
		level:     LevelWarning,
		supported: []types.ManagerType{types.Npm, types.Go},
	}
}

func (r RuleNoAnyTag) GetMessage() string {
	return `Disallow the use of the "any" version tag`
}

func (r RuleNoAnyTag) GetName() string {
	return r.name
}

func (r RuleNoAnyTag) GetLevel() Level {
	return r.level
}

func (r *RuleNoAnyTag) SetLevel(level Level) {
	r.level = level
}

func (r RuleNoAnyTag) IsSupported(t types.ManagerType) bool {
	return slices.Contains(r.supported, t)
}

func (r RuleNoAnyTag) Check(manifests []types.Manifest, info types.PackagesInfo) (mistakes []Mistake, err error) {
	for _, manifest := range manifests {
		for _, dep := range manifest.Dependencies {
			if dep.Version == "*" || dep.Version == "latest" || dep.Version == "" {
				mistakes = append(mistakes, Mistake{
					Rule:        &r,
					Definitions: []types.Definition{dep.Definition},
				})
			}
		}
	}

	return mistakes, nil
}
