package artifacts_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xanzy/go-gitlab"

	"github.com/kilianpaquier/gitlab-storage-cleaner/internal/artifacts/v2"
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

	opts := artifacts.Options{
		PathRegexps:   []*regexp.Regexp{regexp.MustCompile("^project_path$")},
		ThresholdTime: start,
	}

	t.Run("success_e2e", func(t *testing.T) {
		// Arrange
		t.Cleanup(httpmock.Reset)

		// projects endpoint mock
		projectID := 7
		httpmock.RegisterResponder(http.MethodGet, projectsURL,
			httpmock.NewJsonResponderOrPanic(http.StatusOK, []*gitlab.Project{
				{ID: projectID, PathWithNamespace: "project_path"},
				{ID: 8, PathWithNamespace: "not_matching"},
			}).Then(httpmock.NewJsonResponderOrPanic(http.StatusOK, []*gitlab.Project{})))

		// jobs endpoint mock
		jobID := 10
		httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf(jobsURL, projectID),
			httpmock.NewJsonResponderOrPanic(http.StatusOK, []*gitlab.Job{
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
			}).Then(httpmock.NewJsonResponderOrPanic(http.StatusOK, []*gitlab.Job{})))

		// job deletion endpoint mock
		httpmock.RegisterResponder(http.MethodDelete, fmt.Sprintf(deleteURL, projectID, jobID),
			httpmock.NewStringResponder(http.StatusNoContent, ""))

		// expected calls to be made
		expectedCalls := map[string]int{
			"GET " + projectsURL: 2,
			fmt.Sprint("GET ", fmt.Sprintf(jobsURL, projectID)):  2,
			"DELETE " + fmt.Sprintf(deleteURL, projectID, jobID): 1,
		}

		hook := test.NewGlobal()
		t.Cleanup(func() { hook.Reset() })

		// Act
		err := artifacts.Run(ctx, client, opts)

		// Assert
		require.NoError(t, err)
		logs := toString(hook.AllEntries())
		assert.NotContains(t, logs, "failed to retrieve projects")
		assert.Contains(t, logs, "starting project execution")
		assert.NotContains(t, logs, "failed to retrieve project jobs")
		assert.NotContains(t, logs, "failed to delete job's artifacts")
		assert.Contains(t, logs, "ending project execution")
		assert.Equal(t, expectedCalls, httpmock.GetCallCountInfo())
	})
}

// toString transforms a slice of logrus entry into a string concatenation.
func toString(entries []*logrus.Entry) string {
	var buf bytes.Buffer
	for _, entry := range entries {
		b, _ := entry.Bytes()
		buf.Write(b)
		buf.WriteString("\n")
	}
	return buf.String()
}
