package main

import (
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

		linter.Run(p)
	},
}
