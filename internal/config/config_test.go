package config

import (
	"os"
	"path/filepath"
	"testing"

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
	level    types.Level
	value    int
	disabled bool
}

func (m *mockRule) Check(manifests []types.Manifest, info types.PackagesInfo, config types.Config) ([]types.Mistake, error) {
	return nil, nil
}
func (m *mockRule) GetLevel() types.Level                 { return m.level }
func (m *mockRule) GetMessage() string                    { return "mock message" }
func (m *mockRule) GetName() string                       { return m.name }
func (m *mockRule) IsSupported(mt types.ManagerType) bool { return true }
func (m *mockRule) SetLevel(l types.Level)                { m.level = l }
func (m *mockRule) SetValue(v any) error                  { m.value = v.(int); return nil }

// TODO add tests for the Apply function
