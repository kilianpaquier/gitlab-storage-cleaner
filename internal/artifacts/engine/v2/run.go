package artifacts

import (
	"context"
	"fmt"

	"github.com/fogfactory/pipe"
	"github.com/panjf2000/ants/v2"
	gitlab "gitlab.com/gitlab-org/api/client-go"

	"github.com/kilianpaquier/gitlab-storage-cleaner/internal/artifacts/engine"
)

// Run retrieves gitlab projects and filters the one not appropriate with options (paths regexps).
//
// For every appropriate project, it will retrieve jobs and delete outdated artifacts according to input option threshold.
func Run(parent context.Context, client *gitlab.Client, opts ...engine.RunOption) error {
	ro, err := engine.NewRunOptions(opts...)
	if err != nil {
		return fmt.Errorf("new run options: %w", err)
	}
	ctx := context.WithValue(parent, engine.LoggerKey, ro.Logger)

	pools, err := pipe.NewPoolsWithOptions([]int{10, 1000}, ants.WithLogger(engine.GetLogger(ctx)))
	if err != nil {
		return fmt.Errorf("pools initialization: %w", err)
	}
	defer pools.Release()

	piping := NewPipeProjectBuilder().
		Processor(StartProject(ctx)).
		Split(ReadJobs(ctx, client, ro)).
		Processor(DeleteArtifacts(ctx, client, ro)).
		Merge(ObserveCleanup).
		Processor(StopProject(ctx)).
		Build()

	projects := ReadProjects(ctx, client, ro)
	pipe.Run(pools, projects, piping)
	return nil
}
