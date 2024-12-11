package rules

import (
	"slices"

	"github.com/depshubhq/depshub/pkg/types"
)

const MinWeeklyDownloads = 1000

type RuleMinWeeklyDownloads struct {
	name      string
	level     Level
	supported []types.ManagerType
}

func NewRuleMinWeeklyDownloads() *RuleMinWeeklyDownloads {
	return &RuleMinWeeklyDownloads{
		name:      "min-weekly-downloads",
		level:     LevelError,
		supported: []types.ManagerType{types.Npm},
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

func (r RuleMinWeeklyDownloads) IsSupported(t types.ManagerType) bool {
	return slices.Contains(r.supported, t)
}

func (r RuleMinWeeklyDownloads) Check(manifests []types.Manifest, info types.PackagesInfo) ([]Mistake, error) {
	mistakes := []Mistake{}

	for _, manifest := range manifests {
		for _, dep := range manifest.Dependencies {
			if pkg, ok := info[dep.Name]; ok {
				weeklyDownloads := 0

				for _, download := range pkg.Downloads {
					weeklyDownloads += download.Downloads
				}

				if weeklyDownloads < MinWeeklyDownloads {
					mistakes = append(mistakes, Mistake{
						Rule:        &r,
						Definitions: []types.Definition{dep.Definition},
					})
				}
			}
		}
	}

	return mistakes, nil
}
