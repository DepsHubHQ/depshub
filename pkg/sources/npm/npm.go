package npm

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/depshubhq/depshub/pkg/types"
)

type NpmManager struct{}

func (npm NpmManager) FetchPackageData(ctx context.Context, name string) (types.Package, error) {
	var target types.Package

	if err := npm.fetchPackageInfo(ctx, name, &target); err != nil {
		return types.Package{}, err
	}

	if err := npm.fetchLatestVersion(ctx, name, &target); err != nil {
		return types.Package{}, err
	}

	if err := npm.fetchDownloads(ctx, name, &target); err != nil {
		return types.Package{}, err
	}

	return target, nil
}

func (NpmManager) fetchPackageInfo(ctx context.Context, name string, target *types.Package) error {
	url := fmt.Sprintf("https://registry.npmjs.org/%s", name)
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

func (NpmManager) fetchLatestVersion(ctx context.Context, name string, target *types.Package) error {
	url := fmt.Sprintf("https://registry.npmjs.org/%s/latest", name)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("error creating request for %s information from npm registry: %w", name, err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error getting latest %s information from npm registry: %w", name, err)
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(&target.LatestVersion)
}

func (NpmManager) fetchDownloads(ctx context.Context, name string, target *types.Package) error {
	from := time.Now().AddDate(0, -11, 0).Format("2006-01-02")
	to := time.Now().Format("2006-01-02")
	url := fmt.Sprintf("https://api.npmjs.org/downloads/range/%s:%s/%s", from, to, name)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("error creating request for %s downloads information from npm registry: %w", name, err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error getting latest %s downloads information from npm registry: %w", name, err)
	}
	defer resp.Body.Close()

	var downloadsDataTarget types.Package
	if err := json.NewDecoder(resp.Body).Decode(&downloadsDataTarget); err != nil {
		log.Printf("error parsing latest %s downloads data into json: %s", name, err)
	}

	target.Downloads = append(target.Downloads, downloadsDataTarget.Downloads...)
	return nil
}
