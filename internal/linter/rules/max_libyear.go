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
	level     Level
	supported []types.ManagerType
	value     float64
}

func NewRuleMaxLibyear() *RuleMaxLibyear {
	return &RuleMaxLibyear{
		name:      "max-libyear",
		level:     LevelError,
		supported: []types.ManagerType{types.Npm, types.Go},
		value:     DefaultMaxLibyear,
	}
}

func (r RuleMaxLibyear) GetMessage() string {
	return "The total libyear of all dependencies is too high"
}

func (r RuleMaxLibyear) GetName() string {
	return r.name
}

func (r RuleMaxLibyear) GetLevel() Level {
	return r.level
}

func (r *RuleMaxLibyear) SetLevel(level Level) {
	r.level = level
}

func (r *RuleMaxLibyear) SetValue(value any) error {
	if v, ok := value.(float64); ok {
		r.value = v
		return nil
	}
	return ErrInvalidRuleValue
}

func (r RuleMaxLibyear) IsSupported(t types.ManagerType) bool {
	return slices.Contains(r.supported, t)
}

func (r RuleMaxLibyear) Check(manifests []types.Manifest, info types.PackagesInfo) ([]Mistake, error) {
	mistakes := []Mistake{}

	totalLibyear := 0.0

	for _, manifest := range manifests {
		for _, dep := range manifest.Dependencies {
			if pkg, ok := info[dep.Name]; ok {
				if t, ok := pkg.Time[dep.Version]; ok {
					diff := time.Since(t)
					diffHours := diff.Abs().Hours()
					totalLibyear += diffHours / (365 * 24)
				}
			}
		}
	}

	if totalLibyear > r.value {
		mistakes = append(mistakes, Mistake{
			Rule: &r,
			Definitions: []types.Definition{{
				Path: fmt.Sprintf("Allowed libyear: %.2f. Total libyear: %.2f", r.value, totalLibyear),
			}},
		})
	}

	return mistakes, nil
}
