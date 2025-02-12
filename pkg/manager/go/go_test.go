package gomanager

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/depshubhq/depshub/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestGoGetType(t *testing.T) {
	manager := Go{}
	assert.Equal(t, types.Go, manager.GetType())
}

func TestGoManaged(t *testing.T) {
	manager := Go{}
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "go.mod file",
			path:     "path/to/go.mod",
			expected: true,
		},
		{
			name:     "not a go.mod file",
			path:     "path/to/package.json",
			expected: false,
		},
		{
			name:     "directory named go.mod",
			path:     "path/to/go.mod/file",
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

func TestGoDependencies(t *testing.T) {
	manager := Go{}
	testModPath := filepath.Join("testdata", "go.mod")
	fmt.Println(testModPath)

	deps, err := manager.Dependencies(testModPath)
	assert.NoError(t, err)
	assert.NotNil(t, deps)

	// Test for expected dependencies
	expectedDeps := []struct {
		manager types.ManagerType
		name    string
		line    int
		version string
	}{
		{
			manager: types.Go,
			name:    "github.com/charmbracelet/lipgloss",
			line:    6,
			version: "v1.0.0",
		},
		{
			manager: types.Go,
			name:    "github.com/sabhiram/go-gitignore",
			line:    7,
			version: "v0.0.0-20210923224102-525f6e181f06",
		},
		{
			manager: types.Go,
			name:    "github.com/spf13/cobra",
			line:    8,
			version: "v1.8.1",
		},
		{
			manager: types.Go,
			name:    "github.com/stretchr/testify",
			line:    9,
			version: "v1.6.1",
		},
	}

	assert.Equal(t, len(expectedDeps), len(deps))

	for i, exp := range expectedDeps {
		assert.Equal(t, exp.line, deps[i].Line)
		assert.Equal(t, exp.manager, deps[i].Manager)
		assert.Equal(t, exp.name, deps[i].Name)
		assert.Equal(t, exp.version, deps[i].Version)
		assert.False(t, deps[i].Dev)
		assert.Equal(t, testModPath, deps[i].Path)
		assert.NotEmpty(t, deps[i].RawLine)
		assert.Greater(t, deps[i].Line, 0)
	}
}

func TestCleanVersion(t *testing.T) {
	tests := []struct {
		name     string
		version  string
		expected string
	}{
		{
			name:     "simple version",
			version:  "1.0.0",
			expected: "1.0.0",
		},
		{
			name:     "version with v prefix",
			version:  "v1.0.0",
			expected: "v1.0.0",
		},
		{
			name:     "version with comparison operator",
			version:  ">=1.0.0",
			expected: "1.0.0",
		},
		{
			name:     "version with multiple prefixes",
			version:  ">=1.0.0",
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
