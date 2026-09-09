package cmd

import (
	"fmt"

	"github.com/datasplice-labs/datasplice/internal/config"
	"github.com/datasplice-labs/datasplice/internal/pipeline"
	"github.com/spf13/cobra"
)

// ponytail: no separate `init` — there's nothing to prep that `run` can't
// do itself. Package resolution/installation (docs/roadmap.md phase 6)
// belongs right here, as one pass over every step before any row
// processing starts — so a missing package fails clean before anything's
// been written, instead of mid-pipeline.
var runCmd = &cobra.Command{
	Use:   "run",
	Short: fmt.Sprintf("Run the flow described by %s", MainFile),
	RunE: func(cmd *cobra.Command, args []string) error {
		m, v, err := config.Load(MainFile, VariablesFile)
		if err != nil {
			return err
		}
		if _, err := config.ResolveSecrets(m, v); err != nil {
			return err
		}
		return pipeline.Run(m)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
