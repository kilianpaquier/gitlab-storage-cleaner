package artifacts

import (
	"context"
	"fmt"

	pooling "github.com/kilianpaquier/pooling/pkg"
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
	ctx := ro.Context(parent)

	pooler, err := pooling.NewPoolerBuilder().
		// initialize two pools of routines, one for projects and one for jobs
		SetSizes(10, 1000).
		SetOptions(ants.WithLogger(engine.GetLogger(ctx))).
		Build()
	if err != nil {
		return fmt.Errorf("pooler initialization: %w", err)
	}
	defer pooler.Close()

	pooler.Read(ReadProjects(ctx, client, ro))
	return nil
}
