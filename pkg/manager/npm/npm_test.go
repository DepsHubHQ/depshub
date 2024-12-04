package npm

import (
	"path/filepath"
	"testing"

	"github.com/depshubhq/depshub/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestNpmGetType(t *testing.T) {
	manager := Npm{}
	assert.Equal(t, types.Npm, manager.GetType())
}

func TestNpmManaged(t *testing.T) {
	manager := Npm{}
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "package.json file",
			path:     "path/to/package.json",
			expected: true,
		},
		{
			name:     "not a package.json file",
			path:     "path/to/go.mod",
			expected: false,
		},
		{
			name:     "directory named package.json",
			path:     "path/to/package.json/file",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := manager.Managed(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNpmDependencies(t *testing.T) {
	manager := Npm{}
	testPkgPath := filepath.Join("testdata", "package.json")

	deps, err := manager.Dependencies(testPkgPath)
	assert.NoError(t, err)
	assert.NotNil(t, deps)

	// Test dependencies order by line number
	for i := 1; i < len(deps); i++ {
		assert.GreaterOrEqual(t, deps[i].Line, deps[i-1].Line,
			"Dependencies should be sorted by line number")
	}

	// Helper function to find dependency by name
	findDep := func(deps []types.Dependency, name string) *types.Dependency {
		for _, dep := range deps {
			if dep.Name == name {
				return &dep
			}
		}
		return nil
	}

	// Test for expected dependencies
	testCases := []struct {
		name        string
		version     string
		isDev       bool
		shouldExist bool
	}{
		{
			name:        "astro",
			version:     "4.16.10",
			isDev:       false,
			shouldExist: true,
		},
		{
			name:        "typescript",
			version:     "5.7.2",
			isDev:       true,
			shouldExist: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dep := findDep(deps, tc.name)
			if tc.shouldExist {
				assert.NotNil(t, dep, "Dependency %s should exist", tc.name)
				if dep != nil {
					assert.Equal(t, tc.version, dep.Version)
					assert.Equal(t, tc.isDev, dep.Dev)
					assert.Equal(t, testPkgPath, dep.Path)
					assert.NotEmpty(t, dep.RawLine)
					assert.Greater(t, dep.Line, 0)
				}
			} else {
				assert.Nil(t, dep, "Dependency %s should not exist", tc.name)
			}
		})
	}
}

func TestFindLineInfo(t *testing.T) {
	testJSON := []byte(`{
  "name": "test-package",
  "dependencies": {
    "react": "^18.2.0",
    "lodash": "4.17.21"
  },
  "devDependencies": {
    "jest": "^29.0.0",
    "test": "^29.0.0"
  }
}`)

	tests := []struct {
		name         string
		section      string
		key          string
		expectedLine int
		expectedRaw  string
		shouldFind   bool
	}{
		{
			name:         "regular dependency",
			section:      "dependencies",
			key:          "react",
			expectedLine: 3,
			expectedRaw:  `"react": "^18.2.0",`,
			shouldFind:   true,
		},
		{
			name:         "dev dependency",
			section:      "devDependencies",
			key:          "jest",
			expectedLine: 7,
			expectedRaw:  `"jest": "^29.0.0",`,
			shouldFind:   true,
		},
		{
			name:         "dev dependency 2",
			section:      "devDependencies",
			key:          "test",
			expectedLine: 8,
			expectedRaw:  `"test": "^29.0.0"`,
			shouldFind:   true,
		},
		{
			name:         "non-existent dependency",
			section:      "dependencies",
			key:          "nonexistent",
			expectedLine: 0,
			expectedRaw:  "",
			shouldFind:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			line, rawLine := findLineInfo(testJSON, tt.section, tt.key)
			if tt.shouldFind {
				assert.Equal(t, tt.expectedLine, line)
				assert.Equal(t, tt.expectedRaw, rawLine)
			} else {
				assert.Equal(t, 0, line)
				assert.Empty(t, rawLine)
			}
		})
	}
}

func TestCleanVersion(t *testing.T) {
	tests := []struct {
		name     string
		version  string
		expected string
	}{
		{
			name:     "caret version",
			version:  "^1.0.0",
			expected: "1.0.0",
		},
		{
			name:     "tilde version",
			version:  "~1.0.0",
			expected: "1.0.0",
		},
		{
			name:     "star version",
			version:  "*",
			expected: "",
		},
		{
			name:     "greater than version",
			version:  ">1.0.0",
			expected: "1.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cleanVersion(tt.version)
			assert.Equal(t, tt.expected, result)
		})
	}
}
