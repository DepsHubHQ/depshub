package rules

import (
	"fmt"
	"slices"
	"time"

	"github.com/depshubhq/depshub/pkg/types"
)

const DefaultMaxLibyear = 25.0

type RuleMaxLibyear struct {
	name      string
	level     types.Level
	supported []types.ManagerType
	value     float64
}

func NewRuleMaxLibyear() *RuleMaxLibyear {
	return &RuleMaxLibyear{
		name:      "max-libyear",
		level:     types.LevelError,
		supported: []types.ManagerType{types.Npm, types.Go, types.Cargo, types.Pip, types.Hex},
		value:     DefaultMaxLibyear,
	}
}

func (r RuleMaxLibyear) GetMessage() string {
	return "The total libyear of all dependencies is too high"
}

func (r RuleMaxLibyear) GetName() string {
	return r.name
}

func (r RuleMaxLibyear) GetLevel() types.Level {
	return r.level
}

func (r *RuleMaxLibyear) SetLevel(level types.Level) {
	r.level = level
}

func (r *RuleMaxLibyear) SetValue(value any) error {
	if v, ok := value.(float64); ok {
		r.value = v
		return nil
	}
	return types.ErrInvalidRuleValue
}

func (r RuleMaxLibyear) IsSupported(t types.ManagerType) bool {
	return slices.Contains(r.supported, t)
}

func (r *RuleMaxLibyear) Reset() {
	*r = *NewRuleMaxLibyear()
}

func (r RuleMaxLibyear) Check(manifests []types.Manifest, info types.PackagesInfo, c types.Config) ([]types.Mistake, error) {
	mistakes := []types.Mistake{}

	totalLibyear := 0.0
	topPackageContributors := make(map[string]float64)

	for _, manifest := range manifests {
		if !r.IsSupported(manifest.Manager) {
			continue
		}

		for _, dep := range manifest.Dependencies {
			err := c.Apply(manifest.Path, dep.Name, &r)

			if err != nil {
				return nil, err
			}

			if pkg, ok := info[dep.Name]; ok {
				if t, ok := pkg.Time[dep.Version]; ok {
					if t.IsZero() {
						continue
					}

					diff := time.Since(t)
					diffHours := diff.Abs().Hours()
					totalLibyear += diffHours / (365 * 24)

					topPackageContributors[dep.Name] += diffHours / (365 * 24)
				}
			}
		}
	}

	if totalLibyear > r.value {
		message := fmt.Sprintf("The total libyear of all dependencies is too high.\n Allowed libyear: %.2f. Total libyear: %.2f", r.value, totalLibyear)

		message += "\n\nTop outdated packages:"
		for pkg, libyear := range topPackageContributors {
			message += fmt.Sprintf("\n%s: %.2f", pkg, libyear)
		}

		mistakes = append(mistakes, types.Mistake{
			Rule: r,
			Definitions: []types.Definition{{
				Path: message,
			}},
		})
	}

	return mistakes, nil
}
