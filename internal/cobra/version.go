package cobra

import (
	"github.com/spf13/cobra"

	"github.com/kilianpaquier/gitlab-storage-cleaner/internal/build"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show current version",
	Run:   func(_ *cobra.Command, _ []string) { logger.Info(build.GetInfo().String()) },
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
