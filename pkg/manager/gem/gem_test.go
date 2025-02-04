package gem

import (
	"path/filepath"

	"github.com/depshubhq/depshub/pkg/types"
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestPip_GetType(t *testing.T) {
	manager := Gem{}
	assert.Equal(t, types.Gem, manager.GetType())
}

func TestPip_Managed(t *testing.T) {
	manager := Gem{}
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "Gemfile file",
			path:     "path/to/Gemfile",
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

func TestPip_Dependencies(t *testing.T) {
	manager := Gem{}
	testPath := filepath.Join("testdata", "Gemfile")

	dependencies, err := manager.Dependencies(testPath)
	assert.NoError(t, err)

	expected := []types.Dependency{
		{
			Manager: types.Gem,
			Name:    "rails",
			Version: "6.0.2",
			Dev:     false,
			Definition: types.Definition{
				Path:    testPath,
				RawLine: "gem 'rails', '~> 6.0.2', '>= 6.0.2.2'",
				Line:    7,
			},
		},
		{
			Manager: types.Gem,
			Name:    "sqlite3",
			Version: "1.4",
			Dev:     false,
			Definition: types.Definition{
				Path:    testPath,
				RawLine: "gem 'sqlite3', '~> 1.4'",
				Line:    9,
			},
		},
		{
			Manager: types.Gem,
			Name:    "sass-rails",
			Version: "6",
			Dev:     false,
			Definition: types.Definition{
				Path:    testPath,
				RawLine: "gem 'sass-rails', '>= 6'",
				Line:    11,
			},
		},
		{
			Manager: types.Gem,
			Name:    "bootsnap",
			Version: "1.4.2",
			Dev:     false,
			Definition: types.Definition{
				Path:    testPath,
				RawLine: "gem 'bootsnap', '>= 1.4.2', require: false",
				Line:    21,
			},
		},
		{
			Manager: types.Gem,
			Name:    "byebug",
			Version: "",
			Dev:     true,
			Definition: types.Definition{
				Path:    testPath,
				RawLine: "gem 'byebug', platforms: [:mri, :mingw, :x64_mingw]",
				Line:    25,
			},
		},
		{
			Manager: types.Gem,
			Name:    "web-console",
			Version: "3.3.0",
			Dev:     true,
			Definition: types.Definition{
				Path:    testPath,
				RawLine: "gem 'web-console', '>= 3.3.0'",
				Line:    30,
			},
		},
		{
			Manager: types.Gem,
			Name:    "spring",
			Version: "",
			Dev:     true,
			Definition: types.Definition{
				Path:    testPath,
				RawLine: "gem 'spring'",
				Line:    32,
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

func TestPip_LockfilePath(t *testing.T) {
	manager := Gem{}
	tests := []struct {
		name        string
		inputPath   string
		expectError bool
	}{
		{
			name:        "missing lockfile",
			inputPath:   "testdata/requirements.txt",
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
