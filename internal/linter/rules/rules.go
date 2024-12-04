package rules

import "github.com/depshubhq/depshub/pkg/types"

type Level int

const (
	LevelError Level = iota
	LevelWarning
)

type Rule interface {
	GetName() string
	GetMessage() string
	GetLevel() Level
	Check([]types.Manifest, types.PackagesInfo) ([]Mistake, error)
}

type Mistake struct {
	Rule        Rule
	Definitions []types.Definition
}
