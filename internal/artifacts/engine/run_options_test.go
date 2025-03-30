package engine_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kilianpaquier/gitlab-storage-cleaner/internal/artifacts/engine"
)

func TestRunOptions(t *testing.T) {
	t.Run("error_invalid_regexps", func(t *testing.T) {
		// Arrange
		opts := []engine.RunOption{engine.WithPaths(`/\/\`)}

		// Act
		_, err := engine.NewRunOptions(opts...)

		// Assert
		assert.ErrorContains(t, err, `invalid regexp '/\/\'`)
	})

	t.Run("success_defaults", func(t *testing.T) {
		// Arrange
		opts := []engine.RunOption{}

		// Act
		runOptions, err := engine.NewRunOptions(opts...)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, engine.DefaultThresholdDuration, runOptions.ThresholdDuration)
		assert.NotNil(t, engine.GetLogger(runOptions.Context(t.Context())))
	})
}
