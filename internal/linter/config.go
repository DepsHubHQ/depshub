package linter

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/depshubhq/depshub/internal/linter/rules"
	"github.com/spf13/viper"
)

type ConfigFile struct {
	Version       int            `mapstructure:"version"`
	Ignore        []string       `mapstructure:"ignore"`
	ManifestFiles []ManifestFile `mapstructure:"manifest_files"`
}

type Rule struct {
	Name     string      `mapstructure:"name"`
	Disabled bool        `mapstructure:"disabled"`
	Value    int         `mapstructure:"value"`
	Level    rules.Level `mapstructure:"level"`
}

type ManifestFile struct {
	Filter   string   `mapstructure:"filter"`
	Rules    []Rule   `mapstructure:"rules"`
	Packages []string `mapstructure:"packages"`
}

type Config struct {
	path   string
	config ConfigFile
}

func NewConfig(filePath string) (Config, error) {
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

func (c Config) Apply(mistakes []rules.Mistake) []rules.Mistake {
	for _, mistake := range mistakes {
		ignored := false

		for _, ignore := range c.config.Ignore {
			matched, err := doublestar.Match(ignore, mistake.Definitions[0].Path)

			if err != nil {
				fmt.Println(err)
				continue
			}

			if matched {
				ignored = true
				break
			}
		}

		if ignored {
			mistake.Rule.SetLevel(rules.LevelDisabled)
			continue
		}

		for _, configManifestFile := range c.config.ManifestFiles {
			matched, err := doublestar.Match(configManifestFile.Filter, mistake.Definitions[0].Path)

			if err != nil {
				fmt.Println(err)
				continue
			}

			matchByPackageName := false
			if len(configManifestFile.Packages) > 0 {
				for _, p := range configManifestFile.Packages {
					// We should probably include the package information in the mistake struct,
					// instead of just checking the raw line.
					if strings.Contains(mistake.Definitions[0].RawLine, p) {
						matched = true
						matchByPackageName = true
						break
					}
				}
			} else {
				matchByPackageName = true
			}

			if !matched {
				continue
			}

			for _, rule := range configManifestFile.Rules {
				if mistake.Rule.GetName() == rule.Name && matchByPackageName {
					mistake.Rule.SetLevel(rule.Level)

					if rule.Disabled {
						mistake.Rule.SetLevel(rules.LevelDisabled)
						continue
					}

					err := mistake.Rule.SetValue(rule.Value)

					if err != nil {
						fmt.Println(err)
					}
				}
			}
		}
	}

	return filterDisabledRules(mistakes)
}

func filterDisabledRules(mistakes []rules.Mistake) []rules.Mistake {
	var filteredMistakes []rules.Mistake

	for _, mistake := range mistakes {
		if mistake.Rule.GetLevel() != rules.LevelDisabled {
			filteredMistakes = append(filteredMistakes, mistake)
		}
	}
	return filteredMistakes
}
