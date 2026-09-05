package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// ponytail: stub for phase 1 — the macro-area report lands in phase 5,
// see docs/roadmap.md and docs/schema.md#plan-output.
var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Show what apply would do, without doing it",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("plan")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(planCmd)
}
