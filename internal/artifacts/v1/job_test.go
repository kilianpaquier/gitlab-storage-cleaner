package artifacts_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	testlogrus "github.com/kilianpaquier/testlogrus/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xanzy/go-gitlab"

	"github.com/kilianpaquier/gitlab-storage-cleaner/internal/artifacts/v1"
	"github.com/kilianpaquier/gitlab-storage-cleaner/internal/artifacts/v1/tests"
)

func TestNeedCleanup(t *testing.T) {
	t.Run("false_no_artifacts", func(t *testing.T) {
		// Arrange
		job := tests.NewJobBuilder().Build()

		// Act
		clean := job.NeedCleanup(0, time.Time{})

		// Assert
		assert.False(t, clean)
	})

	t.Run("false_already_expired", func(t *testing.T) {
		// Arrange
		now := time.Now()
		job := tests.NewJobBuilder().
			SetArtifacts(artifacts.Artifact{}).
			SetArtifactsExpireAt(now.Add(-1 * time.Hour)).
			Build()

		// Act
		clean := job.NeedCleanup(0, now)

		// Assert
		assert.False(t, clean)
	})

	t.Run("false_not_above_size_threshold", func(t *testing.T) {
		// Arrange
		now := time.Now()
		job := tests.NewJobBuilder().
			SetArtifacts(artifacts.Artifact{}).
			SetArtifactsExpireAt(now.Add(time.Hour)).
			Build()

		// Act
		clean := job.NeedCleanup(10, now)

		// Assert
		assert.False(t, clean)
	})

	t.Run("false_too_recent", func(t *testing.T) {
		// Arrange
		now := time.Now().Add(time.Hour)
		job := tests.NewJobBuilder().
			SetArtifacts(artifacts.Artifact{}).
			SetArtifactsExpireAt(now.Add(-1 * time.Minute)).
			Build()

		// Act
		clean := job.NeedCleanup(0, now)

		// Assert
		assert.False(t, clean)
	})

	t.Run("success_true", func(t *testing.T) {
		// Arrange
		now := time.Now()
		job := tests.NewJobBuilder().
			SetArtifacts(artifacts.Artifact{Size: 1000}).
			SetArtifactsExpireAt(now.Add(time.Hour)).
			Build()

		// Act
		clean := job.NeedCleanup(100, now)

		// Assert
		assert.True(t, clean)
	})

	t.Run("success_true_exact_thresholds", func(t *testing.T) {
		// Arrange
		now := time.Now()
		job := tests.NewJobBuilder().
			SetArtifacts(artifacts.Artifact{Size: 1000}).
			SetArtifactsExpireAt(now.Add(time.Hour)).
			Build()

		// Act
		clean := job.NeedCleanup(1000, now)

		// Assert
		assert.True(t, clean)
	})

	t.Run("success_true_no_expiration", func(t *testing.T) {
		// Arrange
		now := time.Now()
		job := tests.NewJobBuilder().
			SetArtifacts(artifacts.Artifact{Size: 1000}).
			Build()

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

	job := tests.NewJobBuilder().
		SetID(jobID).
		SetProjectID(projectID).
		Build()

	t.Run("error_delete_call", func(t *testing.T) {
		// Arrange
		t.Cleanup(httpmock.Reset)
		httpmock.RegisterResponder(http.MethodDelete, url,
			httpmock.NewStringResponder(http.StatusInternalServerError, "an error"))
		testlogrus.CatchLogs(t)

		// Act
		job.DeleteArtifacts(ctx, client)(nil)

		// Assert
		logs := testlogrus.Logs()
		assert.Contains(t, logs, "an error")
		assert.Contains(t, logs, "failed to delete job's artifacts")
	})

	t.Run("success_deletion", func(t *testing.T) {
		// Arrange
		t.Cleanup(httpmock.Reset)
		httpmock.RegisterResponder(http.MethodDelete, url,
			httpmock.NewStringResponder(http.StatusNoContent, ""))
		testlogrus.CatchLogs(t)

		// Act
		job.DeleteArtifacts(ctx, client)(nil)

		// Assert
		logs := testlogrus.Logs()
		assert.Equal(t, logs, "")
		assert.Equal(t, 1, httpmock.GetTotalCallCount())
	})
}
