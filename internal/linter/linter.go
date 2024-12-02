package linter

import (
	"context"
	"fmt"

	"github.com/depshubhq/depshub/internal/linter/rules"
	"github.com/depshubhq/depshub/pkg/manager"
	"github.com/depshubhq/depshub/pkg/sources/npm"
	"github.com/depshubhq/depshub/pkg/types"
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
		},
	}
}

func (l Linter) Run(path string) (mistakes []rules.Mistake, err error) {
	scanner := manager.New()
	manifests, err := scanner.Scan(path)
	if err != nil {
		return nil, fmt.Errorf("failed to scan manifests: %w", err)
	}

	uniqueDependencies := scanner.Unique(manifests)

	// Create channels for results and errors
	type packageResult struct {
		pkg types.Package
		err error
	}
	resultChan := make(chan packageResult)

	// Launch goroutines for concurrent fetching
	npmManager := npm.NpmManager{}
	background := context.Background()
	activeRequests := 0

	// Use a semaphore to limit concurrent requests
	const maxConcurrent = 30
	sem := make(chan struct{}, maxConcurrent)

	for _, name := range uniqueDependencies {
		activeRequests++

		go func(depName string) {
			sem <- struct{}{} // Acquire semaphore
			defer func() {
				<-sem // Release semaphore
			}()

			npmPackage, err := npmManager.FetchPackageData(background, depName)
			resultChan <- packageResult{
				pkg: npmPackage,
				err: err,
			}
		}(name)
	}

	// Collect results
	var packagesData = make(rules.PackagesInfo)
	for i := 0; i < activeRequests; i++ {
		result := <-resultChan
		if result.err != nil {
			fmt.Printf("Error fetching package data: %s\n", result.err)
			continue
		}
		packagesData[result.pkg.Name] = result.pkg
	}

	fmt.Printf("Successfully fetched %d packages\n", len(packagesData))

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
