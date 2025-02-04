package rules

import (
	"slices"

	"github.com/depshubhq/depshub/pkg/types"
)

const DefaultMaxPatchUpdatesPercent = 60.0

type RuleMaxPatchUpdates struct {
	name      string
	level     types.Level
	supported []types.ManagerType
	value     float64
}

func NewRuleMaxPatchUpdates() *RuleMaxPatchUpdates {
	return &RuleMaxPatchUpdates{
		name:      "max-patch-updates",
		level:     types.LevelError,
		supported: []types.ManagerType{types.Npm, types.Go, types.Cargo, types.Pip, types.Hex, types.Pyproject, types.Maven, types.Gem},
		value:     DefaultMaxPatchUpdatesPercent,
	}
}

func (r RuleMaxPatchUpdates) GetMessage() string {
	return "The total number of patch updates is too high"
}

func (r RuleMaxPatchUpdates) GetName() string {
	return r.name
}

func (r RuleMaxPatchUpdates) GetLevel() types.Level {
	return r.level
}

func (r *RuleMaxPatchUpdates) SetLevel(level types.Level) {
	r.level = level
}

func (r *RuleMaxPatchUpdates) SetValue(value any) error {
	if v, ok := value.(float64); ok {
		r.value = v
		return nil
	}
	return types.ErrInvalidRuleValue
}

func (r RuleMaxPatchUpdates) IsSupported(t types.ManagerType) bool {
	return slices.Contains(r.supported, t)
}

func (r *RuleMaxPatchUpdates) Reset() {
	*r = *NewRuleMaxPatchUpdates()
}

func (r RuleMaxPatchUpdates) Check(manifests []types.Manifest, info types.PackagesInfo, c types.Config) ([]types.Mistake, error) {
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

	if float64(len(definitions))/float64(totalDependencies)*100 > DefaultMaxPatchUpdatesPercent {
		mistakes = append(mistakes, types.Mistake{
			Rule:        r,
			Definitions: definitions,
		})
	}

	return mistakes, nil
}
