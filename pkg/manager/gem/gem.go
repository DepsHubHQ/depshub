package gem

import (
	"bufio"
	"fmt"
	"github.com/depshubhq/depshub/pkg/types"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Gem struct{}

func (Gem) GetType() types.ManagerType {
	return types.Gem
}

func (Gem) Managed(path string) bool {
	return filepath.Base(path) == "Gemfile"
}

func (Gem) Dependencies(path string) ([]types.Dependency, error) {
	var dependencies []types.Dependency

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0

	// Regex patterns for parsing Gemfile entries
	commentPattern := regexp.MustCompile(`^\s*#`)
	gemPattern := regexp.MustCompile(`^\s*gem\s+'([^']+)'(?:,\s*'([^']+)')?(?:,\s*'([^']+)')?`)
	groupPattern := regexp.MustCompile(`^\s*group\s+(:[a-zA-Z0-9_]+(?:,\s*:[a-zA-Z0-9_]+)*)\s+do`)

	currentGroups := []string{}

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || commentPattern.MatchString(line) {
			continue
		}

		// Handle group entries
		if matches := groupPattern.FindStringSubmatch(line); matches != nil {
			currentGroups = strings.Split(strings.ReplaceAll(matches[1], ":", ""), ",")
			continue
		} else if strings.HasPrefix(line, "end") {
			currentGroups = []string{}
			continue
		}

		// Extract gem name and version
		if matches := gemPattern.FindStringSubmatch(line); matches != nil {
			name := matches[1]
			version := ""
			if matches[2] != "" {
				version = cleanVersion(matches[2])
			}

			dependencies = append(dependencies, types.Dependency{
				Manager: types.Gem,
				Name:    name,
				Version: version,
				Dev:     contains(currentGroups, "development") || contains(currentGroups, "test"),
				Definition: types.Definition{
					Path:    path,
					RawLine: line,
					Line:    lineNum,
				},
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return dependencies, nil
}

func cleanVersion(version string) string {
	version = strings.TrimSpace(version)

	// Remove any comments that might be at the end
	if idx := strings.Index(version, "#"); idx != -1 {
		version = version[:idx]
	}

	version = strings.Trim(version, "^~*><= ")

	return version
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func (Gem) LockfilePath(path string) (string, error) {
	lockfilePath := filepath.Join(filepath.Dir(path), "Gemfile.lock")

	if _, err := os.Stat(lockfilePath); os.IsNotExist(err) {
		return "", fmt.Errorf("lockfile not found")
	}

	return lockfilePath, nil
}
