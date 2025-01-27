package hex

import (
	"github.com/depshubhq/depshub/pkg/types"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestHex_GetType(t *testing.T) {
	manager := Hex{}
	assert.Equal(t, types.Hex, manager.GetType())
}

func TestHex_Managed(t *testing.T) {
	manager := Hex{}
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "mix.exs file",
			path:     "path/to/mix.exs",
			expected: true,
		},
		{
			name:     "other file",
			path:     "path/to/other.txt",
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

func TestHex_Dependencies(t *testing.T) {
	manager := Hex{}
	testPath := filepath.Join("testdata", "mix.exs")

	dependencies, err := manager.Dependencies(testPath)
	assert.NoError(t, err)

	expected := []types.Dependency{
		{
			Manager: types.Hex,
			Name:    "phoenix",
			Version: "1.7.18",
			Dev:     false,
			Definition: types.Definition{
				Path:    testPath,
				RawLine: "{:phoenix, \"~> 1.7.18\"},",
				Line:    35,
			},
		},
		{
			Manager: types.Hex,
			Name:    "postgrex",
			Version: "0.0.0",
			Dev:     false,
			Definition: types.Definition{
				Path:    testPath,
				RawLine: "{:postgrex, \">= 0.0.0\"},",
				Line:    36,
			},
		},
		{
			Manager: types.Hex,
			Name:    "phoenix_live_reload",
			Version: "1.2",
			Dev:     false,
			Definition: types.Definition{
				Path:    testPath,
				RawLine: "{:phoenix_live_reload, \"~> 1.2\", only: :dev},",
				Line:    37,
			},
		},
	}

	assert.Equal(t, len(expected), len(dependencies))
	for i, exp := range expected {
		assert.Equal(t, exp.Name, dependencies[i].Name)
		assert.Equal(t, exp.Version, dependencies[i].Version)
		assert.Equal(t, exp.Dev, dependencies[i].Dev)
		assert.Equal(t, exp.Line, dependencies[i].Line)
		assert.Equal(t, exp.RawLine, dependencies[i].RawLine)
	}
}

func TestHex_LockfilePath(t *testing.T) {
	manager := Hex{}
	tests := []struct {
		name        string
		inputPath   string
		expectError bool
	}{
		{
			name:        "missing lockfile",
			inputPath:   "testdata/mix.exs",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lockfilePath, err := manager.LockfilePath(tt.inputPath)
			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, lockfilePath)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, lockfilePath)
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
			name:     "exact version",
			version:  "2.2.3",
			expected: "2.2.3",
		},
		{
			name:     "version with spaces",
			version:  " 2.2.3 ",
			expected: "2.2.3",
		},
		{
			name:     "version >=",
			version:  ">= 0.0.2",
			expected: "0.0.2",
		},
		{
			name:     "version ~>",
			version:  "~> 4.10",
			expected: "4.10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cleanVersion(tt.version)
			assert.Equal(t, tt.expected, result)
		})
	}
}
