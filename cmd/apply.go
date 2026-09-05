package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// ponytail: stub for phase 1 — real execution (spawning packages, running
// the pipeline) lands in phase 7, see docs/roadmap.md.
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply the flow described by the plan",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("apply")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)
}
