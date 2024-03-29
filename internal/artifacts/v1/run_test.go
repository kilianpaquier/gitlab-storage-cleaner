package artifacts_test

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xanzy/go-gitlab"

	"github.com/kilianpaquier/gitlab-storage-cleaner/internal/artifacts/mocks"
	"github.com/kilianpaquier/gitlab-storage-cleaner/internal/artifacts/v1"
	"github.com/kilianpaquier/gitlab-storage-cleaner/internal/artifacts/v1/tests"
	"github.com/kilianpaquier/gitlab-storage-cleaner/internal/testlogs"
)

func TestRun(t *testing.T) {
	thresholdDuration := time.Hour
	start := time.Now().Add(thresholdDuration)
	ctx := context.Background()

	httpmock.Activate()
	t.Cleanup(httpmock.DeactivateAndReset)

	// setup mock client
	client, err := gitlab.NewClient("",
		gitlab.WithHTTPClient(&http.Client{Transport: httpmock.DefaultTransport}),
		gitlab.WithoutRetries(),
	)
	require.NoError(t, err)

	projectsURL := "https://gitlab.com/api/v4/projects"
	jobsURL := "https://gitlab.com/api/v4/projects/%d/jobs"
	deleteURL := "https://gitlab.com/api/v4/projects/%d/jobs/%d/artifacts"

	opts := tests.NewOptionsBuilder().
		SetPathRegexps(regexp.MustCompile("^project_path$")).
		SetThresholdTime(start).
		Build()

	t.Run("success_e2e", func(t *testing.T) {
		// Arrange
		t.Cleanup(httpmock.Reset)

		// projects endpoint mock
		projectID := 7
		mocks.MockPages(t, projectsURL, []*gitlab.Project{
			{ID: projectID, PathWithNamespace: "project_path"},
			{ID: 8, PathWithNamespace: "not_matching"},
		})

		// jobs endpoint mock
		jobID := 10
		mocks.MockPages(t, fmt.Sprintf(jobsURL, projectID), []*gitlab.Job{
			{
				ID: jobID,
				Artifacts: []struct {
					FileType   string `json:"file_type"`
					Filename   string `json:"filename"`
					Size       int    `json:"size"`
					FileFormat string `json:"file_format"`
				}{{}}, // one artifact
				ArtifactsExpireAt: lo.ToPtr(start.Add(time.Hour)), // date is after threshold
			},
			{
				ID: 18,
				Artifacts: []struct {
					FileType   string `json:"file_type"`
					Filename   string `json:"filename"`
					Size       int    `json:"size"`
					FileFormat string `json:"file_format"`
				}{{}}, // one artifact
				ArtifactsExpireAt: lo.ToPtr(start.Add(-time.Hour)), // date is before threshold
			},
			{
				ID: 23,
				Artifacts: []struct {
					FileType   string `json:"file_type"`
					Filename   string `json:"filename"`
					Size       int    `json:"size"`
					FileFormat string `json:"file_format"`
				}{}, // no artifacts
			},
		})

		// job deletion endpoint mock
		httpmock.RegisterResponder(http.MethodDelete, fmt.Sprintf(deleteURL, projectID, jobID),
			httpmock.NewStringResponder(http.StatusNoContent, ""))

		// expected calls to be made
		expectedCalls := map[string]int{
			fmt.Sprint("GET ", projectsURL, " <with_page_one>"):                     1,
			fmt.Sprint("GET ", projectsURL, " <with_page_two>"):                     1,
			fmt.Sprint("GET ", fmt.Sprintf(jobsURL, projectID), " <with_page_one>"): 1,
			fmt.Sprint("GET ", fmt.Sprintf(jobsURL, projectID), " <with_page_two>"): 1,
			"DELETE " + fmt.Sprintf(deleteURL, projectID, jobID):                    1,
		}

		hook := test.NewGlobal()
		t.Cleanup(func() { hook.Reset() })

		// Act
		err := artifacts.Run(ctx, client, *opts)

		// Assert
		assert.NoError(t, err)
		logs := testlogs.ToString(hook.AllEntries())
		assert.NotContains(t, logs, "failed to retrieve projects")
		assert.Contains(t, logs, "running project cleanup")
		assert.NotContains(t, logs, "failed to retrieve project jobs")
		assert.NotContains(t, logs, "failed to delete job's artifacts")
		assert.Contains(t, logs, "ended project cleanup")
		assert.Equal(t, expectedCalls, httpmock.GetCallCountInfo())
	})
}
