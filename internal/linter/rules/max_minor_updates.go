package rules

import (
	"slices"

	"github.com/depshubhq/depshub/pkg/types"
)

const MaxMinorUpdatesPercent = 40.0

type RuleMaxMinorUpdates struct {
	name      string
	level     Level
	supported []types.ManagerType
}

func NewRuleMaxMinorUpdates() *RuleMaxMinorUpdates {
	return &RuleMaxMinorUpdates{
		name:      "max-minor-updates",
		level:     LevelError,
		supported: []types.ManagerType{types.Npm, types.Go},
	}
}

func (r RuleMaxMinorUpdates) GetMessage() string {
	return "The total number of minor updates is too high"
}

func (r RuleMaxMinorUpdates) GetName() string {
	return r.name
}

func (r RuleMaxMinorUpdates) GetLevel() Level {
	return r.level
}

func (r *RuleMaxMinorUpdates) SetLevel(level Level) {
	r.level = level
}

func (r RuleMaxMinorUpdates) IsSupported(t types.ManagerType) bool {
	return slices.Contains(r.supported, t)
}

func (r RuleMaxMinorUpdates) Check(manifests []types.Manifest, info types.PackagesInfo) ([]Mistake, error) {
	mistakes := []Mistake{}
	definitions := []types.Definition{}
	totalDependencies := 0

	for _, manifest := range manifests {
		for _, dep := range manifest.Dependencies {
			if pkg, ok := info[dep.Name]; ok {
				totalDependencies++
				major, minor, patch := parseVersion(dep.Version)

				for v := range pkg.Versions {
					ma, mi, p := parseVersion(v)

					if mi > minor && ma == major && p == patch {
						definitions = append(definitions, dep.Definition)
						break
					}
				}
			}
		}
	}

	if totalDependencies == 0 {
		return mistakes, nil
	}

	if float64(len(definitions))/float64(totalDependencies)*100 > MaxMinorUpdatesPercent {
		mistakes = append(mistakes, Mistake{
			Rule:        &r,
			Definitions: definitions,
		})
	}

	return mistakes, nil
}
