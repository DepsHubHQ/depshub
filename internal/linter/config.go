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

type ConfigType struct {
	Version       int            `mapstructure:"version"`
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

var Config = ConfigType{}

func InitConfig(filePath string) error {
	folder := filepath.Dir(filePath)

	viper.AddConfigPath(folder)
	viper.SetConfigName("depshub")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()

	if err != nil {
		if errors.As(err, &viper.ConfigFileNotFoundError{}) && filePath != "." {
			return err
		}

		if errors.As(err, &viper.ConfigParseError{}) {
			return err
		}
	}

	return viper.Unmarshal(&Config)
}

func ApplyConfig(mistakes []rules.Mistake) []rules.Mistake {
	for _, mistake := range mistakes {
		for _, configManifestFile := range Config.ManifestFiles {
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
