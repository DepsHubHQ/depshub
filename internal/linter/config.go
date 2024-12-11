package linter

import (
	"github.com/spf13/viper"
)

type config struct {
	Version       int            `mapstructure:"version"`
	Rules         []Rule         `mapstructure:"rules"`
	ManifestFiles []ManifestFile `mapstructure:"manifest_files"`
}

type Rule struct {
	Name     string `mapstructure:"name"`
	Disabled bool   `mapstructure:"disabled"`
	Value    int    `mapstructure:"value"`
	Level    string `mapstructure:"level"`
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

var Config = config{}

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
