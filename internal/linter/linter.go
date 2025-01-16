package linter

import (
	"fmt"

	"github.com/depshubhq/depshub/internal/config"
	"github.com/depshubhq/depshub/internal/linter/rules"
	"github.com/depshubhq/depshub/pkg/manager"
	"github.com/depshubhq/depshub/pkg/sources"
	"github.com/depshubhq/depshub/pkg/types"
)

type Linter struct {
	rules []types.Rule
}

func New() Linter {
	return Linter{
		rules: []types.Rule{
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

func (l Linter) Run(path string, configPath string) (mistakes []types.Mistake, err error) {
	config, err := config.New(configPath)

	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	scanner := manager.New(config)
	manifests, err := scanner.Scan(path)
	if err != nil {
		return nil, fmt.Errorf("failed to scan manifests: %w", err)
	}

	uniqueDependencies := scanner.UniqueDependencies(manifests)

	packagesData, err := sources.NewFetcher().Fetch(uniqueDependencies)

	if err != nil {
		fmt.Println("Error: ", err)
	}

	// Run all rules
	for _, rule := range l.rules {
		m, err := rule.Check(manifests, packagesData, config)

		if err != nil {
			return nil, fmt.Errorf("rule check failed: %w", err)
		}
		mistakes = append(mistakes, m...)
	}

	return mistakes, nil
}
