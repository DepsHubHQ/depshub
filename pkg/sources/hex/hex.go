package hex

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/depshubhq/depshub/pkg/types"
	"net/http"
	"time"
)

type HexSource struct{}

type Metadata struct {
	Licenses []string `json:"licenses"`
}

type Downloads struct {
	Week int `json:"week"`
}

type Release struct {
	Version    string    `json:"version"`
	InsertedAt time.Time `json:"inserted_at"`
}

// https://hexpm.docs.apiary.io/#reference/packages/package
type HexPackage struct {
	Name       string    `json:"name"`
	Repository string    `json:"repository"`
	Private    bool      `json:"private"`
	Metadata   Metadata  `json:"meta"`
	Downloads  Downloads `json:"downloads"`
	Releases   []Release `json:"releases"`
}

func (s HexSource) FetchPackageData(ctx context.Context, name string) (types.Package, error) {
	var target HexPackage
	var result types.Package

	if err := s.fetchPackageInfo(ctx, name, &target); err != nil {
		return types.Package{}, err
	}

	// Convert HexPackage to the generic types.Package
	result.Name = target.Name
	result.Versions = make(map[string]types.PackageVersion)
	result.Time = make(map[string]time.Time)

	if len(target.Metadata.Licenses) > 0 {
		result.License = target.Metadata.Licenses[0]
	}

	result.Downloads = []types.Download{
		{Day: time.Now().Format("2006-01-02"), Downloads: target.Downloads.Week},
	}

	for _, version := range target.Releases {
		result.Versions[version.Version] = types.PackageVersion{
			Name:    target.Name,
			Version: version.Version,
		}

		result.Time[version.Version] = version.InsertedAt
	}

	return result, nil
}

func (HexSource) fetchPackageInfo(ctx context.Context, name string, target *HexPackage) error {
	url := fmt.Sprintf("https://hex.pm/api/packages/%s", name)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("error creating request for %s information from Hex registry: %w", name, err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error getting %s information from Hex registry: %w", name, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 || resp.StatusCode == 405 {
		return types.ErrPackageNotFound
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("error getting %s information from Hex registry: %s", name, resp.Status)
	}

	return json.NewDecoder(resp.Body).Decode(target)
}
