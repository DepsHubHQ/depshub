package manager

import "github.com/depshubhq/depshub/pkg/types"

type Manager interface {
	Managed(path string) bool
	Dependencies(path string) ([]types.Dependency, error)
}
