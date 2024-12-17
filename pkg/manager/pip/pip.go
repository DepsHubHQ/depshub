package pipmanager

import (
	"bufio"
	"fmt"
	"github.com/depshubhq/depshub/pkg/types"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Pip struct{}

func (Pip) GetType() types.ManagerType {
	return types.Pip
}

func (Pip) Managed(path string) bool {
	return filepath.Base(path) == "requirements.txt"
}

func (Pip) Dependencies(path string) ([]types.Dependency, error) {
	var dependencies []types.Dependency

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0

	// Regex patterns for parsing requirements.txt entries
	commentPattern := regexp.MustCompile(`^\s*#`)
	versionPattern := regexp.MustCompile(`^([^=<>!~]+)(==|>=|<=|!=|~=|>|<)(.+)`)

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Skip empty lines and comments
		if line == "" || commentPattern.MatchString(line) {
			continue
		}

		// Handle -r requirements inclusion
		if strings.HasPrefix(line, "-r ") {
			continue // Skip requirements file inclusions for now
		}

		// Extract package name and version
		var name, version string
		if matches := versionPattern.FindStringSubmatch(line); matches != nil {
			name = strings.TrimSpace(matches[1])
			version = cleanVersion(matches[3])
		} else {
			// Package with no version specified
			name = strings.TrimSpace(line)
			version = ""
		}

		dependencies = append(dependencies, types.Dependency{
			Manager: types.Pip,
			Name:    name,
			Version: version,
			Dev:     false,
			Definition: types.Definition{
				Path:    path,
				RawLine: line,
				Line:    lineNum,
			},
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return dependencies, nil
}

// Returns the version without any prefix or suffix
func cleanVersion(version string) string {
	version = strings.TrimSpace(version)

	// Remove any comments that might be at the end
	if idx := strings.Index(version, "#"); idx != -1 {
		version = version[:idx]
	}

	version = strings.Trim(version, "^~*><= ")

	return version
}

func (Pip) LockfilePath(path string) (string, error) {
	// Check for requirements.txt in the same directory
	lockfilePath := filepath.Join(filepath.Dir(path), "requirements.lock")
	if _, err := os.Stat(lockfilePath); os.IsNotExist(err) {
		// Some projects use pip-lock instead
		lockfilePath = filepath.Join(filepath.Dir(path), "pip.lock")
		if _, err := os.Stat(lockfilePath); os.IsNotExist(err) {
			return "", fmt.Errorf("lockfile not found")
		}
	}
	return lockfilePath, nil
}
