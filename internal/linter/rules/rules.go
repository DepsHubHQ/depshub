package rules

import (
	"errors"

	"github.com/depshubhq/depshub/pkg/types"
)

type Level string

const (
	LevelError    Level = "error"
	LevelWarning  Level = "warning"
	LevelDisabled Level = "disabled"
)

type Rule interface {
	Check([]types.Manifest, types.PackagesInfo) ([]Mistake, error)
	GetLevel() Level
	GetMessage() string
	GetName() string
	IsSupported(types.ManagerType) bool
	SetLevel(Level)
	SetValue(any) error
}

type Mistake struct {
	Rule        Rule
	Definitions []types.Definition
}

var ErrInvalidRuleValue = errors.New("invalid rule value")
