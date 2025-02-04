package rules

import (
	"slices"

	"github.com/depshubhq/depshub/pkg/types"
)

type RuleLockfile struct {
	name      string
	level     types.Level
	supported []types.ManagerType
}

func NewRuleLockfile() *RuleLockfile {
	return &RuleLockfile{
		name:      "lockfile",
		level:     types.LevelError,
		supported: []types.ManagerType{types.Npm, types.Go, types.Cargo, types.Pip, types.Hex, types.Pyproject, types.Maven, types.Gem},
	}
}

func (r RuleLockfile) GetMessage() string {
	return "The lockfile should be always present"
}

func (r RuleLockfile) GetName() string {
	return r.name
}

func (r RuleLockfile) GetLevel() types.Level {
	return r.level
}

func (r *RuleLockfile) SetLevel(level types.Level) {
	r.level = level
}

func (r *RuleLockfile) SetValue(value any) error {
	return nil
}

func (r *RuleLockfile) Reset() {
	*r = *NewRuleLockfile()
}

func (r RuleLockfile) IsSupported(t types.ManagerType) bool {
	return slices.Contains(r.supported, t)
}

func (r RuleLockfile) Check(manifests []types.Manifest, info types.PackagesInfo, c types.Config) (mistakes []types.Mistake, err error) {
	for _, manifest := range manifests {
		if !r.IsSupported(manifest.Manager) {
			continue
		}
		err := c.Apply(manifest.Path, "", &r)

		if err != nil {
			return nil, err
		}

		if manifest.Lockfile == nil {
			mistakes = append(mistakes, types.Mistake{
				Rule: r,
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
