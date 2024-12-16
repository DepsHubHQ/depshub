package rules

import (
	"slices"

	"github.com/depshubhq/depshub/pkg/types"
)

const DefaultMaxMajorUpdatesPercent = 20.0

type RuleMaxMajorUpdates struct {
	name      string
	level     Level
	supported []types.ManagerType
	value     float64
}

func NewRuleMaxMajorUpdates() *RuleMaxMajorUpdates {
	return &RuleMaxMajorUpdates{
		name:      "max-major-updates",
		level:     LevelError,
		supported: []types.ManagerType{types.Npm, types.Go},
		value:     DefaultMaxMajorUpdatesPercent,
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

func (r *RuleMaxMajorUpdates) SetLevel(level Level) {
	r.level = level
}

func (r *RuleMaxMajorUpdates) SetValue(value any) error {
	if v, ok := value.(float64); ok {
		r.value = v
		return nil
	}
	return ErrInvalidRuleValue
}

func (r RuleMaxMajorUpdates) IsSupported(t types.ManagerType) bool {
	return slices.Contains(r.supported, t)
}

func (r RuleMaxMajorUpdates) Check(manifests []types.Manifest, info types.PackagesInfo) ([]Mistake, error) {
	mistakes := []Mistake{}
	definitions := []types.Definition{}
	totalDependencies := 0

	for _, manifest := range manifests {
		if !r.IsSupported(manifest.Manager) {
			continue
		}

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

	if float64(len(definitions))/float64(totalDependencies)*100 > r.value {
		mistakes = append(mistakes, Mistake{
			Rule:        NewRuleMaxMajorUpdates(),
			Definitions: definitions,
		})
	}

	return mistakes, nil
}
