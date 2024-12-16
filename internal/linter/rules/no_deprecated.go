package rules

import (
	"slices"

	"github.com/depshubhq/depshub/pkg/types"
)

type RuleNoDeprecated struct {
	name      string
	level     Level
	supported []types.ManagerType
}

func NewRuleNoDeprecated() *RuleNoDeprecated {
	return &RuleNoDeprecated{
		name:      "no-deprecated",
		level:     LevelError,
		supported: []types.ManagerType{types.Npm, types.Go, types.Cargo},
	}
}

func (r RuleNoDeprecated) GetMessage() string {
	return "Disallow the use of deprecated package versions"
}

func (r RuleNoDeprecated) GetName() string {
	return r.name
}

func (r RuleNoDeprecated) GetLevel() Level {
	return r.level
}

func (r *RuleNoDeprecated) SetLevel(level Level) {
	r.level = level
}

func (r *RuleNoDeprecated) SetValue(value any) error {
	return nil
}

func (r RuleNoDeprecated) IsSupported(t types.ManagerType) bool {
	return slices.Contains(r.supported, t)
}

func (r RuleNoDeprecated) Check(manifests []types.Manifest, info types.PackagesInfo) ([]Mistake, error) {
	mistakes := []Mistake{}

	for _, manifest := range manifests {
		if !r.IsSupported(manifest.Manager) {
			continue
		}

		for _, dep := range manifest.Dependencies {
			if pkg, ok := info[dep.Name]; ok {
				for _, version := range pkg.Versions {
					if version.Version == dep.Version && version.Deprecated != "" {
						mistakes = append(mistakes, Mistake{
							Rule: NewRuleNoDeprecated(),
							Definitions: []types.Definition{
								dep.Definition,
							},
						})
					}
				}
			}
		}
	}

	return mistakes, nil
}
