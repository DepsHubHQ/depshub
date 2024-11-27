package rules

import (
	"strings"

	"github.com/depshubhq/depshub/pkg/types"
)

type RuleNoPreRelease struct {
	name  string
	level Level
}

func NewRuleNoPreRelease() RuleNoPreRelease {
	return RuleNoPreRelease{
		name:  "no-pre-release",
		level: LevelError,
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

func (r RuleNoPreRelease) Check(manifests []types.Manifest) (mistakes []Mistake, err error) {
	for _, manifest := range manifests {
		for _, deps := range [][]types.Dependency{
			manifest.Dependencies,
		} {
			for i := 0; i < len(deps)-1; i++ {
				version := deps[i].Version

				if strings.Contains(version, "alpha") || strings.Contains(version, "beta") || strings.Contains(version, "rc") {
					mistakes = append(mistakes, Mistake{
						Rule:       r,
						Path:       manifest.Path,
						Definition: deps[i+1].Definition,
					})
				}
			}
		}
	}

	return mistakes, nil
}
