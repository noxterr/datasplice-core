package cmd

import (
	"fmt"

	"github.com/datasplice-labs/datasplice/internal/config"
	"github.com/datasplice-labs/datasplice/internal/pipeline"
	"github.com/spf13/cobra"
)

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Show what run would do, without doing it",
	RunE: func(cmd *cobra.Command, args []string) error {
		m, v, err := config.Load(MainFile, VariablesFile)
		if err != nil {
			return err
		}
		if _, err := config.ResolveSecrets(m, v); err != nil {
			return err
		}

		fmt.Printf("datasplice plan: %q\n\n", m.Name)
		for i, line := range pipeline.Describe(m) {
			fmt.Printf("  %d. %s\n", i+1, line)
		}
		fmt.Printf("\n%d steps. Run `datasplice run` to execute.\n", len(m.Steps))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(planCmd)
}
