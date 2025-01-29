package maven

import (
	"time"

	"github.com/depshubhq/depshub/pkg/types"
	"github.com/edoardottt/depsdev/pkg/depsdev"
)

type MavenSource struct{}

func (MavenSource) FetchPackageData(name string, version string) (types.Package, error) {
	var target types.Package

	info, err := depsdev.NewAPI().GetInfo("maven", name)

	if err != nil {
		return target, err
	}

	v, err := depsdev.NewAPI().GetVersion("maven", name, version)

	if err != nil {
		return target, err
	}

	target.Name = name
	target.License = v.Licenses[0]
	target.Versions = make(map[string]types.PackageVersion)
	target.Time = make(map[string]time.Time)

	for _, v := range info.Versions {
		if target.Versions == nil {
			target.Versions = make(map[string]types.PackageVersion)
		}

		target.Versions[v.VersionKey.Version] = types.PackageVersion{
			Name:    v.VersionKey.Name,
			Version: v.VersionKey.Version,
		}

		if !v.PublishedAt.IsZero() {
			target.Time[v.VersionKey.Version] = v.PublishedAt
		}
	}

	return target, nil
}
