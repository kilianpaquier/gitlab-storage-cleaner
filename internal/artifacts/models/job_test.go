package models_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gitlab "gitlab.com/gitlab-org/api/client-go"

	"github.com/kilianpaquier/gitlab-storage-cleaner/internal/artifacts/models"
)

func TestNeedCleanup(t *testing.T) {
	t.Run("false_no_artifacts", func(t *testing.T) {
		// Arrange
		job := models.Job{}

		// Act
		clean := job.NeedCleanup(0)

		// Assert
		assert.False(t, clean)
	})

	t.Run("false_already_expired", func(t *testing.T) {
		// Arrange
		now := time.Now()
		job := models.Job{
			ArtifactsCount:    1,
			ArtifactsExpireAt: now.Add(-1 * time.Hour),
			CreatedAt:         now.Add(-1 * time.Hour),
		}

		// Act
		clean := job.NeedCleanup(0)

		// Assert
		assert.False(t, clean)
	})

	t.Run("false_too_recent", func(t *testing.T) {
		// Arrange
		now := time.Now().Add(time.Hour)
		job := models.Job{
			ArtifactsCount:    1,
			ArtifactsExpireAt: now.Add(-1 * time.Minute),
			CreatedAt:         now.Add(-5 * time.Minute),
		}

		// Act
		clean := job.NeedCleanup(time.Hour)

		// Assert
		assert.False(t, clean)
	})

	t.Run("success_true_no_creation_date", func(t *testing.T) {
		// Arrange
		job := models.Job{ArtifactsCount: 1}

		// Act
		clean := job.NeedCleanup(0)

		// Assert
		assert.True(t, clean)
	})
}

func TestDeleteArtifacts(t *testing.T) {
	ctx := t.Context()

	httpmock.Activate()
	t.Cleanup(httpmock.DeactivateAndReset)

	// setup mock client
	client, err := gitlab.NewClient("",
		gitlab.WithHTTPClient(&http.Client{Transport: httpmock.DefaultTransport}),
		gitlab.WithoutRetries(),
	)
	require.NoError(t, err)

	projectID := int64(5)
	jobID := int64(5)
	url := fmt.Sprintf("https://gitlab.com/api/v4/projects/%d/jobs/%d/artifacts", projectID, jobID)

	job := models.Job{ID: jobID, ProjectID: projectID}

	t.Run("error_delete_call", func(t *testing.T) {
		// Arrange
		t.Cleanup(httpmock.Reset)
		httpmock.RegisterResponder(http.MethodDelete, url,
			httpmock.NewStringResponder(http.StatusInternalServerError, "an error"))

		// Act
		err := job.DeleteArtifacts(ctx, client)

		// Assert
		assert.ErrorContains(t, err, "delete artifacts")
	})

	t.Run("success_deletion", func(t *testing.T) {
		// Arrange
		t.Cleanup(httpmock.Reset)
		httpmock.RegisterResponder(http.MethodDelete, url,
			httpmock.NewStringResponder(http.StatusNoContent, ""))

		// Act
		err := job.DeleteArtifacts(ctx, client)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, 1, httpmock.GetTotalCallCount())
	})
}

func TestJobFromGitLab(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		now := time.Now()
		projectID := int64(5)
		gitlab := gitlab.Job{
			ID:                1,
			CreatedAt:         lo.ToPtr(now),
			ArtifactsExpireAt: lo.ToPtr(now.Add(time.Hour)),
			Artifacts:         []gitlab.JobArtifact{{}},
		}
		expected := models.Job{
			ArtifactsCount:    1,
			ArtifactsExpireAt: now.Add(time.Hour),
			CreatedAt:         now,
			ID:                1,
			ProjectID:         5,
		}

		// Act
		project := models.JobFromGitLab(projectID, &gitlab)

		// Assert
		assert.Equal(t, expected, project)
	})
}
