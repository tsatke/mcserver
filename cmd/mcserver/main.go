package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Version can be set with the Go linker.
	Version string = "main"
	// AppName is the name of this app, as displayed in the help
	// text of the root command.
	AppName = "mcserver"
)

var (
	rootCmd = &cobra.Command{
		Use:  AppName,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(os.Stdin, os.Stdout)
		},
		Version: Version,
	}
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
