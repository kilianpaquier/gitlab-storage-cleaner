package cobra //nolint:testpackage

import (
	"testing"
	"time"

	"github.com/kilianpaquier/gitlab-storage-cleaner/internal/testutils"
)

func TestArtifactsFlags(t *testing.T) {
	t.Run("missing_required", func(t *testing.T) {
		// Arrange
		t.Setenv("CI_API_V4_URL", "")
		t.Setenv("CI_SERVER_HOST", "")
		cmd := artifactsCmd()

		// Act
		err := cmd.Execute()

		// Assert
		testutils.Error(testutils.Require(t), err)
		testutils.Contains(t, err.Error(), `required flag(s) "paths", "server", "token" not set`)
	})

	t.Run("invalid_env", func(t *testing.T) {
		for _, env := range []string{"CLEANER_DRY_RUN", "CLEANER_THRESHOLD_DURATION"} {
			t.Run(env, func(t *testing.T) {
				// Arrange
				t.Setenv("CI_API_V4_URL", "https://gitlab.example.com/api/v4")
				t.Setenv("CLEANER_PATHS", "path1,path2")
				t.Setenv("GITLAB_TOKEN", "token")
				t.Setenv(env, "invalid")

				cmd := artifactsCmd()

				// Act
				err := cmd.Execute()

				// Assert
				testutils.Error(testutils.Require(t), err)
				testutils.Contains(t, err.Error(), `invalid argument "invalid"`)
			})
		}
	})

	t.Run("from_env", func(t *testing.T) {
		// Arrange
		t.Setenv("CI_API_V4_URL", "https://gitlab.example.com/api/v4")
		t.Setenv("CLEANER_DRY_RUN", "true")
		t.Setenv("CLEANER_PATHS", `^$CI_PROJECT_NAMESPACE\/.*$`)
		t.Setenv("CLEANER_THRESHOLD_DURATION", "72h")
		t.Setenv("GITLAB_TOKEN", "token")

		cmd := artifactsCmd()

		// Act
		err := cmd.Execute()

		// Assert
		testutils.NoError(testutils.Require(t), err)

		server, err := cmd.Flags().GetString(flagServer)
		testutils.NoError(testutils.Require(t), err)
		testutils.Equal(t, "https://gitlab.example.com/api/v4", server)

		token, err := cmd.Flags().GetString(flagToken)
		testutils.NoError(testutils.Require(t), err)
		testutils.Equal(t, "token", token)

		dryRun, err := cmd.Flags().GetBool(flagDryRun)
		testutils.NoError(testutils.Require(t), err)
		testutils.Equal(t, true, dryRun)

		thresholdDuration, err := cmd.Flags().GetDuration(flagThresholdDuration)
		testutils.NoError(testutils.Require(t), err)
		testutils.Equal(t, 72*time.Hour, thresholdDuration)

		paths, err := cmd.Flags().GetStringSlice(flagPaths)
		testutils.NoError(testutils.Require(t), err)
		testutils.Equal(testutils.Require(t), 1, len(paths))
		testutils.Equal(t, `^$CI_PROJECT_NAMESPACE\/.*$`, paths[0])
	})

	t.Run("from_env_alt", func(t *testing.T) {
		// Arrange
		t.Setenv("CI_SERVER_HOST", "https://gitlab.example.com")
		t.Setenv("CLEANER_PATHS", "path1,path2")
		t.Setenv("GL_TOKEN", "token")

		cmd := artifactsCmd()

		// Act
		err := cmd.Execute()

		// Assert
		testutils.NoError(testutils.Require(t), err)

		server, err := cmd.Flags().GetString(flagServer)
		testutils.NoError(testutils.Require(t), err)
		testutils.Equal(t, "https://gitlab.example.com", server)

		token, err := cmd.Flags().GetString(flagToken)
		testutils.NoError(testutils.Require(t), err)
		testutils.Equal(t, "token", token)

		paths, err := cmd.Flags().GetStringSlice(flagPaths)
		testutils.NoError(testutils.Require(t), err)
		testutils.Equal(testutils.Require(t), 2, len(paths))
		for i, p := range []string{"path1", "path2"} {
			testutils.Equal(t, p, paths[i])
		}
	})

	t.Run("flags_override_env", func(t *testing.T) {
		// Arrange
		t.Setenv("CI_API_V4_URL", "https://gitlab.example.com/api/v4")
		t.Setenv("CLEANER_DRY_RUN", "invalid")
		t.Setenv("CLEANER_PATHS", "path1,path2")
		t.Setenv("CLEANER_THRESHOLD_DURATION", "92h")
		t.Setenv("GITLAB_TOKEN", "token")

		cmd := artifactsCmd()
		cmd.SetArgs([]string{"--" + flagPaths, `^$CI_PROJECT_NAMESPACE\/.*$`, "--" + flagDryRun, "--" + flagThresholdDuration, "72h"})

		// Act
		err := cmd.Execute()

		// Assert
		testutils.NoError(testutils.Require(t), err)

		server, err := cmd.Flags().GetString(flagServer)
		testutils.NoError(testutils.Require(t), err)
		testutils.Equal(t, "https://gitlab.example.com/api/v4", server)

		token, err := cmd.Flags().GetString(flagToken)
		testutils.NoError(testutils.Require(t), err)
		testutils.Equal(t, "token", token)

		dryRun, err := cmd.Flags().GetBool(flagDryRun)
		testutils.NoError(testutils.Require(t), err)
		testutils.Equal(t, true, dryRun)

		thresholdDuration, err := cmd.Flags().GetDuration(flagThresholdDuration)
		testutils.NoError(testutils.Require(t), err)
		testutils.Equal(t, 72*time.Hour, thresholdDuration)

		paths, err := cmd.Flags().GetStringSlice(flagPaths)
		testutils.NoError(testutils.Require(t), err)
		testutils.Equal(testutils.Require(t), 1, len(paths))
		testutils.Equal(t, `^$CI_PROJECT_NAMESPACE\/.*$`, paths[0])
	})
}
