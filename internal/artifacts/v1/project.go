package artifacts

import (
	"context"
	"regexp"

	pooling "github.com/kilianpaquier/pooling/pkg"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

// Project is a simplified view of a gitlab project with only useful information used during artifacts command.
type Project struct {
	ID                int
	PathWithNamespace string
}

// ReadProjects reads all projects from gitlab api and send them into the output channel.
// The output channel is closed once all projects were sent into it.
func ReadProjects(ctx context.Context, client *gitlab.Client, opts Options) <-chan pooling.PoolerFunc {
	// un-buffered channel to avoid too many pages in memory
	tasks := make(chan pooling.PoolerFunc)

	projectsOpts := &gitlab.ListProjectsOptions{
		ListOptions: gitlab.ListOptions{
			Page:    1,
			PerPage: 100,
		},

		Archived:             lo.ToPtr(false),
		IncludePendingDelete: lo.ToPtr(false),
		Membership:           lo.ToPtr(true),
		// only maintainers can cleanup job artifacts
		MinAccessLevel: lo.ToPtr(gitlab.MaintainerPermissions),
		Simple:         lo.ToPtr(true),
	}

	go func() {
		defer close(tasks)
		for {
			// retrieve next page of projects
			projects, _, err := client.Projects.ListProjects(projectsOpts, gitlab.WithContext(ctx))
			if err != nil {
				logrus.WithContext(ctx).
					WithError(err).
					Warn("failed to retrieve projects")
				break
			}
			projectsOpts.Page++

			// stop infinite loop
			if len(projects) == 0 {
				break
			}

			// send all projects for cleanup and iterate to next page
			for _, p := range projects {
				project := Project{
					ID:                p.ID,
					PathWithNamespace: p.PathWithNamespace,
				}

				// confirm that project is inside cleanup slice
				_, found := lo.Find(opts.PathRegexps, func(reg *regexp.Regexp) bool {
					return reg.MatchString(project.PathWithNamespace)
				})
				if found {
					tasks <- project.CleanArtifacts(ctx, client, opts)
				}
			}
		}
	}()

	return tasks
}

// CleanArtifacts returns the function to clean artifacts a specific project.
//
// This function retrieves all project's jobs and send them into pooling PoolerFunc input channel.
func (p Project) CleanArtifacts(ctx context.Context, client *gitlab.Client, opts Options) pooling.PoolerFunc {
	log := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"project_id":   p.ID,
		"project_path": p.PathWithNamespace,
	})

	return func(funcs chan<- pooling.PoolerFunc) {
		// handle dry run option
		if opts.DryRun {
			log.Info("should run project cleanup")
		} else {
			log.Info("running project cleanup")
		}

		jobsOpts := &gitlab.ListJobsOptions{
			ListOptions: gitlab.ListOptions{
				Page:    1,
				PerPage: 100,
			},
		}

		for {
			jobs, _, err := client.Jobs.ListProjectJobs(p.ID, jobsOpts, gitlab.WithContext(ctx))
			if err != nil {
				log.WithError(err).Warn("failed to retrieve project jobs")
				break
			}

			// stop infinite loop
			if len(jobs) == 0 {
				break
			}
			jobsOpts.Page++

			for _, j := range jobs {
				job := Job{
					Artifacts: func() []Artifact {
						artifacts := make([]Artifact, 0, len(j.Artifacts))
						for _, artifact := range j.Artifacts {
							artifacts = append(artifacts, Artifact{Size: artifact.Size})
						}
						return artifacts
					}(),
					ArtifactsExpireAt: lo.FromPtr(j.ArtifactsExpireAt),
					ID:                j.ID,
					ProjectID:         p.ID,
				}

				// check that the job needs to be cleaned up
				if !opts.DryRun && job.NeedCleanup(opts.ThresholdSize, opts.ThresholdTime) {
					funcs <- job.DeleteArtifacts(ctx, client)
				}
			}
		}

		log.Info("ended project cleanup")
	}
}
