package manager

import (
	"context"
	"fmt"

	"github.com/depshubhq/depshub/internal/linter/rules"
	"github.com/depshubhq/depshub/pkg/sources/npm"
	"github.com/depshubhq/depshub/pkg/types"
)

type fetcher struct{}

func NewFetcher() fetcher {
	return fetcher{}
}

func (f fetcher) Fetch(uniqueDependencies []string) (rules.PackagesInfo, error) {
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

	return packagesData, nil
}
