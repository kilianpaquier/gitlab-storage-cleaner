package artifacts_test

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	pooling "github.com/kilianpaquier/pooling/pkg"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/api/client-go"

	"github.com/kilianpaquier/gitlab-storage-cleaner/internal/artifacts/v1"
)

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

	opts := artifacts.Options{
		PathRegexps: []*regexp.Regexp{
			regexp.MustCompile("^hey_.*$"),
			regexp.MustCompile("^hoï_.*$"),
		},
	}

	t.Run("error_list_projects", func(t *testing.T) {
		// Arrange
		t.Cleanup(httpmock.Reset)
		httpmock.RegisterResponder(http.MethodGet, url,
			httpmock.NewStringResponder(http.StatusInternalServerError, "an error"))

		hook := test.NewGlobal()
		t.Cleanup(func() { hook.Reset() })

		// Act
		projects := artifacts.ReadProjects(ctx, client, opts)

		// Assert
		// verify channel first because it will block until its closed
		assert.Empty(t, lo.ChannelToSlice(projects))
		logs := toString(hook.AllEntries())
		assert.Contains(t, logs, "an error")
		assert.Contains(t, logs, "failed to retrieve projects")
	})

	t.Run("success_populate_channel", func(t *testing.T) {
		// Arrange
		t.Cleanup(httpmock.Reset)
		httpmock.RegisterResponder(http.MethodGet, url,
			httpmock.NewJsonResponderOrPanic(http.StatusOK, []*gitlab.Project{
				{ID: 7, PathWithNamespace: "hey_one"},
				{ID: 8, PathWithNamespace: "hey_two"},
				{ID: 9, PathWithNamespace: "two_hey"},
				{ID: 10, PathWithNamespace: "hoï_one"},
			}).Then(httpmock.NewJsonResponderOrPanic(http.StatusOK, []*gitlab.Project{})))

		hook := test.NewGlobal()
		t.Cleanup(func() { hook.Reset() })

		// Act
		projects := artifacts.ReadProjects(ctx, client, opts)

		// Assert
		// verify channel first because it will block until its closed
		assert.Len(t, lo.ChannelToSlice(projects), 3)
		logs := toString(hook.AllEntries())
		assert.Equal(t, "", logs)
	})
}

func TestCleanArtifacts(t *testing.T) {
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
	project := artifacts.Project{ID: projectID}
	url := fmt.Sprintf("https://gitlab.com/api/v4/projects/%d/jobs", projectID)

	t.Run("error_list_jobs", func(t *testing.T) {
		// Arrange
		t.Cleanup(httpmock.Reset)
		httpmock.RegisterResponder(http.MethodGet, url,
			httpmock.NewStringResponder(http.StatusInternalServerError, "an error"))

		hook := test.NewGlobal()
		t.Cleanup(func() { hook.Reset() })

		// Act
		project.CleanArtifacts(ctx, client, artifacts.Options{})(nil)

		// Assert
		logs := toString(hook.AllEntries())
		assert.Contains(t, logs, "an error")
		assert.Contains(t, logs, "failed to retrieve project jobs")
	})

	t.Run("success_dry_run_populate_channel", func(t *testing.T) {
		// Arrange
		t.Cleanup(httpmock.Reset)
		httpmock.RegisterResponder(http.MethodGet, url,
			httpmock.NewJsonResponderOrPanic(http.StatusOK, []*gitlab.Job{{ID: 7}, {ID: 8}}).
				Then(httpmock.NewJsonResponderOrPanic(http.StatusOK, []*gitlab.Job{})))

		opts := artifacts.Options{DryRun: true}

		jobs := make(chan pooling.PoolerFunc, 10)
		t.Cleanup(func() { close(jobs) })

		hook := test.NewGlobal()
		t.Cleanup(func() { hook.Reset() })

		// Act
		project.CleanArtifacts(ctx, client, opts)(jobs)

		// Assert
		assert.Empty(t, jobs) // no elements since dry run doesn't send into channel
		assert.Equal(t, 2, httpmock.GetTotalCallCount())
		logs := toString(hook.AllEntries())
		assert.Contains(t, logs, "should run project cleanup")
	})

	t.Run("success_populate_channel", func(t *testing.T) {
		// Arrange
		after := time.Now().Add(time.Hour)
		t.Cleanup(httpmock.Reset)
		httpmock.RegisterResponder(http.MethodGet, url,
			httpmock.NewJsonResponderOrPanic(http.StatusOK, []*gitlab.Job{
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
			}).Then(httpmock.NewJsonResponderOrPanic(http.StatusOK, []*gitlab.Job{})))

		opts := artifacts.Options{ThresholdTime: time.Now()}

		jobs := make(chan pooling.PoolerFunc, 10)
		t.Cleanup(func() { close(jobs) })

		hook := test.NewGlobal()
		t.Cleanup(func() { hook.Reset() })

		// Act
		project.CleanArtifacts(ctx, client, opts)(jobs)

		// Assert
		assert.Len(t, jobs, 2) // two elements, one for each job
		assert.Equal(t, 2, httpmock.GetTotalCallCount())
		logs := toString(hook.AllEntries())
		assert.Contains(t, logs, "running project cleanup")
		assert.Contains(t, logs, "ended project cleanup")
	})
}
