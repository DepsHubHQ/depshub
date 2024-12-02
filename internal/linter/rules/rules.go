package rules

import "github.com/depshubhq/depshub/pkg/types"

type Level int

const (
	LevelError Level = iota
	LevelWarning
)

type PackagesInfo = map[string]types.Package

type Rule interface {
	GetName() string
	GetMessage() string
	GetLevel() Level
	Check([]types.Manifest, PackagesInfo) ([]Mistake, error)
}

type Mistake struct {
	Rule        Rule
	Definitions []types.Definition
}
