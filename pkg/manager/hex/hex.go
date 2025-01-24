package hex

import (
	"fmt"
	"github.com/depshubhq/depshub/pkg/types"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Hex struct{}

func (Hex) GetType() types.ManagerType {
	return types.Hex
}

func (Hex) Managed(path string) bool {
	return filepath.Base(path) == "mix.exs"
}

func (Hex) Dependencies(path string) ([]types.Dependency, error) {
	var dependencies []types.Dependency

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	content := string(file)
	lines := strings.Split(content, "\n")

	// Regex pattern to match dependencies in mix.exs
	// Matches lines like: `{:ecto, "~> 3.7"}`
	dependencyPattern := regexp.MustCompile(`^\s*\{\s*:(\w+),\s*"([^"]+)"`)
	inDepsBlock := false

	for lineNum, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// Detect start of deps block
		if strings.HasPrefix(trimmedLine, "defp deps do") {
			inDepsBlock = true
			continue
		}

		// Detect end of deps block
		if inDepsBlock && trimmedLine == "end" {
			break
		}

		if inDepsBlock {
			// Skip lines with git dependencies
			if strings.Contains(trimmedLine, "git:") {
				continue
			}

			matches := dependencyPattern.FindStringSubmatch(trimmedLine)
			if matches == nil {
				continue
			}

			name := matches[1]
			version := cleanVersion(matches[2])

			dependencies = append(dependencies, types.Dependency{
				Manager: types.Hex,
				Name:    name,
				Version: version,
				Dev:     false,
				Definition: types.Definition{
					Path:    path,
					RawLine: strings.TrimSpace(line),
					Line:    lineNum + 1, // Line number starts from 1
				},
			})
		}
	}

	return dependencies, nil
}

// Returns the version without any prefix or suffix
func cleanVersion(version string) string {
	// Trim spaces
	version = strings.TrimSpace(version)

	// Remove constraints like >=, ~>, etc.
	version = strings.TrimPrefix(version, ">=")
	version = strings.TrimPrefix(version, "~>")

	return strings.TrimSpace(version)
}

func (Hex) LockfilePath(path string) (string, error) {
	// Check for requirements.txt in the same directory
	lockfilePath := filepath.Join(filepath.Dir(path), "mix.lock")

	if _, err := os.Stat(lockfilePath); err == nil {
		return lockfilePath, nil
	}

	return "", fmt.Errorf("no lockfile found")
}
