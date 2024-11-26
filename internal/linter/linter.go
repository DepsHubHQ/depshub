package linter

import (
	"fmt"
	"github.com/depshubhq/depshub/pkg/manager"
)

func Run(path string) {
	fmt.Printf("Linting your project... %s\n", path)

	scanner := manager.NewScanner()
	res := scanner.Scan(path)

	fmt.Println(res)
}
