package pyproject

import (
	"path/filepath"
	"testing"

	"github.com/depshubhq/depshub/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestPyproject_GetType(t *testing.T) {
	manager := Pyproject{}
	assert.Equal(t, types.Pyproject, manager.GetType())
}

func TestPyproject_Managed(t *testing.T) {
	manager := Pyproject{}
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "requirements.txt file",
			path:     "path/to/pyproject.toml",
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

func TestPyproject_Dependencies(t *testing.T) {
	manager := Pyproject{}
	testPath := filepath.Join("testdata", "pyproject.toml")

	dependencies, err := manager.Dependencies(testPath)
	assert.NoError(t, err)

	expected := []types.Dependency{
		{
			Manager: types.Pyproject,
			Name:    "test",
			Version: "2.26.0",
			Dev:     false,
			Definition: types.Definition{
				Path:    testPath,
				RawLine: "test = \"2.26.0\"",
				Line:    18,
			},
		},
		{
			Manager: types.Pyproject,
			Name:    "test2",
			Version: "2.26.0",
			Dev:     false,
			Definition: types.Definition{
				Path:    testPath,
				RawLine: "test2 = \">=2.26.0\"",
				Line:    19,
			},
		},
		{
			Manager: types.Pyproject,
			Name:    "requests",
			Version: "2.26.0",
			Dev:     false,
			Definition: types.Definition{
				Path:    testPath,
				RawLine: "requests = \"^2.26.0\"",
				Line:    20,
			},
		},
		{
			Manager: types.Pyproject,
			Name:    "numpy",
			Version: "1.21",
			Dev:     false,
			Definition: types.Definition{
				Path:    testPath,
				RawLine: "numpy = \">=1.21,<2.0\"",
				Line:    21,
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
