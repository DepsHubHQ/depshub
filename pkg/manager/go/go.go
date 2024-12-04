package gomanager

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/depshubhq/depshub/pkg/types"
)

type Go struct{}

func (Go) GetType() types.ManagerType {
	return types.Go
}

func (Go) Managed(path string) bool {
	return filepath.Base(path) == "go.mod"
}

func (Go) Dependencies(path string) ([]types.Dependency, error) {
	return nil, nil
}

func (Go) LockfilePath(path string) (string, error) {
	lockfilePath := filepath.Join(filepath.Dir(path), "go.sum")

	if _, err := os.Stat(lockfilePath); os.IsNotExist(err) {
		return "", fmt.Errorf("lockfile not found")
	}

	return lockfilePath, nil
}
