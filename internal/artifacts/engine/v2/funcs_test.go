package artifacts_test

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gitlab "gitlab.com/gitlab-org/api/client-go"

	"github.com/kilianpaquier/gitlab-storage-cleaner/internal/artifacts/engine"
	artifacts "github.com/kilianpaquier/gitlab-storage-cleaner/internal/artifacts/engine/v2"
	"github.com/kilianpaquier/gitlab-storage-cleaner/internal/artifacts/models"
)

func TestReadProjects(t *testing.T) {
	ctx := context.Background()

	httpmock.Activate()
	t.Cleanup(httpmock.DeactivateAndReset)

	// setup mock client
	client, err := gitlab.NewClient("",
		gitlab.WithHTTPClient(&http.Client{Transport: httpmock.DefaultTransport}),
		gitlab.WithoutRetries())
	require.NoError(t, err)

	runOptions, err := engine.NewRunOptions(engine.WithPaths("^hey_.*$", "^hoï_.*$"))
	require.NoError(t, err)

	t.Run("error_list_projects", func(t *testing.T) {
		// Arrange
		t.Cleanup(httpmock.Reset)
		httpmock.RegisterResponder(http.MethodGet, projectsURL,
			httpmock.NewStringResponder(http.StatusInternalServerError, "an error"))

		var buf strings.Builder
		ctx := context.WithValue(ctx, engine.LoggerKey, engine.NewTestLogger(&buf))

		// Act
		projects := artifacts.ReadProjects(ctx, client, runOptions)

		// Assert
		// verify channel first because it will block until its closed
		assert.Empty(t, lo.ChannelToSlice(projects))
		logs := buf.String()
		assert.Contains(t, logs, "an error")
		assert.Contains(t, logs, "failed to retrieve projects")
	})

	t.Run("success_populate_channel", func(t *testing.T) {
		// Arrange
		t.Cleanup(httpmock.Reset)
		httpmock.RegisterResponder(http.MethodGet, projectsURL,
			httpmock.NewJsonResponderOrPanic(http.StatusOK, []*gitlab.Project{
				{ID: 7, PathWithNamespace: "hey_one"},
				{ID: 8, PathWithNamespace: "hey_two"},
				{ID: 9, PathWithNamespace: "two_hey"},
				{ID: 10, PathWithNamespace: "hoï_one"},
			}).Then(httpmock.NewJsonResponderOrPanic(http.StatusOK, []*gitlab.Project{})))

		var buf strings.Builder
		ctx := context.WithValue(ctx, engine.LoggerKey, engine.NewTestLogger(&buf))

		// Act
		projects := artifacts.ReadProjects(ctx, client, runOptions)

		// Assert
		// verify channel first because it will block until its closed
		assert.Len(t, lo.ChannelToSlice(projects), 3)
		assert.Contains(t, buf.String(), "skipping project cleaning project_id=9 project_path=two_hey")
	})
}

