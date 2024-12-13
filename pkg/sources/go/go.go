package gosource

import (
	"time"

	"github.com/depshubhq/depshub/pkg/types"
	"github.com/edoardottt/depsdev/pkg/depsdev"
)

type GoSource struct{}

func (GoSource) FetchPackageData(name string, version string) (types.Package, error) {
	var target types.Package

	info, err := depsdev.NewAPI().GetInfo("go", name)

	if err != nil {
		return target, err
	}

	v, err := depsdev.NewAPI().GetVersion("go", name, version)

	if err != nil {
		return target, err
	}

	target.Name = name
	target.License = v.Licenses[0]
	target.Time = map[string]time.Time{version: v.PublishedAt}

	for _, v := range info.Versions {
		if target.Versions == nil {
			target.Versions = make(map[string]types.PackageVersion)
		}

		target.Versions[v.VersionKey.Version] = types.PackageVersion{
			Name:    v.VersionKey.Name,
			Version: v.VersionKey.Version,
		}
	}

	return target, nil
}
