package linter

import (
	"fmt"

	"github.com/depshubhq/depshub/internal/linter/rules"
	"github.com/depshubhq/depshub/pkg/manager"
)

type Linter struct {
	rules []rules.Rule
}

func New() Linter {
	return Linter{
		rules: []rules.Rule{
			rules.NewRuleSorted(),
			rules.NewRuleNoAnyTag(),
			rules.NewRuleNoDuplicates(),
			rules.NewRuleNoUnstable(),
			rules.NewRuleNoPreRelease(),
			rules.NewRuleLockfile(),
			rules.NewRuleNoMultipleVersions(),
			rules.NewRuleMaxPackageAge(),
			rules.NewRuleNoDeprecated(),
			rules.NewRuleAllowedLicenses(),
		},
	}
}

func (l Linter) Run(path string) (mistakes []rules.Mistake, err error) {
	scanner := manager.New()
	manifests, err := scanner.Scan(path)
	if err != nil {
		return nil, fmt.Errorf("failed to scan manifests: %w", err)
	}

	uniqueDependencies := scanner.UniqueDependencies(manifests)

	packagesData, err := manager.NewFetcher().Fetch(uniqueDependencies)

	// Run all rules
	for _, rule := range l.rules {
		m, err := rule.Check(manifests, packagesData)
		if err != nil {
			return nil, fmt.Errorf("rule check failed: %w", err)
		}
		mistakes = append(mistakes, m...)
	}

	return mistakes, nil
}