func TestReadJobs(t *testing.T) {
	ctx := context.Background()

	httpmock.Activate()
	t.Cleanup(httpmock.DeactivateAndReset)

	// setup mock client
	client, err := gitlab.NewClient("",
		gitlab.WithHTTPClient(&http.Client{Transport: httpmock.DefaultTransport}),
		gitlab.WithoutRetries())
	require.NoError(t, err)

	project := artifacts.Project{Project: models.Project{ID: 5}}

	t.Run("error_list_jobs", func(t *testing.T) {
		// Arrange
		t.Cleanup(httpmock.Reset)
		httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf(jobsURL, project.ID),
			httpmock.NewStringResponder(http.StatusInternalServerError, "an error"))

		var buf strings.Builder
		ctx := context.WithValue(ctx, engine.LoggerKey, engine.NewTestLogger(&buf))

		// Act
		artifacts.ReadJobs(ctx, client, engine.RunOptions{})(project, nil)

		// Assert
		logs := buf.String()
		assert.Contains(t, logs, "an error")
		assert.Contains(t, logs, "failed to retrieve project jobs")
	})

	t.Run("success_populate_channel", func(t *testing.T) {
		// Arrange
		now := time.Now()
		start := now.Add(-2 * time.Hour) // jobs are old and artifacts not cleaned yet

		t.Cleanup(httpmock.Reset)
		httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf(jobsURL, project.ID),
			httpmock.NewJsonResponderOrPanic(http.StatusOK, []*gitlab.Job{
				{
					ID:                7,
					ArtifactsExpireAt: lo.ToPtr(now.Add(time.Hour)),
					CreatedAt:         &start,
					Artifacts: []struct {
						FileType   string `json:"file_type"`
						Filename   string `json:"filename"`
						Size       int    `json:"size"`
						FileFormat string `json:"file_format"`
					}{{}}, // at least one element to need cleanup
				},
				{
					ID:                8,
					ArtifactsExpireAt: lo.ToPtr(now.Add(time.Hour)),
					CreatedAt:         &start,
					Artifacts: []struct {
						FileType   string `json:"file_type"`
						Filename   string `json:"filename"`
						Size       int    `json:"size"`
						FileFormat string `json:"file_format"`
					}{{}}, // at least one element to need cleanup
				},
			}).Then(httpmock.NewJsonResponderOrPanic(http.StatusOK, []*gitlab.Job{})))

		ro, _ := engine.NewRunOptions(engine.WithThresholdDuration(time.Second))

		jobs := make(chan models.Job, 10)
		t.Cleanup(func() { close(jobs) })

		// Act
		artifacts.ReadJobs(ctx, client, ro)(project, jobs)

		// Assert
		assert.Len(t, jobs, 2) // two elements, one for each job
		assert.Equal(t, 2, httpmock.GetTotalCallCount())
	})
}

func TestDeleteArtifacts(t *testing.T) {
	ctx := context.Background()

	httpmock.Activate()
	t.Cleanup(httpmock.DeactivateAndReset)

	// setup mock client
	client, err := gitlab.NewClient("",
		gitlab.WithHTTPClient(&http.Client{Transport: httpmock.DefaultTransport}),
		gitlab.WithoutRetries())
	require.NoError(t, err)

	job := models.Job{ID: 7, ProjectID: 5}

	t.Run("success_dry_run", func(t *testing.T) {
		// Arrange
		var buf strings.Builder
		ctx := context.WithValue(ctx, engine.LoggerKey, engine.NewTestLogger(&buf))

		// Act
		job := artifacts.DeleteArtifacts(ctx, client, engine.RunOptions{DryRun: true})(job)

		// Assert
		assert.False(t, job.Cleaned)
		assert.Contains(t, buf.String(), "running in dry run mode, skipping job's artifacts deletion")
	})

	t.Run("error_delete_artifacts", func(t *testing.T) {
		// Arrange
		t.Cleanup(httpmock.Reset)
		httpmock.RegisterResponder(http.MethodDelete, fmt.Sprintf(artifactsURL, job.ProjectID, job.ID),
			httpmock.NewStringResponder(http.StatusInternalServerError, "an error"))

		var buf strings.Builder
		ctx := context.WithValue(ctx, engine.LoggerKey, engine.NewTestLogger(&buf))

		// Act
		job := artifacts.DeleteArtifacts(ctx, client, engine.RunOptions{})(job)

		// Assert
		assert.False(t, job.Cleaned)
		logs := buf.String()
		assert.Contains(t, logs, "an error")
		assert.Contains(t, logs, "failed to delete job's artifacts")
	})

	t.Run("success_delete_artifacts", func(t *testing.T) {
		// Arrange
		t.Cleanup(httpmock.Reset)
		httpmock.RegisterResponder(http.MethodDelete, fmt.Sprintf(artifactsURL, job.ProjectID, job.ID),
			httpmock.NewStringResponder(http.StatusNoContent, ""))

		var buf strings.Builder
		ctx := context.WithValue(ctx, engine.LoggerKey, engine.NewTestLogger(&buf))

		// Act
		job := artifacts.DeleteArtifacts(ctx, client, engine.RunOptions{})(job)

		// Assert
		assert.Empty(t, buf.String())
		assert.True(t, job.Cleaned)
	})
}

func TestObserveCleanup(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		jobs := make(chan models.Job, 10)
		jobs <- models.Job{}
		jobs <- models.Job{Cleaned: true}
		jobs <- models.Job{Cleaned: true}
		close(jobs)

		project := artifacts.Project{}

		// Act
		project = artifacts.ObserveCleanup(project, jobs)

		// Assert
		assert.Equal(t, 2, project.JobsCleaned)
	})
}
