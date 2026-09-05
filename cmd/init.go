package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// ponytail: stub for phase 1 (CLI skeleton) — real behavior (resolving
// packages, downloading binaries) lands in a later phase, see docs/roadmap.md.
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Set up what's required to run datasplice",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("init")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
