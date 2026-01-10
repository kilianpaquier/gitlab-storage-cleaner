package artifacts_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/samber/lo"
	gitlab "gitlab.com/gitlab-org/api/client-go"

	"github.com/kilianpaquier/gitlab-storage-cleaner/internal/artifacts/engine"
	artifacts "github.com/kilianpaquier/gitlab-storage-cleaner/internal/artifacts/engine/v2"
	"github.com/kilianpaquier/gitlab-storage-cleaner/internal/testutils"
)

const (
	projectsURL  = "https://gitlab.com/api/v4/projects"
	jobsURL      = "https://gitlab.com/api/v4/projects/%d/jobs"
	artifactsURL = "https://gitlab.com/api/v4/projects/%d/jobs/%d/artifacts"
)

func TestRun(t *testing.T) {
	now := time.Now()
	ctx := t.Context()

	httpmock.Activate()
	t.Cleanup(httpmock.DeactivateAndReset)

	// setup mock client
	client, err := gitlab.NewClient("",
		gitlab.WithHTTPClient(&http.Client{Transport: httpmock.DefaultTransport}),
		gitlab.WithoutRetries(),
	)
	testutils.NoError(testutils.Require(t), err)

	var buf strings.Builder
	opts := []engine.RunOption{
		engine.WithLogger(engine.NewTestLogger(&buf)),
		engine.WithPaths("^project_path$"),
		engine.WithThresholdDuration(time.Hour),
	}

	t.Run("success_e2e", func(t *testing.T) {
		// Arrange
		t.Cleanup(httpmock.Reset)

		// projects endpoint mock
		projectID := int64(7)
		httpmock.RegisterResponder(http.MethodGet, projectsURL,
			httpmock.NewJsonResponderOrPanic(http.StatusOK, []*gitlab.Project{
				{ID: projectID, PathWithNamespace: "project_path"},
				{ID: 8, PathWithNamespace: "not_matching"},
			}).Then(httpmock.NewJsonResponderOrPanic(http.StatusOK, []*gitlab.Project{})))

		// jobs endpoint mock
		jobID := int64(10)
		httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf(jobsURL, projectID),
			httpmock.NewJsonResponderOrPanic(http.StatusOK, []*gitlab.Job{
				{
					ID:                jobID,
					Artifacts:         []gitlab.JobArtifact{{}},          // one artifact
					ArtifactsExpireAt: lo.ToPtr(now.Add(time.Hour)),      // artifacts not expired
					CreatedAt:         lo.ToPtr(now.Add(-2 * time.Hour)), // job is old
				},
				{
					ID:                18,
					Artifacts:         []gitlab.JobArtifact{{}},          // one artifact
					ArtifactsExpireAt: lo.ToPtr(now.Add(-time.Hour)),     // artifacts already expired
					CreatedAt:         lo.ToPtr(now.Add(-2 * time.Hour)), // job is old
				},
				{
					ID:        23,
					Artifacts: []gitlab.JobArtifact{}, // no artifacts
				},
			}).Then(httpmock.NewJsonResponderOrPanic(http.StatusOK, []*gitlab.Job{})))

		// job deletion endpoint mock
		httpmock.RegisterResponder(http.MethodDelete, fmt.Sprintf(artifactsURL, projectID, jobID),
			httpmock.NewStringResponder(http.StatusNoContent, ""))

		// expected calls to be made
		expectedCalls := map[string]int{
			"GET " + projectsURL: 2,
			fmt.Sprint("GET ", fmt.Sprintf(jobsURL, projectID)):     2,
			"DELETE " + fmt.Sprintf(artifactsURL, projectID, jobID): 1,
		}

		// Act
		err := artifacts.Run(ctx, client, opts...)

		// Assert
		testutils.NoError(testutils.Require(t), err)
		for k, v := range expectedCalls {
			actual, ok := httpmock.GetCallCountInfo()[k]
			testutils.True(t, ok)
			testutils.Equal(t, v, actual)
		}
		logs := buf.String()
		testutils.Contains(t, logs, "starting project execution")
		testutils.Contains(t, logs, "ending project execution")
		testutils.NotContains(t, logs, "failed to retrieve projects")
		testutils.NotContains(t, logs, "failed to retrieve project jobs")
		testutils.NotContains(t, logs, "failed to delete job's artifacts")
	})
}
