package rules

import (
	"github.com/depshubhq/depshub/pkg/types"
)

const MaxMajorUpdatesPercent = 20.0

type RuleMaxMajorUpdates struct {
	name  string
	level Level
}

func NewRuleMaxMajorUpdates() RuleMaxMajorUpdates {
	return RuleMaxMajorUpdates{
		name:  "max-major-updates",
		level: LevelError,
	}
}

func (r RuleMaxMajorUpdates) GetMessage() string {
	return "The total number of major updates is too high"
}

func (r RuleMaxMajorUpdates) GetName() string {
	return r.name
}

func (r RuleMaxMajorUpdates) GetLevel() Level {
	return r.level
}

func (r RuleMaxMajorUpdates) Check(manifests []types.Manifest, info PackagesInfo) ([]Mistake, error) {
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

					if ma > major && p == patch && mi == minor {
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

	if float64(len(definitions))/float64(totalDependencies)*100 > MaxMajorUpdatesPercent {
		mistakes = append(mistakes, Mistake{
			Rule:        r,
			Definitions: definitions,
		})
	}

	return mistakes, nil
}
