package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type ManagerType int

const (
	Npm ManagerType = iota
	Go
	Cargo
	Pip
)

var ErrPackageNotFound = errors.New("package not found")
var ErrPackageUnpublished = errors.New("package unpublished")

type Manifest struct {
	Manager      ManagerType
	Path         string
	Dependencies []Dependency
	*Lockfile
}

type Lockfile struct {
	Path string
}

type Dependency struct {
	Manager ManagerType
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

type Download struct {
	Day       string
	Downloads int
}

// A map of package names to package information.
type PackagesInfo = map[string]Package

type Package struct {
	Name      string
	Versions  map[string]PackageVersion
	Time      map[string]time.Time `json:"time"`
	License   string
	Downloads []Download
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
		Time    map[string]json.RawMessage `json:"time"`
		License json.RawMessage            `json:"license"`
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
