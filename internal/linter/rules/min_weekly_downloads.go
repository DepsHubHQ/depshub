package rules

import (
	"slices"

	"github.com/depshubhq/depshub/pkg/types"
)

const DefaultMinWeeklyDownloads = 1000

type RuleMinWeeklyDownloads struct {
	name      string
	level     Level
	supported []types.ManagerType
	value     int
}

func NewRuleMinWeeklyDownloads() *RuleMinWeeklyDownloads {
	return &RuleMinWeeklyDownloads{
		name:      "min-weekly-downloads",
		level:     LevelError,
		supported: []types.ManagerType{types.Npm, types.Cargo},
		value:     DefaultMinWeeklyDownloads,
	}
}

func (r RuleMinWeeklyDownloads) GetMessage() string {
	return "Minimum weekly downloads not met"
}

func (r RuleMinWeeklyDownloads) GetName() string {
	return r.name
}

func (r RuleMinWeeklyDownloads) GetLevel() Level {
	return r.level
}

func (r *RuleMinWeeklyDownloads) SetLevel(level Level) {
	r.level = level
}

func (r *RuleMinWeeklyDownloads) SetValue(value any) error {
	if v, ok := value.(int); ok {
		r.value = v
		return nil
	}
	return ErrInvalidRuleValue
}

func (r RuleMinWeeklyDownloads) IsSupported(t types.ManagerType) bool {
	return slices.Contains(r.supported, t)
}

func (r RuleMinWeeklyDownloads) Check(manifests []types.Manifest, info types.PackagesInfo) ([]Mistake, error) {
	mistakes := []Mistake{}

	for _, manifest := range manifests {
		if !r.IsSupported(manifest.Manager) {
			continue
		}

		for _, dep := range manifest.Dependencies {
			if pkg, ok := info[dep.Name]; ok {
				weeklyDownloads := 0

				for _, download := range pkg.Downloads {
					weeklyDownloads += download.Downloads
				}

				if weeklyDownloads < r.value {
					mistakes = append(mistakes, Mistake{
						Rule:        NewRuleMinWeeklyDownloads(),
						Definitions: []types.Definition{dep.Definition},
					})
				}
			}
		}
	}

	return mistakes, nil
}
