package pypi

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/depshubhq/depshub/pkg/types"
	"net/http"
	"time"
)

type PyPISource struct{}

type Release struct {
	Filename         string `json:"filename"`
	PackageType      string `json:"packagetype"`
	Size             int    `json:"size"`
	UploadTime       string `json:"upload_time"`
	UploadTimeISOStr string `json:"upload_time_iso_8601"`
	URL              string `json:"url"`
	Yanked           bool   `json:"yanked"`
	YankedReason     string `json:"yanked_reason"`
}

type Info struct {
	Author       string   `json:"author"`
	Description  string   `json:"description"`
	License      string   `json:"license"`
	Name         string   `json:"name"`
	Summary      string   `json:"summary"`
	Version      string   `json:"version"`
	RequiresDist []string `json:"requires_dist"`
}

type PyPIPackage struct {
	Info     Info                 `json:"info"`
	Releases map[string][]Release `json:"releases"`
}

func (s PyPISource) FetchPackageData(ctx context.Context, name string) (types.Package, error) {
	var target PyPIPackage
	var result types.Package

	if err := s.fetchPackageInfo(ctx, name, &target); err != nil {
		return types.Package{}, err
	}

	// Convert PyPIPackage to the generic types.Package
	result.Name = target.Info.Name
	result.License = target.Info.License
	result.Versions = make(map[string]types.PackageVersion)
	result.Time = make(map[string]time.Time)

	// Process each version
	for version, releases := range target.Releases {
		if len(releases) > 0 {
			deprecated := ""
			if releases[0].Yanked {
				deprecated = "yanked"
				if releases[0].YankedReason != "" {
					deprecated = fmt.Sprintf("yanked: %s", releases[0].YankedReason)
				}
			}

			pv := types.PackageVersion{
				Name:       target.Info.Name,
				Version:    version,
				Deprecated: deprecated,
			}
			result.Versions[version] = pv

			// Parse upload time
			if uploadTime, err := time.Parse(time.RFC3339, releases[0].UploadTimeISOStr); err == nil {
				result.Time[version] = uploadTime
			}
		}
	}

	// PyPI doesn't provide direct download counts in the API response
	// You might want to fetch this separately if needed
	result.Downloads = []types.Download{}

	return result, nil
}

func (PyPISource) fetchPackageInfo(ctx context.Context, name string, target *PyPIPackage) error {
	url := fmt.Sprintf("https://pypi.org/pypi/%s/json", name)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("error creating request for %s information from PyPI registry: %w", name, err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error getting %s information from PyPI registry: %w", name, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 || resp.StatusCode == 405 {
		return types.ErrPackageNotFound
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("error getting %s information from PyPI registry: %s", name, resp.Status)
	}

	return json.NewDecoder(resp.Body).Decode(target)
}
