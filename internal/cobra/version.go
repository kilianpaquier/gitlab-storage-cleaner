package cobra

import "github.com/spf13/cobra"

var (
	// version is substituted with -ldflags
	version = "v0.0.0"

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Show current version",
		Run:   func(_ *cobra.Command, _ []string) { logger.Info(version) },
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}
