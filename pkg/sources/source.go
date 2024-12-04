package sources

import "go/types"

type Source interface {
	FetchPackageData(depName string) (types.Package, error)
}
