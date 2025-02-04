package rules

import (
	"slices"
	"strings"

	"github.com/depshubhq/depshub/pkg/types"
)

type RuleNoPreRelease struct {
	name      string
	level     types.Level
	supported []types.ManagerType
}

func NewRuleNoPreRelease() *RuleNoPreRelease {
	return &RuleNoPreRelease{
		name:      "no-pre-release",
		level:     types.LevelError,
		supported: []types.ManagerType{types.Npm, types.Go, types.Cargo, types.Pip, types.Hex, types.Pyproject, types.Maven},
	}
}

func (r RuleNoPreRelease) GetMessage() string {
	return `Disallow the use of "alpha", "beta", "rc", etc. version tags`
}

func (r RuleNoPreRelease) GetName() string {
	return r.name
}

func (r RuleNoPreRelease) GetLevel() types.Level {
	return r.level
}

func (r *RuleNoPreRelease) SetLevel(level types.Level) {
	r.level = level
}

func (r *RuleNoPreRelease) SetValue(value any) error {
	return nil
}

func (r RuleNoPreRelease) IsSupported(t types.ManagerType) bool {
	return slices.Contains(r.supported, t)
}

func (r *RuleNoPreRelease) Reset() {
	*r = *NewRuleNoPreRelease()
}

func (r RuleNoPreRelease) Check(manifests []types.Manifest, info types.PackagesInfo, c types.Config) (mistakes []types.Mistake, err error) {
	for _, manifest := range manifests {
		if !r.IsSupported(manifest.Manager) {
			continue
		}

		for _, dep := range manifest.Dependencies {
			err := c.Apply(manifest.Path, dep.Name, &r)

			if err != nil {
				return nil, err
			}

			version := dep.Version

			if strings.Contains(version, "alpha") || strings.Contains(version, "beta") || strings.Contains(version, "rc") {
				mistakes = append(mistakes, types.Mistake{
					Rule: r,
					Definitions: []types.Definition{
						dep.Definition,
					},
				})
			}
		}
	}

	return mistakes, nil
}
