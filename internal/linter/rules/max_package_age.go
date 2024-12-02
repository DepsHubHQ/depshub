package rules

import (
	"time"

	"github.com/depshubhq/depshub/pkg/types"
)

const MaxPackageAge = 36

type RuleMaxPackageAge struct {
	name  string
	level Level
}

func NewRuleMaxPackageAge() RuleMaxPackageAge {
	return RuleMaxPackageAge{
		name:  "max-package-age",
		level: LevelError,
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

func (r RuleMaxPackageAge) Check(manifests []types.Manifest, info PackagesInfo) (mistakes []Mistake, err error) {
	for _, manifest := range manifests {
		for _, dep := range manifest.Dependencies {

			if pkg, ok := info[dep.Name]; ok {
				for version, t := range pkg.Time {
					if version == dep.Version && t.Before(time.Now().AddDate(0, -MaxPackageAge, 0)) {
						mistakes = append(mistakes, Mistake{
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
