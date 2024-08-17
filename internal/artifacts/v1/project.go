package artifacts

import (
	"context"
	"regexp"

	pooling "github.com/kilianpaquier/pooling/pkg"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"
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
			for _, project := range projects {
				p := Project{
					ID:                project.ID,
					PathWithNamespace: project.PathWithNamespace,
				}

				// confirm that project is inside cleanup slice
				_, found := lo.Find(opts.PathRegexps, func(reg *regexp.Regexp) bool {
					return reg.MatchString(p.PathWithNamespace)
				})
				if found {
					tasks <- p.CleanArtifacts(ctx, client, opts)
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

			for _, job := range jobs {
				j := Job{
					Artifacts: func() []Artifact {
						artifacts := make([]Artifact, 0, len(job.Artifacts))
						for _, artifact := range job.Artifacts {
							artifacts = append(artifacts, Artifact{Size: uint64(artifact.Size)})
						}
						return artifacts
					}(),
					ArtifactsExpireAt: lo.FromPtr(job.ArtifactsExpireAt),
					ID:                job.ID,
					ProjectID:         p.ID,
				}

				// check that the job needs to be cleaned up
				if !opts.DryRun && j.NeedCleanup(opts.ThresholdSize, opts.ThresholdTime) {
					funcs <- j.DeleteArtifacts(ctx, client)
				}
			}
		}

		log.Info("ended project cleanup")
	}
}
