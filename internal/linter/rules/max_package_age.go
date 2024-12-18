package rules

import (
	"slices"
	"time"

	"github.com/depshubhq/depshub/pkg/types"
)

const DefaultMaxPackageAge = 36

type RuleMaxPackageAge struct {
	name      string
	level     Level
	supported []types.ManagerType
	value     int
}

func NewRuleMaxPackageAge() *RuleMaxPackageAge {
	return &RuleMaxPackageAge{
		name:      "max-package-age",
		level:     LevelError,
		supported: []types.ManagerType{types.Npm, types.Go, types.Cargo, types.Pip},
		value:     DefaultMaxPackageAge,
	}
}

func (r RuleMaxPackageAge) GetMessage() string {
	return `Disallow the use of any package that is older than a certain age (in months).`
}

func (r RuleMaxPackageAge) GetName() string {
	return r.name
}

func (r RuleMaxPackageAge) GetLevel() Level {
	return r.level
}

func (r *RuleMaxPackageAge) SetLevel(level Level) {
	r.level = level
}

func (r *RuleMaxPackageAge) SetValue(value any) error {
	if v, ok := value.(int); ok {
		r.value = v
		return nil
	}
	return ErrInvalidRuleValue
}

func (r RuleMaxPackageAge) IsSupported(t types.ManagerType) bool {
	return slices.Contains(r.supported, t)
}

func (r RuleMaxPackageAge) Check(manifests []types.Manifest, info types.PackagesInfo) ([]Mistake, error) {
	mistakes := []Mistake{}

	for _, manifest := range manifests {
		if !r.IsSupported(manifest.Manager) {
			continue
		}

		for _, dep := range manifest.Dependencies {

			if pkg, ok := info[dep.Name]; ok {
				for version, t := range pkg.Time {
					if version == dep.Version && t.Before(time.Now().AddDate(0, -r.value, 0)) {
						mistakes = append(mistakes, Mistake{
							Rule: NewRuleMaxPackageAge(),
							Definitions: []types.Definition{
								dep.Definition,
							},
						})
					}
				}
			}
		}
	}

	return mistakes, nil
}
