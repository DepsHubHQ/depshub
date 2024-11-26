package npm

import (
	"path/filepath"

	"github.com/depshubhq/depshub/pkg/types"
)

type Npm struct{}

func (Npm) Managed(path string) bool {
	fileName := filepath.Base(path)
	return fileName == "package.json"
}

func (Npm) Dependencies(path string) ([]types.Dependency, error) {
	return []types.Dependency{
		{Name: "test"},
	}, nil
}
