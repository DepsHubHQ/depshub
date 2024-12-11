package rules

import (
	"slices"

	"github.com/depshubhq/depshub/pkg/types"
)

type RuleLockfile struct {
	name      string
	level     Level
	supported []types.ManagerType
}

func NewRuleLockfile() *RuleLockfile {
	return &RuleLockfile{
		name:      "lockfile",
		level:     LevelError,
		supported: []types.ManagerType{types.Npm, types.Go},
	}
}

func (r RuleLockfile) GetMessage() string {
	return "The lockfile should be always present"
}

func (r RuleLockfile) GetName() string {
	return r.name
}

func (r RuleLockfile) GetLevel() Level {
	return r.level
}

func (r *RuleLockfile) SetLevel(level Level) {
	r.level = level
}

func (r *RuleLockfile) SetValue(value any) error {
	return nil
}

func (r RuleLockfile) IsSupported(t types.ManagerType) bool {
	return slices.Contains(r.supported, t)
}

func (r RuleLockfile) Check(manifests []types.Manifest, info types.PackagesInfo) (mistakes []Mistake, err error) {
	for _, manifest := range manifests {
		if manifest.Lockfile == nil {
			mistakes = append(mistakes, Mistake{
				Rule: NewRuleLockfile(),
				Definitions: []types.Definition{
					{
						Path: manifest.Path,
					},
				},
			})
		}
	}

	return mistakes, nil
}
