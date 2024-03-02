package artifacts_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/kilianpaquier/gitlab-storage-cleaner/internal/artifacts/v1"
)

func TestEnsureDefaults(t *testing.T) {
	t.Run("error_no_options", func(t *testing.T) {
		// Arrange
		opts := artifacts.Options{}

		// Act
		err := opts.EnsureDefaults()

		// Assert
		assert.ErrorContains(t, err, "'Paths' failed on the 'required' tag")
		assert.ErrorContains(t, err, "'ThresholdDuration' failed on the 'required' tag")
		assert.ErrorContains(t, err, "'ThresholdSize' failed on the 'required' tag")
	})

	t.Run("error_invalid_regexp", func(t *testing.T) {
		// Arrange
		opts := artifacts.Options{
			Paths:             []string{"\\"},
			ThresholdDuration: time.Second,
		}

		// Act
		err := opts.EnsureDefaults()

		// Assert
		assert.ErrorContains(t, err, "invalid regexp '\\'")
	})

	t.Run("error_invalid_duration", func(t *testing.T) {
		// Arrange
		opts := artifacts.Options{
			Paths:             []string{".*"},
			ThresholdDuration: -time.Second,
		}

		// Act
		err := opts.EnsureDefaults()

		// Assert
		assert.ErrorContains(t, err, "'ThresholdDuration' failed on the 'gt' tag")
	})

	t.Run("success", func(t *testing.T) {
		// Arrange
		opts := artifacts.Options{
			Paths:             []string{".*", ".+"},
			ThresholdDuration: time.Second,
			ThresholdSize:     uint64(1),
		}

		// Act
		err := opts.EnsureDefaults()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, []*regexp.Regexp{regexp.MustCompile(".*"), regexp.MustCompile(".+")}, opts.PathRegexps)
		assert.WithinDuration(t, time.Now().Add(time.Second), opts.ThresholdTime, 10*time.Millisecond)
	})
}
