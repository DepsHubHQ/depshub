package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
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
			errorsStyle := lipgloss.Color("9")
			errors := lipgloss.NewStyle().
				Foreground(errorsStyle).
				Render(fmt.Sprintf("%d errors found", len(mistakes)))

			fmt.Printf("%s:\n\n", errors)

			for _, mistake := range mistakes {
				fmt.Printf("- %s - %s \n\n", mistake.Rule.GetName(), mistake.Rule.GetMessage())

				pStyle := lipgloss.Color("86")
				p := lipgloss.NewStyle().
					Foreground(pStyle).
					Render(mistake.Path)

				lineNumberStyle := lipgloss.Color("5")
				lineNumber := lipgloss.NewStyle().
					Foreground(lineNumberStyle).
					Render(fmt.Sprintf("%d", mistake.Line))

				rawLineStyle := lipgloss.Color("3")
				rawLine := lipgloss.NewStyle().Align(lipgloss.Center).Foreground(rawLineStyle).Render(mistake.RawLine)

				fmt.Printf("%s \n", p)
				fmt.Printf("   %s %s\n\n", lineNumber, rawLine)
			}
		}
	},
}
