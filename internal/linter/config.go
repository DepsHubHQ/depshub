package linter

import (
	"fmt"

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
	Glob     string    `mapstructure:"glob"`
	Rules    []Rule    `mapstructure:"rules"`
	Packages []Package `mapstructure:"packages"`
}

type Package struct {
	Name  string `mapstructure:"name"`
	Rules []Rule `mapstructure:"rules"`
}

var Config = ConfigType{}

func InitConfig() error {
	viper.AddConfigPath(".")
	viper.SetConfigName("depshub")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()

	if err != nil {
		return err
	}

	return viper.Unmarshal(&Config)
}

func ApplyConfig(mistakes []rules.Mistake) []rules.Mistake {
	for _, mistake := range mistakes {
		for _, configManifestFile := range Config.ManifestFiles {
			matched, err := doublestar.Match(configManifestFile.Glob, mistake.Definitions[0].Path)

			if err != nil {
				fmt.Println(err)
				continue
			}

			if !matched {
				continue
			}

			for _, rule := range configManifestFile.Rules {
				if mistake.Rule.GetName() == rule.Name {
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
