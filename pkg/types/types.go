package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

var ErrPackageNotFound = errors.New("package not found")
var ErrPackageUnpublished = errors.New("package unpublished")

type Manifest struct {
	Path         string
	Dependencies []Dependency
	*Lockfile
}

type Lockfile struct {
	Path string
}

type Dependency struct {
	Name    string
	Version string
	Dev     bool
	Definition
}

type Definition struct {
	Path    string
	RawLine string
	Line    int
}

type PackageVersion struct {
	Name       string `json:"name"`
	Version    string `json:"version"`
	Deprecated string `json:"deprecated"`
}

type Repository struct {
	Type      string `json:"type"`
	URL       string `json:"url"`
	Directory string `json:"directory"`
}

type Maintainer struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Download struct {
	Day       string
	Downloads int
}

type Package struct {
	Name          string
	Versions      map[string]PackageVersion
	LatestVersion PackageVersion
	Repository    Repository
	Time          map[string]time.Time `json:"time"`
	License       string
	Homepage      string
	Readme        string
	Maintainers   []Maintainer
	Keywords      []string
	Description   string
	Downloads     []Download
}

func (pv *PackageVersion) UnmarshalJSON(data []byte) error {
	// Create an auxiliary struct with Deprecated as json.RawMessage
	aux := struct {
		Name       string          `json:"name"`
		Version    string          `json:"version"`
		Deprecated json.RawMessage `json:"deprecated"`
	}{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Copy the standard fields
	pv.Name = aux.Name
	pv.Version = aux.Version

	// Handle the Deprecated field based on its type
	if len(aux.Deprecated) == 0 {
		pv.Deprecated = ""
		return nil
	}

	// Try to unmarshal as boolean first
	var boolValue bool
	if err := json.Unmarshal(aux.Deprecated, &boolValue); err == nil {
		if boolValue {
			pv.Deprecated = "deprecated"
		} else {
			pv.Deprecated = ""
		}
		return nil
	}

	// If not boolean, try to unmarshal as string
	var stringValue string
	if err := json.Unmarshal(aux.Deprecated, &stringValue); err == nil {
		pv.Deprecated = stringValue
		return nil
	}

	// If neither boolean nor string, return error
	return fmt.Errorf("deprecated field must be either boolean or string")
}

// Custom UnmarshalJSON for PackageData
func (pd *Package) UnmarshalJSON(data []byte) error {
	type Alias Package // Create an alias to avoid recursion
	aux := &struct {
		Repository json.RawMessage            `json:"repository"`
		Time       map[string]json.RawMessage `json:"time"`
		License    json.RawMessage            `json:"license"`
		*Alias
	}{
		Alias: (*Alias)(pd),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if err := checkUnpublished(aux.Time); err != nil {
		return err
	}

	if err := handleRepository(aux.Repository, pd); err != nil {
		return err
	}

	if err := handleLicense(aux.License, pd); err != nil {
		return err
	}

	return handleTime(aux.Time, pd)
}

func checkUnpublished(timeMap map[string]json.RawMessage) error {
	if _, ok := timeMap["unpublished"]; ok {
		return ErrPackageUnpublished
	}
	return nil
}

func handleRepository(repoData json.RawMessage, pd *Package) error {
	if len(repoData) == 0 {
		return nil // No repository data
	}

	var repos []Repository
	if err := json.Unmarshal(repoData, &repos); err == nil && len(repos) > 0 {
		pd.Repository = repos[0]
		return nil
	}

	var repoObj Repository
	if err := json.Unmarshal(repoData, &repoObj); err == nil {
		pd.Repository = repoObj
		return nil
	}

	var repoStr string
	if err := json.Unmarshal(repoData, &repoStr); err == nil {
		pd.Repository = Repository{
			Type:      "git",
			URL:       repoStr,
			Directory: "",
		}
		return nil
	}

	return fmt.Errorf("invalid repository format")
}

func handleLicense(licenseData json.RawMessage, pd *Package) error {
	if len(licenseData) == 0 {
		return nil // No license data
	}

	if err := json.Unmarshal(licenseData, &pd.License); err == nil {
		return nil
	}

	var license struct {
		Type string `json:"type"`
		URL  string `json:"url"`
	}
	if err := json.Unmarshal(licenseData, &license); err != nil {
		return fmt.Errorf("invalid license format: %w", err)
	}

	pd.License = license.Type
	return nil
}

func handleTime(timeMap map[string]json.RawMessage, pd *Package) error {
	pd.Time = make(map[string]time.Time)
	for key, rawValue := range timeMap {
		if key == "unpublished" {
			continue // Skip unpublished as we've already handled it
		}
		var t time.Time
		if err := json.Unmarshal(rawValue, &t); err != nil {
			return fmt.Errorf("invalid time format for key %s: %w", key, err)
		}
		pd.Time[key] = t
	}
	return nil
}
