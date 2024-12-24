package models_test

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	gitlab "gitlab.com/gitlab-org/api/client-go"

	"github.com/kilianpaquier/gitlab-storage-cleaner/internal/artifacts/models"
)

func TestMatches(t *testing.T) {
	t.Run("does_not_match", func(t *testing.T) {
		// Arrange
		project := models.Project{PathWithNamespace: "john.doe"}
		regexps := []*regexp.Regexp{
			regexp.MustCompile("^hoï_.*$"),
			regexp.MustCompile("^hey_.*$"),
		}

		// Act
		matches := project.Matches(regexps...)

		// Assert
		assert.False(t, matches)
	})

	t.Run("matches", func(t *testing.T) {
		// Arrange
		project := models.Project{PathWithNamespace: "hey_john.doe"}
		regexps := []*regexp.Regexp{
			regexp.MustCompile("^hoï_.*$"),
			regexp.MustCompile("^hey_.*$"),
		}

		// Act
		matches := project.Matches(regexps...)

		// Assert
		assert.True(t, matches)
	})
}

func TestProjectFromGitLab(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		gitlab := gitlab.Project{ID: 1, PathWithNamespace: "john.doe"}
		expected := models.Project{ID: 1, PathWithNamespace: "john.doe"}

		// Act
		project := models.ProjectFromGitLab(&gitlab)

		// Assert
		assert.Equal(t, expected, project)
	})
}
