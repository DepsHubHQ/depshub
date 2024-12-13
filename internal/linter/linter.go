package linter

import (
	"fmt"

	"github.com/depshubhq/depshub/internal/linter/rules"
	"github.com/depshubhq/depshub/pkg/manager"
	"github.com/depshubhq/depshub/pkg/sources"
)

type Linter struct {
	rules []rules.Rule
}

func New() Linter {
	return Linter{
		rules: []rules.Rule{
			rules.NewRuleAllowedLicenses(),
			rules.NewRuleLockfile(),
			rules.NewRuleMaxLibyear(),
			rules.NewRuleMaxMajorUpdates(),
			rules.NewRuleMaxMinorUpdates(),
			rules.NewRuleMaxPackageAge(),
			rules.NewRuleMaxPatchUpdates(),
			rules.NewRuleMinWeeklyDownloads(),
			rules.NewRuleNoAnyTag(),
			rules.NewRuleNoDeprecated(),
			rules.NewRuleNoDuplicates(),
			rules.NewRuleNoMultipleVersions(),
			rules.NewRuleNoPreRelease(),
			rules.NewRuleNoUnstable(),
			rules.NewRuleSorted(),
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

	packagesData, err := sources.NewFetcher().Fetch(uniqueDependencies)

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
