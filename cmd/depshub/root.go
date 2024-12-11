package main

import (
	"fmt"
	"os"

	"github.com/depshubhq/depshub/internal/linter"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "depshub",
	Short: "DepsHub is a tool to manage your dependencies",
	Long: `DepsHub is a tool to manage your dependencies.
It helps you to keep track of your dependencies, 
and to update them when new versions are available.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func Execute() {
	rootCmd.Version = fmt.Sprintf("%s", version)

	if err := linter.InitConfig(); err != nil {
		fmt.Println(err)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
