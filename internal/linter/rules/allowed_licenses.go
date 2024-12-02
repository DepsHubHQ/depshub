package rules

import (
	"slices"

	"github.com/depshubhq/depshub/pkg/types"
)

var AllowedLicenses = []string{
	// No license is also allowed
	"",
	"MIT",
	"Apache-2.0",
}

type RuleAllowedLicenses struct {
	name  string
	level Level
}

func NewRuleAllowedLicenses() RuleAllowedLicenses {
	return RuleAllowedLicenses{
		name:  "allowed-licenses",
		level: LevelError,
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

func (r RuleAllowedLicenses) Check(manifests []types.Manifest, info PackagesInfo) ([]Mistake, error) {
	mistakes := []Mistake{}

	for _, manifest := range manifests {
		for _, dep := range manifest.Dependencies {
			if pkg, ok := info[dep.Name]; ok {
				if !slices.Contains(AllowedLicenses, pkg.License) {
					mistakes = append(mistakes, Mistake{
						Rule:        r,
						Definitions: []types.Definition{dep.Definition},
					})
				}
			}
		}
	}

	return mistakes, nil
}
