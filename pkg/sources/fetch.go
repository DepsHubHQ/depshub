package sources

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/depshubhq/depshub/pkg/sources/crates"
	"github.com/depshubhq/depshub/pkg/sources/go"
	"github.com/depshubhq/depshub/pkg/sources/hex"
	"github.com/depshubhq/depshub/pkg/sources/maven"
	"github.com/depshubhq/depshub/pkg/sources/npm"
	"github.com/depshubhq/depshub/pkg/sources/pypi"
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
	cratesSource := crates.CratesSource{}
	pypiSource := pypi.PyPISource{}
	hexSource := hex.HexSource{}
	mavenSource := maven.MavenSource{}

	background := context.Background()

	// Use a semaphore to limit concurrent requests
	sem := make(chan struct{}, MaxConcurrent)
	c, err := NewFileCache("dependencies")

	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	for _, dep := range uniqueDependencies {
		wg.Add(1)

		go func() {
			defer wg.Done()

			sem <- struct{}{} // Acquire semaphore
			defer func() {
				<-sem // Release semaphore
			}()

			var packageInfo types.Package
			var err error

			key := fmt.Sprintf("%d-%s", dep.Manager, dep.Name)

			exists, err := c.Get(key, &packageInfo)

			if err != nil {
				fmt.Printf("Error getting cache: %s\n", err)
			}

			if !exists {
				switch dep.Manager {
				case types.Npm:
					packageInfo, err = npmSource.FetchPackageData(background, dep.Name)
				case types.Go:
					packageInfo, err = goSource.FetchPackageData(dep.Name, dep.Version)
				case types.Cargo:
					packageInfo, err = cratesSource.FetchPackageData(background, dep.Name)
				case types.Pip, types.Pyproject:
					packageInfo, err = pypiSource.FetchPackageData(background, dep.Name)
				case types.Hex:
					packageInfo, err = hexSource.FetchPackageData(background, dep.Name)
				case types.Maven:
					packageInfo, err = mavenSource.FetchPackageData(dep.Name, dep.Version)
				}

				if err != nil {
					fmt.Printf("Error fetching package data: %s\n", err)
				} else {
					c.Set(key, packageInfo, 48*time.Hour)
				}
			}

			resultChan <- packageResult{
				pkg: packageInfo,
				err: err,
			}
		}()
	}

	// Start a goroutine to close resultChan after all workers are done
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results
	var packagesData = make(types.PackagesInfo)

	for result := range resultChan {
		if result.err != nil {
			fmt.Printf("Error fetching package data: %s\n", result.err)
			continue
		}
		packagesData[result.pkg.Name] = result.pkg
	}

	return packagesData, nil
}
