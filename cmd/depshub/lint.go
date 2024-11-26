package main

import (
	"fmt"

	"github.com/depshubhq/depshub/internal/linter"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(lintCmd)
}

var lintCmd = &cobra.Command{
	Use:   "lint [flags] [path]",
	Short: "Run the linter on your project",
	Long:  `Run the linter on your project to find issues in your code.`,
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		var p = "."

		if len(args) > 0 {
			p = args[0]
		}

		lint := linter.New()
		mistakes, err := lint.Run(p)

		if err != nil {
			fmt.Printf("Error: %s", err)
		}

		if len(mistakes) != 0 {
			fmt.Printf("Found %d mistakes:\n", len(mistakes))
			for _, mistake := range mistakes {
				fmt.Printf("- %s - %s \n\n", mistake.Rule.GetName(), mistake.Rule.GetMessage())
				fmt.Printf("   %s:\n", mistake.Path)
				fmt.Printf("   %d %s\n\n", mistake.Line, mistake.RawLine)
			}
		}
	},
}
