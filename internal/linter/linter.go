package linter

import (
	"fmt"

	"github.com/depshubhq/depshub/internal/linter/rules"
	"github.com/depshubhq/depshub/pkg/manager"
)

type Linter struct {
	rules []rules.Rule
}

func New() Linter {
	return Linter{
		rules: []rules.Rule{
			rules.NewRuleSorted(),
			rules.NewRuleNoAnyTag(),
		},
	}
}

func (l Linter) Run(path string) (mistakes []rules.Mistake, err error) {
	scanner := manager.New()
	manifests, err := scanner.Scan(path)

	if err != nil {
		fmt.Printf("Error: %s", err)
	}

	for _, rule := range l.rules {
		m, err := rule.Check(manifests)

		if err != nil {
			return nil, err
		}

		mistakes = append(mistakes, m...)
	}

	return mistakes, nil
}
