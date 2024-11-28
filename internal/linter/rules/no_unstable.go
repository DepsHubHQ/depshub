package rules

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/depshubhq/depshub/pkg/types"
)

type RuleNoUnstable struct {
	name  string
	level Level
}

func NewRuleNoUnstable() RuleNoUnstable {
	return RuleNoUnstable{
		name:  "no-unstable",
		level: LevelError,
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

func (r RuleNoUnstable) Check(manifests []types.Manifest) (mistakes []Mistake, err error) {
	for _, manifest := range manifests {
		for _, dep := range manifest.Dependencies {
			// Define regex pattern for x.x.x where x is one or more digits
			pattern := regexp.MustCompile(`\d+\.\d+\.\d+`)

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

			fmt.Println("package, majorVersion", dep.Name, majorVersion)
			if majorVersion < 1 {
				mistakes = append(mistakes, Mistake{
					Rule:        r,
					Definitions: []types.Definition{dep.Definition},
				})
			}
		}
	}

	return mistakes, nil
}
