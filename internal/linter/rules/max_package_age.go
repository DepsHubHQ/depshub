package rules

import (
	"slices"
	"time"

	"github.com/depshubhq/depshub/pkg/types"
)

const DefaultMaxPackageAge = 36

type RuleMaxPackageAge struct {
	name      string
	level     types.Level
	supported []types.ManagerType
	value     int
}

func NewRuleMaxPackageAge() *RuleMaxPackageAge {
	return &RuleMaxPackageAge{
		name:      "max-package-age",
		level:     types.LevelError,
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

func (r RuleMaxPackageAge) GetLevel() types.Level {
	return r.level
}

func (r *RuleMaxPackageAge) SetLevel(level types.Level) {
	r.level = level
}

func (r *RuleMaxPackageAge) SetValue(value any) error {
	if v, ok := value.(int); ok {
		r.value = v
		return nil
	}
	return types.ErrInvalidRuleValue
}

func (r RuleMaxPackageAge) IsSupported(t types.ManagerType) bool {
	return slices.Contains(r.supported, t)
}

func (r *RuleMaxPackageAge) Reset() {
	*r = *NewRuleMaxPackageAge()
}

func (r RuleMaxPackageAge) Check(manifests []types.Manifest, info types.PackagesInfo, c types.Config) ([]types.Mistake, error) {
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

				for version, t := range pkg.Time {
					if version == dep.Version && t.Before(time.Now().AddDate(0, -r.value, 0)) {
						mistakes = append(mistakes, types.Mistake{
							Rule: r,
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
