package engine_test

import (
	"testing"

	"github.com/kilianpaquier/gitlab-storage-cleaner/internal/artifacts/engine"
	"github.com/kilianpaquier/gitlab-storage-cleaner/internal/testutils"
)

func TestRunOptions(t *testing.T) {
	t.Run("error_invalid_regexps", func(t *testing.T) {
		// Arrange
		opts := []engine.RunOption{engine.WithPaths(`/\/\`)}

		// Act
		_, err := engine.NewRunOptions(opts...)

		// Assert
		testutils.Error(testutils.Require(t), err)
		testutils.Contains(t, err.Error(), `invalid regexp '/\/\'`)
	})

	t.Run("success_defaults", func(t *testing.T) {
		// Arrange
		opts := []engine.RunOption{}

		// Act
		runOptions, err := engine.NewRunOptions(opts...)

		// Assert
		testutils.NoError(testutils.Require(t), err)
		testutils.Equal(t, engine.DefaultThresholdDuration, runOptions.ThresholdDuration)
		testutils.NotNil(t, engine.GetLogger(runOptions.Context(t.Context())))
	})
}
