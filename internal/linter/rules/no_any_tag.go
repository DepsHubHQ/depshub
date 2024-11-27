package rules

import (
	"github.com/depshubhq/depshub/pkg/types"
)

type RuleNoAnyTag struct {
	name  string
	level Level
}

func NewRuleNoAnyTag() RuleNoAnyTag {
	return RuleNoAnyTag{
		name:  "no-any-tag",
		level: LevelWarning,
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

func (r RuleNoAnyTag) Check(manifests []types.Manifest) (mistakes []Mistake, err error) {
	for _, manifest := range manifests {
		for _, deps := range [][]types.Dependency{
			manifest.Dependencies,
		} {
			for i := 0; i < len(deps)-1; i++ {
				if deps[i].Version == "*" || deps[i].Version == "latest" || deps[i].Version == "" {
					mistakes = append(mistakes, Mistake{
						Rule:       r,
						Path:       manifest.Path,
						Definition: &deps[i+1].Definition,
					})
				}
			}
		}
	}

	return mistakes, nil
}
