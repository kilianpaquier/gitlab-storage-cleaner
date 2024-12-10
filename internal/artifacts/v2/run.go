package artifacts

import (
	"context"
	"fmt"

	"github.com/fogfactory/pipe"
	"github.com/panjf2000/ants/v2"
	"github.com/sirupsen/logrus"
	"gitlab.com/gitlab-org/api/client-go"
)

// Run retrieves gitlab projects and filters the one not appropriate with options (paths regexps).
//
// For every appropriate project, it will retrieve jobs and delete outdated artifacts according to input option threshold.
func Run(ctx context.Context, client *gitlab.Client, opts Options) error {
	log := logrus.WithContext(ctx)

	sizes := []int{10, 1000}
	pools, err := pipe.NewPoolsWithOptions(sizes, ants.WithLogger(log))
	if err != nil {
		return fmt.Errorf("pools initialization: %w", err)
	}
	defer pools.Release()

	piping := NewPipeProjectBuilder().
		Processor((Project).Start).
		Split(SplitProject(client, opts)).
		Processor(DeleteArtifacts(client, opts)).
		Merge((Project).Merge).
		Processor((Project).Stop).
		Build()

	projects := ReadProjects(ctx, client, opts)
	pipe.Run(pools, projects, piping)
	return nil
}
