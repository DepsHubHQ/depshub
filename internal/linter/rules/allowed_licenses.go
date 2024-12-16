package rules

import (
	"slices"

	"github.com/depshubhq/depshub/pkg/types"
)

var DefaultAllowedLicenses = []string{"", "MIT", "Apache-2.0"}

type RuleAllowedLicenses struct {
	name      string
	level     Level
	supported []types.ManagerType
	value     []string
}

func NewRuleAllowedLicenses() *RuleAllowedLicenses {
	return &RuleAllowedLicenses{
		name:      "allowed-licenses",
		level:     LevelError,
		supported: []types.ManagerType{types.Npm, types.Go, types.Cargo},
		value:     DefaultAllowedLicenses,
	}
}

func (r RuleAllowedLicenses) GetMessage() string {
	return `The license of the package is not allowed.`
}

func (r RuleAllowedLicenses) GetName() string {
	return r.name
}

func (r RuleAllowedLicenses) GetLevel() Level {
	return r.level
}

func (r RuleAllowedLicenses) IsSupported(t types.ManagerType) bool {
	return slices.Contains(r.supported, t)
}

func (r *RuleAllowedLicenses) SetLevel(level Level) {
	r.level = level
}

func (r *RuleAllowedLicenses) SetValue(value any) error {
	if v, ok := value.([]string); ok {
		r.value = v
		return nil
	}
	return ErrInvalidRuleValue
}

func (r RuleAllowedLicenses) Check(manifests []types.Manifest, info types.PackagesInfo) ([]Mistake, error) {
	mistakes := []Mistake{}

	for _, manifest := range manifests {
		if !r.IsSupported(manifest.Manager) {
			continue
		}

		for _, dep := range manifest.Dependencies {
			if pkg, ok := info[dep.Name]; ok {
				if !slices.Contains(r.value, pkg.License) {
					mistakes = append(mistakes, Mistake{
						Rule:        NewRuleAllowedLicenses(),
						Definitions: []types.Definition{dep.Definition},
					})
				}
			}
		}
	}

	return mistakes, nil
}
