package rules

import (
	"slices"

	"github.com/depshubhq/depshub/pkg/types"
)

const DefaultMaxMinorUpdatesPercent = 40.0

type RuleMaxMinorUpdates struct {
	name      string
	level     types.Level
	supported []types.ManagerType
	value     float64
}

func NewRuleMaxMinorUpdates() *RuleMaxMinorUpdates {
	return &RuleMaxMinorUpdates{
		name:      "max-minor-updates",
		level:     types.LevelError,
		supported: []types.ManagerType{types.Npm, types.Go, types.Cargo, types.Pip, types.Hex, types.Pyproject},
		value:     DefaultMaxMinorUpdatesPercent,
	}
}

func (r RuleMaxMinorUpdates) GetMessage() string {
	return "The total number of minor updates is too high"
}

func (r RuleMaxMinorUpdates) GetName() string {
	return r.name
}

func (r RuleMaxMinorUpdates) GetLevel() types.Level {
	return r.level
}

func (r *RuleMaxMinorUpdates) SetLevel(level types.Level) {
	r.level = level
}

func (r *RuleMaxMinorUpdates) SetValue(value any) error {
	if v, ok := value.(float64); ok {
		r.value = v
		return nil
	}
	return types.ErrInvalidRuleValue
}

func (r RuleMaxMinorUpdates) IsSupported(t types.ManagerType) bool {
	return slices.Contains(r.supported, t)
}

func (r *RuleMaxMinorUpdates) Reset() {
	*r = *NewRuleMaxMinorUpdates()
}

func (r RuleMaxMinorUpdates) Check(manifests []types.Manifest, info types.PackagesInfo, c types.Config) ([]types.Mistake, error) {
	mistakes := []types.Mistake{}
	definitions := []types.Definition{}
	totalDependencies := 0

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

	if float64(len(definitions))/float64(totalDependencies)*100 > r.value {
		mistakes = append(mistakes, types.Mistake{
			Rule:        r,
			Definitions: definitions,
		})
	}

	return mistakes, nil
}
