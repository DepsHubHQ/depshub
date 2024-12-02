package rules

import (
	"github.com/depshubhq/depshub/pkg/types"
)

type RuleLockfile struct {
	name  string
	level Level
}

func NewRuleLockfile() RuleLockfile {
	return RuleLockfile{
		name:  "lockfile",
		level: LevelError,
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

func (r RuleLockfile) Check(manifests []types.Manifest, info PackagesInfo) (mistakes []Mistake, err error) {
	for _, manifest := range manifests {
		if manifest.Lockfile == nil {
			mistakes = append(mistakes, Mistake{
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
