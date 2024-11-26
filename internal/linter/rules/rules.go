package rules

import "github.com/depshubhq/depshub/pkg/types"

type Rule interface {
	GetName() string
	GetMessage() string
	Check([]types.Manifest) ([]Mistake, error)
}

type Mistake struct {
	Rule Rule
	Path string
	types.Definition
}
