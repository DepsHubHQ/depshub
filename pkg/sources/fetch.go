package sources

import (
	"context"
	"fmt"

	"github.com/depshubhq/depshub/pkg/sources/go"
	"github.com/depshubhq/depshub/pkg/sources/npm"
	"github.com/depshubhq/depshub/pkg/types"
)

type fetcher struct{}

func NewFetcher() fetcher {
	return fetcher{}
}

const MaxConcurrent = 30

func (f fetcher) Fetch(uniqueDependencies []types.Dependency) (types.PackagesInfo, error) {
	// Create channels for results and errors
	type packageResult struct {
		pkg types.Package
		err error
	}
	resultChan := make(chan packageResult)

	// Launch goroutines for concurrent fetching
	npmSource := npm.NpmSource{}
	goSource := gosource.GoSource{}
	background := context.Background()
	activeRequests := 0

	// Use a semaphore to limit concurrent requests
	sem := make(chan struct{}, MaxConcurrent)

	for _, dep := range uniqueDependencies {
		activeRequests++

		go func(dep types.Dependency) {
			sem <- struct{}{} // Acquire semaphore
			defer func() {
				<-sem // Release semaphore
			}()

			var packageInfo types.Package
			var err error

			switch dep.Manager {
			case types.Npm:
				packageInfo, err = npmSource.FetchPackageData(background, dep.Name)
			case types.Go:
				packageInfo, err = goSource.FetchPackageData(dep.Name, dep.Version)
			}

			resultChan <- packageResult{
				pkg: packageInfo,
				err: err,
			}
		}(dep)
	}

	// Collect results
	var packagesData = make(types.PackagesInfo)

	for range activeRequests {
		result := <-resultChan
		if result.err != nil {
			fmt.Printf("Error fetching package data: %s\n", result.err)
			continue
		}
		packagesData[result.pkg.Name] = result.pkg
	}

	return packagesData, nil
}
