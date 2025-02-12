package rules

import (
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/depshubhq/depshub/pkg/types"
)

type RuleNoUnstable struct {
	name      string
	level     types.Level
	supported []types.ManagerType
}

func NewRuleNoUnstable() *RuleNoUnstable {
	return &RuleNoUnstable{
		name:      "no-unstable",
		level:     types.LevelError,
		supported: []types.ManagerType{types.Npm, types.Go, types.Cargo, types.Pip, types.Hex, types.Pyproject, types.Maven},
	}
}

func (r RuleNoUnstable) GetMessage() string {
	return `Disallow the use of unstable versions (< 1.0.0)`
}

func (r RuleNoUnstable) GetName() string {
	return r.name
}

func (r RuleNoUnstable) GetLevel() types.Level {
	return r.level
}

func (r *RuleNoUnstable) SetLevel(level types.Level) {
	r.level = level
}

func (r *RuleNoUnstable) SetValue(value any) error {
	return nil
}

func (r RuleNoUnstable) IsSupported(t types.ManagerType) bool {
	return slices.Contains(r.supported, t)
}

func (r *RuleNoUnstable) Reset() {
	*r = *NewRuleNoUnstable()
}

func (r RuleNoUnstable) Check(manifests []types.Manifest, info types.PackagesInfo, c types.Config) (mistakes []types.Mistake, err error) {
	for _, manifest := range manifests {
		if !r.IsSupported(manifest.Manager) {
			continue
		}
		for _, dep := range manifest.Dependencies {
			err := c.Apply(manifest.Path, dep.Name, &r)

			if err != nil {
				return nil, err
			}

			// Define regex pattern for x.x.x or x.x where x is one or more digits
			pattern := regexp.MustCompile(`\d+\.\d+(\.\d+)?`)

			match := pattern.FindString(dep.Version)
			if match == "" {
				continue
			}

			// Split version string into components
			parts := strings.Split(match, ".")

			// Parse major version
			majorVersion, err := strconv.Atoi(parts[0])
			if err != nil {
				continue
			}

			if majorVersion < 1 {
				mistakes = append(mistakes, types.Mistake{
					Rule:        r,
					Definitions: []types.Definition{dep.Definition},
				})
			}
		}
	}

	return mistakes, nil
}
