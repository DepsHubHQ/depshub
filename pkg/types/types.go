package types

import "time"

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
	Name    string `json:"name"`
	Version string `json:"version"`
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
