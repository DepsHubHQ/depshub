package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/depshubhq/depshub/internal/linter/rules"
	"github.com/depshubhq/depshub/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewConfig tests various scenarios for config loading
func TestNewConfig(t *testing.T) {
	t.Run("valid config file", func(t *testing.T) {
		// Create a temporary config file
		dir := t.TempDir()
		configPath := filepath.Join(dir, "depshub.yaml")
		configContent := []byte(`
version: 1
ignore:
    - "test-ignore"
manifest_files:
  - filter: "*.lock"
    rules:
      - name: "test-rule"
        level: "warning"
        value: 42
    packages:
      - "test-package"
`)
		err := os.WriteFile(configPath, configContent, 0644)
		require.NoError(t, err)

		// Test loading the config
		config, err := New(configPath)
		require.NoError(t, err)
		assert.Equal(t, 1, config.config.Version)
		assert.Len(t, config.config.ManifestFiles, 1)
		assert.Equal(t, "*.lock", config.config.ManifestFiles[0].Filter)
		assert.Equal(t, "test-rule", config.config.ManifestFiles[0].Rules[0].Name)
		assert.Equal(t, "test-ignore", config.config.Ignore[0])
	})

	t.Run("config file not found", func(t *testing.T) {
		_, err := New("nonexistent/path/depshub.yaml")
		assert.Error(t, err)
	})

	t.Run("invalid config file", func(t *testing.T) {
		dir := t.TempDir()
		configPath := filepath.Join(dir, "depshub.yaml")
		invalidContent := []byte(`invalid: yaml: content`)
		err := os.WriteFile(configPath, invalidContent, 0644)
		require.NoError(t, err)

		_, err = New(configPath)
		assert.Error(t, err)
	})
}

// mockRule implements the Rule interface for testing
type mockRule struct {
	name     string
	level    rules.Level
	value    int
	disabled bool
}

func (m *mockRule) Check(manifests []types.Manifest, info types.PackagesInfo) ([]rules.Mistake, error) {
	return nil, nil
}
func (m *mockRule) GetLevel() rules.Level                 { return m.level }
func (m *mockRule) GetMessage() string                    { return "mock message" }
func (m *mockRule) GetName() string                       { return m.name }
func (m *mockRule) IsSupported(mt types.ManagerType) bool { return true }
func (m *mockRule) SetLevel(l rules.Level)                { m.level = l }
func (m *mockRule) SetValue(v any) error                  { m.value = v.(int); return nil }

// TestApply tests the Apply method with various scenarios
func TestApply(t *testing.T) {
	t.Run("ignore matching paths", func(t *testing.T) {
		config := Config{
			config: ConfigFile{
				Ignore: []string{"vendor/**/*.lock"},
				ManifestFiles: []ManifestFile{
					{
						Filter: "*.lock",
						Rules: []Rule{
							{
								Name:  "test-rule",
								Level: rules.LevelWarning,
							},
						},
					},
				},
			},
		}
		rule := &mockRule{name: "test-rule", level: rules.LevelError}
		mistakes := []rules.Mistake{
			{
				Rule: rule,
				Definitions: []types.Definition{
					{Path: "vendor/some/path/test.lock", RawLine: "test content"},
				},
			},
		}
		result := config.Apply(mistakes)
		assert.Len(t, result, 0) // Ignored paths should be filtered out
	})

	t.Run("apply rules with matching filter", func(t *testing.T) {
		config := Config{
			config: ConfigFile{
				ManifestFiles: []ManifestFile{
					{
						Filter: "*.lock",
						Rules: []Rule{
							{
								Name:  "test-rule",
								Level: rules.LevelWarning,
								Value: 42,
							},
						},
					},
				},
			},
		}

		rule := &mockRule{name: "test-rule", level: rules.LevelError}
		mistakes := []rules.Mistake{
			{
				Rule: rule,
				Definitions: []types.Definition{
					{Path: "test.lock", RawLine: "test content"},
				},
			},
		}

		result := config.Apply(mistakes)
		assert.Len(t, result, 1)
		assert.Equal(t, rules.LevelWarning, result[0].Rule.GetLevel())
	})

	t.Run("filter by package name", func(t *testing.T) {
		config := Config{
			config: ConfigFile{
				ManifestFiles: []ManifestFile{
					{
						Filter:   "*.lock",
						Packages: []string{"test-package"},
						Rules: []Rule{
							{
								Name:  "test-rule",
								Level: rules.LevelWarning,
							},
						},
					},
				},
			},
		}

		rule := &mockRule{name: "test-rule", level: rules.LevelError}
		mistakes := []rules.Mistake{
			{
				Rule: rule,
				Definitions: []types.Definition{
					{Path: "test.lock", RawLine: "test-package"},
				},
			},
		}

		result := config.Apply(mistakes)
		assert.Len(t, result, 1)
		assert.Equal(t, rules.LevelWarning, result[0].Rule.GetLevel())
	})

	t.Run("disable rule", func(t *testing.T) {
		config := Config{
			config: ConfigFile{
				ManifestFiles: []ManifestFile{
					{
						Filter: "*.lock",
						Rules: []Rule{
							{
								Name:     "test-rule",
								Disabled: true,
							},
						},
					},
				},
			},
		}

		rule := &mockRule{name: "test-rule", level: rules.LevelError}
		mistakes := []rules.Mistake{
			{
				Rule: rule,
				Definitions: []types.Definition{
					{Path: "test.lock", RawLine: "test content"},
				},
			},
		}

		result := config.Apply(mistakes)
		assert.Len(t, result, 0) // Disabled rules should be filtered out
	})
}

// TestFilterDisabledRules tests the filterDisabledRules function
func TestFilterDisabledRules(t *testing.T) {
	rule1 := &mockRule{name: "rule1", level: rules.LevelError}
	rule2 := &mockRule{name: "rule2", level: rules.LevelDisabled}
	rule3 := &mockRule{name: "rule3", level: rules.LevelWarning}

	mistakes := []rules.Mistake{
		{Rule: rule1},
		{Rule: rule2},
		{Rule: rule3},
	}

	filtered := filterDisabledRules(mistakes)
	assert.Len(t, filtered, 2)
	assert.Equal(t, rules.LevelError, filtered[0].Rule.GetLevel())
	assert.Equal(t, rules.LevelWarning, filtered[1].Rule.GetLevel())
}
