package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// ponytail: stub for phase 1 — schema checks and the --fix/--verbose flags
// land in phase 3, see docs/roadmap.md.
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate configuration and files",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("validate")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}
