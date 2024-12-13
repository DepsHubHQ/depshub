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
		configPath, err := cmd.Flags().GetString("config")

		if err != nil {
			fmt.Println(err)
			return
		}
		config, err := linter.NewConfig(configPath)

		if err != nil {
			fmt.Println(err)
			return
		}

		var p = "."

		if len(args) > 0 {
			p = args[0]
		}

		lint := linter.New()
		mistakes, err := lint.Run(p)

		if err != nil {
			fmt.Printf("Error: %s", err)
			return
		}

		mistakes = config.Apply(mistakes)

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

			fmt.Printf("Found %s and %s:\n", e, w)
		} else if errorsCount != 0 {
			e := errors.Render(fmt.Sprintf("Found %d %s", errorsCount, pluralizedError))

			fmt.Printf("%s:\n", e)
		} else if warningsCount != 0 {
			w := warnings.Render(fmt.Sprintf("Found %d %s", warningsCount, pluralizedWarning))

			fmt.Printf("%s:\n", w)
		}

		for _, mistake := range mistakes {
			if mistake.Rule.GetLevel() == rules.LevelDisabled {
				continue
			}

			name := fmt.Sprintf("[%s]", mistake.Rule.GetName())

			if mistake.Rule.GetLevel() == rules.LevelError {
				name = errors.Render(fmt.Sprintf("[%s]", mistake.Rule.GetName()))
			} else {
				name = warnings.Render(fmt.Sprintf("[%s]", mistake.Rule.GetName()))
			}

			fmt.Printf("\n - %s - %s \n", name, mistake.Rule.GetMessage())

			var style = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("8"))

			lineNumberStyle := lipgloss.Color("8")

			for _, definition := range mistake.Definitions {
				path := lipgloss.NewStyle().
					Foreground(lineNumberStyle).
					Render(fmt.Sprintf("%s", definition.Path))

				rawLineStyle := lipgloss.Color("110")
				rawLine := lipgloss.NewStyle().Align(lipgloss.Center).Foreground(rawLineStyle).Render(definition.RawLine)

				lineNumber := lipgloss.NewStyle().
					Foreground(lineNumberStyle).
					Render(fmt.Sprintf("%d |", definition.Line))

				if definition.Line == 0 {
					fmt.Println(style.Render(fmt.Sprintf(" %s", path)))
				} else {
					fmt.Println(style.Render(fmt.Sprintf(" %s\n\n %s %s", path, lineNumber, rawLine)))
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
