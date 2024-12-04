package rules

import (
	"fmt"
	"time"

	"github.com/depshubhq/depshub/pkg/types"
)

const MaxLibyear = 25.0

type RuleMaxLibyear struct {
	name  string
	level Level
}

func NewRuleMaxLibyear() RuleMaxLibyear {
	return RuleMaxLibyear{
		name:  "max-libyear",
		level: LevelError,
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

func (r RuleMaxLibyear) Check(manifests []types.Manifest, info PackagesInfo) ([]Mistake, error) {
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

	if totalLibyear > MaxLibyear {
		mistakes = append(mistakes, Mistake{
			Rule: r,
			Definitions: []types.Definition{{
				Path: fmt.Sprintf("Allowed libyear: %.2f. Total libyear: %.2f", MaxLibyear, totalLibyear),
			}},
		})
	}

	return mistakes, nil
}