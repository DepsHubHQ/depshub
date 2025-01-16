package rules

import (
	"slices"

	"github.com/depshubhq/depshub/pkg/types"
)

const DefaultMinWeeklyDownloads = 1000

type RuleMinWeeklyDownloads struct {
	name      string
	level     types.Level
	supported []types.ManagerType
	value     int
}

func NewRuleMinWeeklyDownloads() *RuleMinWeeklyDownloads {
	return &RuleMinWeeklyDownloads{
		name:      "min-weekly-downloads",
		level:     types.LevelError,
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

func (r RuleMinWeeklyDownloads) GetLevel() types.Level {
	return r.level
}

func (r *RuleMinWeeklyDownloads) SetLevel(level types.Level) {
	r.level = level
}

func (r *RuleMinWeeklyDownloads) SetValue(value any) error {
	if v, ok := value.(int); ok {
		r.value = v
		return nil
	}
	return types.ErrInvalidRuleValue
}

func (r RuleMinWeeklyDownloads) IsSupported(t types.ManagerType) bool {
	return slices.Contains(r.supported, t)
}

func (r *RuleMinWeeklyDownloads) Reset() {
	r = NewRuleMinWeeklyDownloads()
}

func (r RuleMinWeeklyDownloads) Check(manifests []types.Manifest, info types.PackagesInfo, c types.Config) ([]types.Mistake, error) {
	mistakes := []types.Mistake{}

	for _, manifest := range manifests {
		if !r.IsSupported(manifest.Manager) {
			continue
		}

		for _, dep := range manifest.Dependencies {
			if pkg, ok := info[dep.Name]; ok {
				err := c.Apply(manifest.Path, dep.Name, &r)

				if err != nil {
					return nil, err
				}

				weeklyDownloads := 0

				for _, download := range pkg.Downloads {
					weeklyDownloads += download.Downloads
				}

				if weeklyDownloads < r.value {
					mistakes = append(mistakes, types.Mistake{
						Rule:        r,
						Definitions: []types.Definition{dep.Definition},
					})
				}
			}
		}
	}

	return mistakes, nil
}
