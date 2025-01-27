package rules

import (
	"slices"

	"github.com/depshubhq/depshub/pkg/types"
)

var DefaultAllowedLicenses = []string{"", "MIT", "Apache-2.0"}

type RuleAllowedLicenses struct {
	name      string
	level     types.Level
	supported []types.ManagerType
	value     []string
}

func NewRuleAllowedLicenses() *RuleAllowedLicenses {
	return &RuleAllowedLicenses{
		name:      "allowed-licenses",
		level:     types.LevelError,
		supported: []types.ManagerType{types.Npm, types.Go, types.Cargo, types.Pip, types.Hex},
		value:     DefaultAllowedLicenses,
	}
}

func (r RuleAllowedLicenses) GetMessage() string {
	return `The license of the package is not allowed.`
}

func (r RuleAllowedLicenses) GetName() string {
	return r.name
}

func (r RuleAllowedLicenses) GetLevel() types.Level {
	return r.level
}

func (r RuleAllowedLicenses) IsSupported(t types.ManagerType) bool {
	return slices.Contains(r.supported, t)
}

func (r *RuleAllowedLicenses) SetLevel(level types.Level) {
	r.level = level
}

func (r *RuleAllowedLicenses) SetValue(value any) error {
	if v, ok := value.([]any); ok {
		for _, i := range v {
			val, ok := i.(string)
			if !ok {
				return types.ErrInvalidRuleValue
			}

			r.value = append(r.value, val)
		}

		return nil
	}
	return types.ErrInvalidRuleValue
}

func (r *RuleAllowedLicenses) Reset() {
	*r = *NewRuleAllowedLicenses()
}

func (r RuleAllowedLicenses) Check(manifests []types.Manifest, info types.PackagesInfo, c types.Config) ([]types.Mistake, error) {
	mistakes := []types.Mistake{}

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
				if !slices.Contains(r.value, pkg.License) {
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
