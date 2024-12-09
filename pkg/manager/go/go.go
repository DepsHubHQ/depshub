package gomanager

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/depshubhq/depshub/pkg/types"
	"golang.org/x/mod/modfile"
)

type Go struct{}

func (Go) GetType() types.ManagerType {
	return types.Go
}

func (Go) Managed(path string) bool {
	return filepath.Base(path) == "go.mod"
}

func (Go) Dependencies(path string) ([]types.Dependency, error) {
	var dependencies []types.Dependency

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	mod, err := modfile.Parse("go.mod", file, nil)

	if err != nil {
		return nil, err
	}

	for _, require := range mod.Require {
		if require.Indirect {
			continue
		}

		rawLine := []byte{}
		line := 0

		if require.Syntax != nil {
			rawLine = file[require.Syntax.Start.Byte:require.Syntax.End.Byte]
			line = require.Syntax.Start.Line
		}

		dependencies = append(dependencies, types.Dependency{
			Manager: types.Go,
			Name:    require.Mod.Path,
			Version: cleanVersion(require.Mod.Version),
			Dev:     false,
			Definition: types.Definition{
				Path:    path,
				RawLine: string(rawLine),
				Line:    line,
			},
		})
	}

	return dependencies, nil
}

// Returns the version without any prefix or suffix
func cleanVersion(version string) string {
	return strings.Trim(version, "^~*><= ")
}

func (Go) LockfilePath(path string) (string, error) {
	lockfilePath := filepath.Join(filepath.Dir(path), "go.sum")

	if _, err := os.Stat(lockfilePath); os.IsNotExist(err) {
		return "", fmt.Errorf("lockfile not found")
	}

	return lockfilePath, nil
}
