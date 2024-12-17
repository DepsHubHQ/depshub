package crates

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/depshubhq/depshub/pkg/types"
)

type CratesSource struct{}

type Version struct {
	ID          int       `json:"id"`
	Num         string    `json:"num"`
	License     string    `json:"license"`
	Yanked      bool      `json:"yanked"`
	YankMessage string    `json:"yank_message"`
	CreatedAt   time.Time `json:"created_at"`
}

type Crate struct {
	RecentDownloads int    `json:"recent_downloads"`
	DefaultVersion  string `json:"default_version"`
	Versions        []int  `json:"versions"`
}

type CratePackage struct {
	Name     string
	Crate    Crate
	Versions []Version `json:"versions"`
}

func (s CratesSource) FetchPackageData(ctx context.Context, name string) (types.Package, error) {
	var target CratePackage
	var result types.Package

	if err := s.fetchPackageInfo(ctx, name, &target); err != nil {
		return types.Package{}, err
	}

	// Convert CratePackage to the generic types.Package
	result.Name = target.Name
	var currentVersion Version

	for _, version := range target.Versions {
		if version.ID == target.Crate.Versions[0] {
			currentVersion = version
		}
	}
	result.License = currentVersion.License
	result.Versions = make(map[string]types.PackageVersion)
	result.Time = make(map[string]time.Time)

	for _, version := range target.Versions {
		deprecated := ""

		if version.Yanked {
			deprecated = "yanked"
		}

		if len(version.YankMessage) > 0 {
			deprecated = fmt.Sprintf("yanked: %s", version.YankMessage)
		}

		pv := types.PackageVersion{
			Name:       target.Name,
			Version:    version.Num,
			Deprecated: deprecated,
		}

		result.Versions[pv.Version] = pv
		result.Time[pv.Version] = version.CreatedAt
	}

	result.Downloads = []types.Download{
		{Day: time.Now().Format("2006-01-02"), Downloads: target.Crate.RecentDownloads},
	}

	return result, nil
}

func (CratesSource) fetchPackageInfo(ctx context.Context, name string, target *CratePackage) error {
	url := fmt.Sprintf("https://crates.io/api/v1/crates/%s", name)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("error creating request for %s information from npm registry: %w", name, err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error getting %s information from npm registry: %w", name, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 || resp.StatusCode == 405 {
		return types.ErrPackageNotFound
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("error getting %s information from npm registry: %s", name, resp.Status)
	}

	return json.NewDecoder(resp.Body).Decode(target)
}
