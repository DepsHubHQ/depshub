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

func (r RuleNoPreRelease) Check(manifests []types.Manifest, info types.PackagesInfo) (mistakes []Mistake, err error) {
	for _, manifest := range manifests {
		for _, dep := range manifest.Dependencies {
			version := dep.Version

			if strings.Contains(version, "alpha") || strings.Contains(version, "beta") || strings.Contains(version, "rc") {
				mistakes = append(mistakes, Mistake{
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
