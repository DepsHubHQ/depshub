package maven

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/depshubhq/depshub/pkg/types"
	"github.com/vifraa/gopom"
)

type Maven struct{}

func (Maven) GetType() types.ManagerType {
	return types.Maven
}

func (Maven) Managed(path string) bool {
	path = strings.ToLower(path)
	return filepath.Base(path) == "pom.xml"
}

func (Maven) Dependencies(path string) ([]types.Dependency, error) {
	var dependencies []types.Dependency

	fileBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	parsedPom, err := gopom.Parse(path)
	if err != nil {
		return nil, err
	}

	if parsedPom.Dependencies == nil {
		return []types.Dependency{}, nil
	}

	for _, dep := range *parsedPom.Dependencies {
		if dep.GroupID == nil || dep.ArtifactID == nil {
			continue
		}
		name := fmt.Sprintf("%s:%s", *dep.GroupID, *dep.ArtifactID)
		line, rawLine := findLineInfo(fileBytes, *dep.ArtifactID)
		version := ""
		if dep.Version != nil {
			version = *dep.Version
		} else {
			// Try to find version in dependencyManagement
			if parsedPom.DependencyManagement != nil && parsedPom.DependencyManagement.Dependencies != nil {
				for _, depMan := range *parsedPom.DependencyManagement.Dependencies {
					if depMan.GroupID != nil && depMan.ArtifactID != nil && *depMan.ArtifactID == *dep.ArtifactID {
						version = *depMan.Version
					}
				}
			}
		}

		dependencies = append(dependencies, types.Dependency{
			Manager: types.Maven,
			Name:    name,
			//  TODO We should use the version from the lockfile instead
			Version: cleanVersion(version),
			// FIXME We should use the scope to determine if it's a dev dependency
			Dev: false,
			Definition: types.Definition{
				Path:    path,
				RawLine: rawLine,
				Line:    line,
			},
		})
	}

	if parsedPom.DependencyManagement == nil || parsedPom.DependencyManagement.Dependencies == nil {
		return dependencies, nil
	}

	for _, dep := range *parsedPom.DependencyManagement.Dependencies {
		if dep.GroupID == nil || dep.ArtifactID == nil {
			continue
		}
		name := fmt.Sprintf("%s:%s", *dep.GroupID, *dep.ArtifactID)
		line, rawLine := findLineInfo(fileBytes, *dep.ArtifactID)
		version := ""
		if dep.Version != nil {
			version = *dep.Version
		}

		dependencies = append(dependencies, types.Dependency{
			Manager: types.Maven,
			Name:    name,
			//  TODO We should use the version from the lockfile instead
			Version: cleanVersion(version),
			// FIXME We should use the scope to determine if it's a dev dependency
			Dev: false,
			Definition: types.Definition{
				Path:    path,
				RawLine: rawLine,
				Line:    line,
			},
		})
	}

	return dependencies, nil
}

// Maven doesn't have a lock file but it's not an error
func (Maven) LockfilePath(path string) (string, error) {
	return "", nil
}

// Returns the version without any prefix or suffix
func cleanVersion(version string) string {
	return strings.Trim(version, "v^~*><= ")
}

func findLineInfo(data []byte, key string) (line int, rawLine string) {
	lines := bytes.Split(data, []byte{'\n'})

	for i, line := range lines {
		trimmed := bytes.TrimSpace(line)

		// Look for our key while in the correct section
		if bytes.Contains(trimmed, []byte(fmt.Sprintf("<artifactId>%s</artifactId>", key))) {
			return i + 1, string(trimmed)
		}
	}

	return 0, ""
}
