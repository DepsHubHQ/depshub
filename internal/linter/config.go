package linter

import (
	"github.com/depshubhq/depshub/internal/linter/rules"
	"github.com/spf13/viper"
)

type ConfigType struct {
	Version       int            `mapstructure:"version"`
	Rules         []Rule         `mapstructure:"rules"`
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
		// First, apply global rules
		for _, rule := range Config.Rules {
			if mistake.Rule.GetName() == rule.Name {
				mistake.Rule.SetLevel(rule.Level)

				if rule.Disabled {
					mistake.Rule.SetLevel(rules.LevelDisabled)
				}
			}
		}
	}

	// Filter out disabled rules

	mistakes = filterDisabledRules(mistakes)

	return mistakes
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
