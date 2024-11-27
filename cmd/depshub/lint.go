package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/depshubhq/depshub/internal/linter"
	"github.com/depshubhq/depshub/internal/linter/rules"
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

		errorsCount := 0
		warningsCount := 0

		for _, mistake := range mistakes {
			if mistake.Rule.GetLevel() == rules.LevelError {
				errorsCount++
			} else {
				warningsCount++
			}
		}

		errorsStyle := lipgloss.Color("9")
		errors := lipgloss.NewStyle().
			Foreground(errorsStyle)

		warningsStyle := lipgloss.Color("11")
		warnings := lipgloss.NewStyle().
			Foreground(warningsStyle)

		pluralizedError := pluralize(errorsCount, "error", "errors")
		pluralizedWarning := pluralize(warningsCount, "warning", "warnings")

		if errorsCount != 0 && warningsCount != 0 {
			e := errors.Render(fmt.Sprintf("%d %s", errorsCount, pluralizedError))
			w := warnings.Render(fmt.Sprintf("%d %s", warningsCount, pluralizedWarning))

			fmt.Printf("%s and %s found:\n", e, w)
		} else if errorsCount != 0 {
			e := errors.Render(fmt.Sprintf("%d %s found", errorsCount, pluralizedError))

			fmt.Printf("%s:\n", e)
		} else if warningsCount != 0 {
			w := warnings.Render(fmt.Sprintf("%d %s found", warningsCount, pluralizedWarning))

			fmt.Printf("%s:\n", w)
		}

		mistakesMap := make(map[string][]rules.Mistake)

		for _, mistake := range mistakes {
			mistakesMap[mistake.Path] = append(mistakesMap[mistake.Path], mistake)
		}

		for path, mistakes := range mistakesMap {
			pStyle := lipgloss.Color("86")
			p := lipgloss.NewStyle().
				Foreground(pStyle).
				Render(path)

			fmt.Printf("\n %s", p)

			for _, mistake := range mistakes {
				name := fmt.Sprintf("[%s]", mistake.Rule.GetName())

				if mistake.Rule.GetLevel() == rules.LevelError {
					name = errors.Render(fmt.Sprintf("[%s]", mistake.Rule.GetName()))
				} else {
					name = warnings.Render(fmt.Sprintf("[%s]", mistake.Rule.GetName()))
				}

				fmt.Printf("\n - %s - %s \n", name, mistake.Rule.GetMessage())

				if mistake.Definition != nil {
					rawLineStyle := lipgloss.Color("110")
					rawLine := lipgloss.NewStyle().Align(lipgloss.Center).Foreground(rawLineStyle).Render(mistake.RawLine)

					var style = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("8"))
					lineNumberStyle := lipgloss.Color("8")
					lineNumber := lipgloss.NewStyle().
						Foreground(lineNumberStyle).
						Render(fmt.Sprintf("%d", mistake.Line))
					fmt.Println(style.Render(fmt.Sprintf("   %s %s", lineNumber, rawLine)))
				}
			}
		}
	},
}

func pluralize(n int, singular, plural string) string {
	if n == 1 {
		return singular
	}
	return plural
}
