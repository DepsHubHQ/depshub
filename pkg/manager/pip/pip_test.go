package pipmanager

import (
	"github.com/depshubhq/depshub/pkg/types"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestPip_GetType(t *testing.T) {
	manager := Pip{}
	assert.Equal(t, types.Pip, manager.GetType())
}

func TestPip_Managed(t *testing.T) {
	manager := Pip{}
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "requirements.txt file",
			path:     "path/to/requirements.txt",
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
	manager := Pip{}
	testPath := filepath.Join("testdata", "requirements.txt")

	dependencies, err := manager.Dependencies(testPath)
	assert.NoError(t, err)

	expected := []types.Dependency{
		{
			Manager: types.Pip,
			Name:    "Flask",
			Version: "2.2.3",
			Dev:     false,
			Definition: types.Definition{
				Path:    testPath,
				RawLine: "Flask==2.2.3",
				Line:    1,
			},
		},
		{
			Manager: types.Pip,
			Name:    "requests",
			Version: "2.28.1",
			Dev:     false,
			Definition: types.Definition{
				Path:    testPath,
				RawLine: "requests>=2.28.1",
				Line:    2,
			},
		},
		{
			Manager: types.Pip,
			Name:    "numpy",
			Version: "",
			Dev:     false,
			Definition: types.Definition{
				Path:    testPath,
				RawLine: "numpy",
				Line:    3,
			},
		},
		{
			Manager: types.Pip,
			Name:    "pandas",
			Version: "1.5.3",
			Dev:     false,
			Definition: types.Definition{
				Path:    testPath,
				RawLine: "pandas<=1.5.3 # test comment",
				Line:    4,
			},
		},
		{
			Manager: types.Pip,
			Name:    "gunicorn",
			Version: "20.1.0",
			Dev:     false,
			Definition: types.Definition{
				Path:    testPath,
				RawLine: "gunicorn==20.1.0",
				Line:    5,
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
	manager := Pip{}
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
			name:     "version with comment",
			version:  "2.2.3 # latest stable",
			expected: "2.2.3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cleanVersion(tt.version)
			assert.Equal(t, tt.expected, result)
		})
	}
}
