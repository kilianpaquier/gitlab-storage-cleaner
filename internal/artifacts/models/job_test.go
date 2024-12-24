package models_test

import (
	"context"
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
		clean := job.NeedCleanup(0, time.Time{})

		// Assert
		assert.False(t, clean)
	})

	t.Run("false_already_expired", func(t *testing.T) {
		// Arrange
		now := time.Now()
		job := models.Job{
			Artifacts:         []models.Artifact{{}},
			ArtifactsExpireAt: now.Add(-1 * time.Hour),
		}

		// Act
		clean := job.NeedCleanup(0, now)

		// Assert
		assert.False(t, clean)
	})

	t.Run("false_not_above_size_threshold", func(t *testing.T) {
		// Arrange
		now := time.Now()
		job := models.Job{
			Artifacts:         []models.Artifact{{}},
			ArtifactsExpireAt: now.Add(time.Hour),
		}

		// Act
		clean := job.NeedCleanup(10, now)

		// Assert
		assert.False(t, clean)
	})

	t.Run("false_too_recent", func(t *testing.T) {
		// Arrange
		now := time.Now().Add(time.Hour)
		job := models.Job{
			Artifacts:         []models.Artifact{{}},
			ArtifactsExpireAt: now.Add(-1 * time.Minute),
		}

		// Act
		clean := job.NeedCleanup(0, now)

		// Assert
		assert.False(t, clean)
	})

	t.Run("success_true", func(t *testing.T) {
		// Arrange
		now := time.Now()
		job := models.Job{
			Artifacts:         []models.Artifact{{Size: 1000}},
			ArtifactsExpireAt: now.Add(time.Hour),
		}

		// Act
		clean := job.NeedCleanup(100, now)

		// Assert
		assert.True(t, clean)
	})

	t.Run("success_true_exact_thresholds", func(t *testing.T) {
		// Arrange
		now := time.Now()
		job := models.Job{
			Artifacts:         []models.Artifact{{Size: 1000}},
			ArtifactsExpireAt: now.Add(time.Hour),
		}

		// Act
		clean := job.NeedCleanup(1000, now)

		// Assert
		assert.True(t, clean)
	})

	t.Run("success_true_no_expiration", func(t *testing.T) {
		// Arrange
		now := time.Now()
		job := models.Job{Artifacts: []models.Artifact{{Size: 1000}}}

		// Act
		clean := job.NeedCleanup(100, now)

		// Assert
		assert.True(t, clean)
	})
}

func TestDeleteArtifacts(t *testing.T) {
	ctx := context.Background()

	httpmock.Activate()
	t.Cleanup(httpmock.DeactivateAndReset)

	// setup mock client
	client, err := gitlab.NewClient("",
		gitlab.WithHTTPClient(&http.Client{Transport: httpmock.DefaultTransport}),
		gitlab.WithoutRetries(),
	)
	require.NoError(t, err)

	projectID := 5
	jobID := 5
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
		projectID := 5
		gitlab := gitlab.Job{
			ID:                1,
			ArtifactsExpireAt: lo.ToPtr(now),
			Artifacts: []struct {
				FileType   string `json:"file_type"`
				Filename   string `json:"filename"`
				Size       int    `json:"size"`
				FileFormat string `json:"file_format"`
			}{{Size: 1}},
		}
		expected := models.Job{
			ID:                1,
			ProjectID:         5,
			ArtifactsExpireAt: now,
			Artifacts:         []models.Artifact{{Size: 1}},
		}

		// Act
		project := models.JobFromGitLab(projectID, &gitlab)

		// Assert
		assert.Equal(t, expected, project)
	})
}
