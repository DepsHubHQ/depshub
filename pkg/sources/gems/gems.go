package gems

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/depshubhq/depshub/pkg/types"
)

// Types for the RubyGems API response
type GemVersionInfo struct {
	Authors         string            `json:"authors"`
	BuiltAt         time.Time         `json:"built_at"`
	CreatedAt       time.Time         `json:"created_at"`
	DownloadsCount  int               `json:"downloads_count"`
	Metadata        map[string]string `json:"metadata"`
	Number          string            `json:"number"`
	Summary         string            `json:"summary"`
	Platform        string            `json:"platform"`
	RubyGemsVersion string            `json:"rubygems_version"`
	RubyVersion     *string           `json:"ruby_version"`
	Prerelease      bool              `json:"prerelease"`
	Licenses        []string          `json:"licenses"`
}

type GemsSource struct{}

func (GemsSource) FetchPackageData(name string, version string) (types.Package, error) {
	var target types.Package
	var gemVersions []GemVersionInfo

	url := fmt.Sprintf("https://rubygems.org/api/v1/versions/%s.json", name)
	resp, err := http.Get(url)
	if err != nil {
		return target, fmt.Errorf("failed to fetch data from RubyGems API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return target, fmt.Errorf("non-OK response from RubyGems API: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&gemVersions); err != nil {
		return target, fmt.Errorf("failed to parse RubyGems response: %w", err)
	}

	target.Name = name
	target.Versions = make(map[string]types.PackageVersion)
	target.Time = make(map[string]time.Time)

	for _, gemVersion := range gemVersions {
		target.Versions[gemVersion.Number] = types.PackageVersion{
			Name:    name,
			Version: gemVersion.Number,
		}

		if !gemVersion.CreatedAt.IsZero() {
			target.Time[gemVersion.Number] = gemVersion.CreatedAt
		}

		// Use the license from the matching version if available
		if gemVersion.Number == version && len(gemVersion.Licenses) > 0 {
		}

		if gemVersion.Number == version {
			if len(gemVersion.Licenses) > 0 {
				target.License = gemVersion.Licenses[0]
			}

			target.Downloads = []types.Download{
				{Day: time.Now().Format("2006-01-02"), Downloads: gemVersion.DownloadsCount},
			}
		}
	}

	return target, nil
}
