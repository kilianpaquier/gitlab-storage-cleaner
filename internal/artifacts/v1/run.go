package artifacts

import (
	"context"
	"fmt"

	pooling "github.com/kilianpaquier/pooling/pkg"
	"github.com/panjf2000/ants/v2"
	"github.com/sirupsen/logrus"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

// Run retrieves gitlab projects and filters the one not appropriate with options (paths regexps).
//
// For every appropriate project, it will retrieve jobs and delete outdated artifacts according to input option threshold.
func Run(ctx context.Context, client *gitlab.Client, opts Options) error {
	log := logrus.WithContext(ctx)

	pooler, err := pooling.NewPoolerBuilder().
		// initialize two pools of routines, one for projects and one for jobs
		SetSizes(10, 500).
		SetOptions(ants.WithLogger(log)).
		Build()
	if err != nil {
		return fmt.Errorf("pooler initialization: %w", err)
	}
	defer pooler.Close()

	pooler.Read(ReadProjects(ctx, client, opts))
	return nil
}
