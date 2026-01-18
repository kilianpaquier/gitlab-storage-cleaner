package cobra

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	gitlab "gitlab.com/gitlab-org/api/client-go"

	"github.com/kilianpaquier/gitlab-storage-cleaner/internal/artifacts/engine"
	artifacts "github.com/kilianpaquier/gitlab-storage-cleaner/internal/artifacts/engine/v2"
)

const envPrefix = "cleaner-"

const (
	flagDryRun            = "dry-run"
	flagPaths             = "paths"
	flagServer            = "server"
	flagThresholdDuration = "threshold-duration"
	flagToken             = "token"
)

// artifactsCmd creates a new cobra command for cleaning GitLab artifacts.
func artifactsCmd() *cobra.Command { //nolint:gocognit,funlen
	var (
		token  string
		server string
		dryRun bool
		paths  []string
	)
	thresholdDuration := 7 * 24 * time.Hour

	cmd := &cobra.Command{
		Use:   "artifacts",
		Short: "Clean artifacts of provided project(s)' gitlab storage",
		Args: func(cmd *cobra.Command, _ []string) (err error) {
			// validate dry run environment variable
			if !cmd.Flags().Changed(flagDryRun) {
				if env := getenv(envPrefix + flagDryRun); env != "" {
					dr, err := strconv.ParseBool(env)
					if err != nil {
						return fmt.Errorf(`invalid argument %q for "--%s" flag: %w`, env, flagDryRun, err)
					}
					dryRun = dr
				}
			}

			if !cmd.Flags().Changed(flagPaths) {
				if env := getenv(envPrefix + flagPaths); env != "" {
					paths = strings.Split(env, ",")
				}
			}

			// validate threshold duration environment variable
			if !cmd.Flags().Changed(flagThresholdDuration) {
				if env := getenv(envPrefix + flagThresholdDuration); env != "" {
					td, err := time.ParseDuration(env)
					if err != nil {
						return fmt.Errorf(`invalid argument %q for "--%s" flag: %w`, env, flagThresholdDuration, err)
					}
					thresholdDuration = td
				}
			}

			var missings []string
			if len(paths) == 0 {
				missings = append(missings, `"`+flagPaths+`"`)
			}
			if server == "" {
				missings = append(missings, `"`+flagServer+`"`)
			}
			if token == "" {
				missings = append(missings, `"`+flagToken+`"`)
			}
			if len(missings) > 0 {
				return fmt.Errorf("required flag(s) %s not set", strings.Join(missings, ", "))
			}
			return nil
		},
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

	// gitlab token
	cmd.Flags().StringVar(&token, flagToken, coalesce(os.Getenv("GITLAB_TOKEN"), os.Getenv("GL_TOKEN")), "gitlab read/write token with maintainer rights to delete artifacts")

	// gitlab server
	cmd.Flags().StringVar(&server, flagServer, coalesce(os.Getenv("CI_API_V4_URL"), os.Getenv("CI_SERVER_HOST")), "gitlab server host")

	// dry run
	cmd.Flags().BoolVar(&dryRun, flagDryRun, false, "truthy if run must not delete jobs' artifacts but only list matched projects")

	// projects filtering options
	cmd.Flags().StringSliceVar(&paths, flagPaths, nil, "list of valid regexps to match project path (with namespace)")

	// threshold duration
	cmd.Flags().DurationVar(&thresholdDuration, flagThresholdDuration, thresholdDuration,
		"threshold duration (positive) where, jobs older than command execution time minus this threshold will be deleted")

	return cmd
}
