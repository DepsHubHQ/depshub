package rules

import (
	"github.com/depshubhq/depshub/pkg/types"
)

const MaxPatchUpdatesPercent = 60.0

type RuleMaxPatchUpdates struct {
	name  string
	level Level
}

func NewRuleMaxPatchUpdates() RuleMaxPatchUpdates {
	return RuleMaxPatchUpdates{
		name:  "max-patch-updates",
		level: LevelError,
	}
}

func (r RuleMaxPatchUpdates) GetMessage() string {
	return "The total number of patch updates is too high"
}

func (r RuleMaxPatchUpdates) GetName() string {
	return r.name
}

func (r RuleMaxPatchUpdates) GetLevel() Level {
	return r.level
}

func (r RuleMaxPatchUpdates) Check(manifests []types.Manifest, info types.PackagesInfo) ([]Mistake, error) {
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

					if p > patch && ma == major && mi == minor {
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

	if float64(len(definitions))/float64(totalDependencies)*100 > MaxPatchUpdatesPercent {
		mistakes = append(mistakes, Mistake{
			Rule:        r,
			Definitions: definitions,
		})
	}

	return mistakes, nil
}
