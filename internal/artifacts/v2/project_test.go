package artifacts_test

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	testlogrus "github.com/kilianpaquier/testlogrus/pkg"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xanzy/go-gitlab"

	"github.com/kilianpaquier/gitlab-storage-cleaner/internal/artifacts/mocks"
	"github.com/kilianpaquier/gitlab-storage-cleaner/internal/artifacts/v2"
	"github.com/kilianpaquier/gitlab-storage-cleaner/internal/artifacts/v2/tests"
)

func TestSplitProject(t *testing.T) {
	httpmock.Activate()
	t.Cleanup(httpmock.DeactivateAndReset)

	// setup mock client
	client, err := gitlab.NewClient("",
		gitlab.WithHTTPClient(&http.Client{Transport: httpmock.DefaultTransport}),
		gitlab.WithoutRetries(),
	)
	require.NoError(t, err)

	projectID := 5
	project := tests.NewProjectBuilder().
		SetID(projectID).
		Build()
	url := fmt.Sprintf("https://gitlab.com/api/v4/projects/%d/jobs", projectID)

	t.Run("error_list_jobs", func(t *testing.T) {
		// Arrange
		t.Cleanup(httpmock.Reset)
		httpmock.RegisterResponder(http.MethodGet, url,
			httpmock.NewStringResponder(http.StatusInternalServerError, "an error"))

		opts := tests.NewOptionsBuilder().Build()
		testlogrus.CatchLogs(t)

		// Act
		artifacts.SplitProject(client, *opts)(*project, nil)

		// Assert
		logs := testlogrus.Logs()
		assert.Contains(t, logs, "an error")
		assert.Contains(t, logs, "failed to retrieve project jobs")
	})

	t.Run("success_populate_channel", func(t *testing.T) {
		// Arrange
		after := time.Now().Add(time.Hour)

		mocks.MockPages(t, url, []*gitlab.Job{
			{
				ID:                7,
				ArtifactsExpireAt: &after,
				Artifacts: []struct {
					FileType   string `json:"file_type"`
					Filename   string `json:"filename"`
					Size       int    `json:"size"`
					FileFormat string `json:"file_format"`
				}{{}}, // at least one element to need cleanup
			},
			{
				ID:                8,
				ArtifactsExpireAt: &after,
				Artifacts: []struct {
					FileType   string `json:"file_type"`
					Filename   string `json:"filename"`
					Size       int    `json:"size"`
					FileFormat string `json:"file_format"`
				}{{}}, // at least one element to need cleanup
			},
		})

		opts := tests.NewOptionsBuilder().
			SetThresholdTime(time.Now()).
			Build()
		jobs := make(chan artifacts.Job, 10)
		t.Cleanup(func() { close(jobs) })

		// Act
		artifacts.SplitProject(client, *opts)(*project, jobs)

		// Assert
		assert.Equal(t, 2, httpmock.GetTotalCallCount())
		assert.Len(t, jobs, 2) // two elements, one for each job
	})
}

func TestMergeProject(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		jobs := make(chan artifacts.Job, 10)
		jobs <- *tests.NewJobBuilder().Build()
		jobs <- *tests.NewJobBuilder().
			SetCleaned(true).
			Build()
		jobs <- *tests.NewJobBuilder().
			SetCleaned(true).
			Build()
		close(jobs)

		p := artifacts.Project{}

		// Act
		p = p.Merge(jobs)

		// Assert
		assert.Equal(t, 2, p.JobsCleaned)
	})
}

func TestReadProjects(t *testing.T) {
	ctx := context.Background()

	httpmock.Activate()
	t.Cleanup(httpmock.DeactivateAndReset)

	// setup mock client
	client, err := gitlab.NewClient("",
		gitlab.WithHTTPClient(&http.Client{Transport: httpmock.DefaultTransport}),
		gitlab.WithoutRetries(),
	)
	require.NoError(t, err)

	url := "https://gitlab.com/api/v4/projects"

	opts := tests.NewOptionsBuilder().
		SetPathRegexps(
			regexp.MustCompile("^hey_.*$"),
			regexp.MustCompile("^hoï_.*$"),
		).
		Build()

	t.Run("error_list_projects", func(t *testing.T) {
		// Arrange
		t.Cleanup(httpmock.Reset)
		httpmock.RegisterResponder(http.MethodGet, url,
			httpmock.NewStringResponder(http.StatusInternalServerError, "an error"))
		testlogrus.CatchLogs(t)

		// Act
		projects := artifacts.ReadProjects(ctx, client, *opts)

		// Assert
		// verify channel first because it will block until its closed
		assert.Empty(t, lo.ChannelToSlice(projects))
		logs := testlogrus.Logs()
		assert.Contains(t, logs, "an error")
		assert.Contains(t, logs, "failed to retrieve projects")
	})

	t.Run("success_populate_channel", func(t *testing.T) {
		// Arrange
		mocks.MockPages(t, url, []*gitlab.Project{
			{ID: 7, PathWithNamespace: "hey_one"},
			{ID: 8, PathWithNamespace: "hey_two"},
			{ID: 9, PathWithNamespace: "two_hey"},
			{ID: 10, PathWithNamespace: "hoï_one"},
		})
		testlogrus.CatchLogs(t)

		// Act
		projects := artifacts.ReadProjects(ctx, client, *opts)

		// Assert
		// verify channel first because it will block until its closed
		assert.Len(t, lo.ChannelToSlice(projects), 3)
		logs := testlogrus.Logs()
		assert.Equal(t, logs, "")
	})
}
