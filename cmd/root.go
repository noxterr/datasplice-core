package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "datasplice",
	Short: "Ingest, transform and export data via a YAML-defined flow",
}

// Execute runs the root command; main just calls this.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
