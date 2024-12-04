package gomanager

import "github.com/depshubhq/depshub/pkg/types"

type Go struct{}

func (Go) GetType() types.ManagerType {
	return types.Go
}

func (Go) Managed(path string) bool {
	return false
}

func (Go) Dependencies(path string) ([]types.Dependency, error) {
	return nil, nil
}

func (Go) LockfilePath(path string) (string, error) {
	return "", nil
}
