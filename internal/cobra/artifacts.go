package cobra

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"

	"github.com/kilianpaquier/gitlab-storage-cleaner/internal/artifacts/v2"
)

var (
	cleanOpts = artifacts.Options{}

	server string
	token  string

	cleanCmd = &cobra.Command{
		Use:   "artifacts",
		Short: "Clean artifacts of provided project(s)' gitlab storage",
		Run: func(cmd *cobra.Command, _ []string) {
			ctx := cmd.Context()

			// check gitlab client
			client, err := gitlab.NewClient(token, gitlab.WithBaseURL(server), gitlab.WithoutRetries())
			if err != nil {
				fatal(ctx, err)
			}

			// ensure options are all here
			if err := cleanOpts.EnsureDefaults(); err != nil {
				fatal(ctx, err)
			}

			// run artifacts clean command
			if err := artifacts.Run(ctx, client, cleanOpts); err != nil {
				fatal(ctx, err)
			}
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
	cleanCmd.Flags().DurationVar(
		&cleanOpts.ThresholdDuration, "threshold-duration", 7*24*time.Hour,
		"threshold duration (positive) where, from now, jobs artifacts expiration is after will be cleaned up",
	)

	// threshold size
	cleanCmd.Flags().Uint64Var(
		&cleanOpts.ThresholdSize, "threshold-size", 1000000,
		"threshold size (in bytes) where jobs artifacts size sum is bigger will be cleaned up",
	)

	// dry run
	cleanCmd.Flags().BoolVar(
		&cleanOpts.DryRun, "dry-run", false,
		"truthy if run must not delete jobs' artifacts but only list matched projects")

	// projects filtering options
	cleanCmd.Flags().StringSliceVar(
		&cleanOpts.Paths, "paths", []string{},
		"list of valid regexps to match project path (with namespace)")
	_ = cleanCmd.MarkFlagRequired("paths")
}
