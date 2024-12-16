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
	level     Level
	supported []types.ManagerType
}

func NewRuleNoUnstable() *RuleNoUnstable {
	return &RuleNoUnstable{
		name:      "no-unstable",
		level:     LevelError,
		supported: []types.ManagerType{types.Npm, types.Go, types.Cargo, types.Cargo},
	}
}

func (r RuleNoUnstable) GetMessage() string {
	return `Disallow the use of unstable versions (< 1.0.0)`
}

func (r RuleNoUnstable) GetName() string {
	return r.name
}

func (r RuleNoUnstable) GetLevel() Level {
	return r.level
}

func (r *RuleNoUnstable) SetLevel(level Level) {
	r.level = level
}

func (r *RuleNoUnstable) SetValue(value any) error {
	return nil
}

func (r RuleNoUnstable) IsSupported(t types.ManagerType) bool {
	return slices.Contains(r.supported, t)
}

func (r RuleNoUnstable) Check(manifests []types.Manifest, info types.PackagesInfo) (mistakes []Mistake, err error) {
	for _, manifest := range manifests {
		if !r.IsSupported(manifest.Manager) {
			continue
		}
		for _, dep := range manifest.Dependencies {
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
				mistakes = append(mistakes, Mistake{
					Rule:        NewRuleNoUnstable(),
					Definitions: []types.Definition{dep.Definition},
				})
			}
		}
	}

	return mistakes, nil
}
