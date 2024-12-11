package rules

import (
	"slices"
	"strings"

	"github.com/depshubhq/depshub/pkg/types"
)

type RuleNoPreRelease struct {
	name      string
	level     Level
	supported []types.ManagerType
}

func NewRuleNoPreRelease() *RuleNoPreRelease {
	return &RuleNoPreRelease{
		name:      "no-pre-release",
		level:     LevelError,
		supported: []types.ManagerType{types.Npm, types.Go},
	}
}

func (r RuleNoPreRelease) GetMessage() string {
	return `Disallow the use of "alpha", "beta", "rc", etc. version tags`
}

func (r RuleNoPreRelease) GetName() string {
	return r.name
}

func (r RuleNoPreRelease) GetLevel() Level {
	return r.level
}

func (r *RuleNoPreRelease) SetLevel(level Level) {
	r.level = level
}

func (r *RuleNoPreRelease) SetValue(value any) error {
	return nil
}

func (r RuleNoPreRelease) IsSupported(t types.ManagerType) bool {
	return slices.Contains(r.supported, t)
}

func (r RuleNoPreRelease) Check(manifests []types.Manifest, info types.PackagesInfo) (mistakes []Mistake, err error) {
	for _, manifest := range manifests {
		for _, dep := range manifest.Dependencies {
			version := dep.Version

			if strings.Contains(version, "alpha") || strings.Contains(version, "beta") || strings.Contains(version, "rc") {
				mistakes = append(mistakes, Mistake{
					Rule: NewRuleNoPreRelease(),
					Definitions: []types.Definition{
						dep.Definition,
					},
				})
			}
		}
	}

	return mistakes, nil
}
