package config

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/depshubhq/depshub/pkg/types"
	"github.com/spf13/viper"
)

type Config struct {
	path   string
	config ConfigFile
}

type ConfigFile struct {
	Version       int            `mapstructure:"version"`
	Ignore        []string       `mapstructure:"ignore"`
	ManifestFiles []ManifestFile `mapstructure:"manifest_files"`
}

type Rule struct {
	Name     string      `mapstructure:"name"`
	Disabled bool        `mapstructure:"disabled"`
	Value    any         `mapstructure:"value"`
	Level    types.Level `mapstructure:"level"`
}

type ManifestFile struct {
	Filter   string   `mapstructure:"filter"`
	Rules    []Rule   `mapstructure:"rules"`
	Packages []string `mapstructure:"packages"`
}

func New(filePath string) (Config, error) {
	folder := filepath.Dir(filePath)

	viper.AddConfigPath(folder)
	viper.SetConfigName("depshub")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()

	if err != nil {
		if errors.As(err, &viper.ConfigFileNotFoundError{}) && filePath != "." {
			return Config{}, err
		}

		if errors.As(err, &viper.ConfigParseError{}) {
			return Config{}, err
		}
	}

	c := Config{path: filePath}

	err = viper.Unmarshal(&c.config)

	if err != nil {
		return Config{}, err
	}

	return c, nil
}

// Checks if a path is ignored by the config
func (c Config) Ignored(path string) (bool, error) {
	ignored := false

	for _, ignore := range c.config.Ignore {
		matched, err := doublestar.Match(ignore, path)

		if err != nil {
			return false, err
		}

		if matched {
			ignored = true
			break
		}
	}

	return ignored, nil
}

func (c Config) Apply(manifestPath string, packageName string, rule types.Rule) error {
	// Reset the to the default state before applying any settings
	rule.Reset()

	// Iterate through manifest files in config
	for _, mf := range c.config.ManifestFiles {
		// Check if manifest path matches the filter
		matched, err := doublestar.Match(mf.Filter, manifestPath)
		if err != nil {
			return fmt.Errorf("invalid filter pattern %q: %w", mf.Filter, err)
		}
		if !matched {
			continue
		}

		// Check if package is in the packages list (if specified)
		if len(mf.Packages) > 0 {
			packageMatch := false
			for _, pkg := range mf.Packages {
				if pkg == packageName {
					packageMatch = true
					break
				}
			}
			if !packageMatch {
				continue
			}
		}

		// Look for matching rule by name
		for _, configRule := range mf.Rules {
			if configRule.Name == rule.GetName() {
				if configRule.Disabled {
					rule.SetLevel(types.LevelDisabled)
				}

				// Apply level if specified
				if configRule.Level != "" {
					rule.SetLevel(configRule.Level)
				}

				// Apply value if specified
				if configRule.Value != nil {
					if err := rule.SetValue(configRule.Value); err != nil {
						return fmt.Errorf("failed to set value for rule %q: %w", rule.GetName(), err)
					}
				}

				break
			}
		}
	}

	return nil
}
