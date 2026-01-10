package cobra

import (
	"time"

	"github.com/spf13/cobra"
	gitlab "gitlab.com/gitlab-org/api/client-go"

	"github.com/kilianpaquier/gitlab-storage-cleaner/internal/artifacts/engine"
	artifacts "github.com/kilianpaquier/gitlab-storage-cleaner/internal/artifacts/engine/v2"
)

var (
	dryRun            bool
	paths             []string
	server            string
	thresholdDuration time.Duration
	token             string

	cleanCmd = &cobra.Command{
		Use:   "artifacts",
		Short: "Clean artifacts of provided project(s)' gitlab storage",
		RunE: func(cmd *cobra.Command, _ []string) error {
			// check gitlab client
			client, err := gitlab.NewClient(token, gitlab.WithBaseURL(server), gitlab.WithoutRetries())
			if err != nil {
				return err
			}

			opts := []engine.RunOption{
				engine.WithDryRun(dryRun),
				engine.WithLogger(engine.NewSlogLogger(logger)),
				engine.WithPaths(paths...),
				engine.WithThresholdDuration(thresholdDuration),
			}

			return artifacts.Run(cmd.Context(), client, opts...)
		},
	}
)

func init() {
	rootCmd.AddCommand(cleanCmd)

	// gitlab token
	cleanCmd.Flags().StringVar(&token, "token", "", "gitlab read/write token with maintainer rights to delete artifacts")
	_ = cleanCmd.MarkFlagRequired("token")

	// gitlab server
	cleanCmd.Flags().StringVar(&server, "server", "", "gitlab server host")
	_ = cleanCmd.MarkFlagRequired("server")

	// threshold duration
	cleanCmd.Flags().DurationVar(&thresholdDuration, "threshold-duration", engine.DefaultThresholdDuration,
		"threshold duration (positive) where, jobs older than command execution time minus this threshold will be deleted")

	// dry run
	cleanCmd.Flags().BoolVar(&dryRun, "dry-run", false, "truthy if run must not delete jobs' artifacts but only list matched projects")

	// projects filtering options
	cleanCmd.Flags().StringSliceVar(&paths, "paths", []string{}, "list of valid regexps to match project path (with namespace)")
	_ = cleanCmd.MarkFlagRequired("paths")
}
