package artifacts

import (
	"context"
	"regexp"
	"time"

	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"
)

// Project is a simplified view of a gitlab project with only useful information used during artifacts command.
type Project struct {
	ID                int
	JobsCleaned       int
	PathWithNamespace string

	executionDuration time.Duration
	executionStart    time.Time
}

// Start starts the timer for Project execution and logs the Project start execution.
func (p Project) Start() Project {
	p.executionStart = time.Now() // nolint:revive
	logrus.WithFields(logrus.Fields{
		"project_id":   p.ID,
		"project_path": p.PathWithNamespace,
	}).Info("starting project execution")
	return p
}

// Stop stops the project timer execution and logs the Project execution result.
func (p Project) Stop() Project {
	p.executionDuration = time.Since(p.executionStart) // nolint:revive
	logrus.WithFields(logrus.Fields{
		"execution_duration": p.executionDuration,
		"jobs_cleaned":       p.JobsCleaned,
		"project_id":         p.ID,
		"project_path":       p.PathWithNamespace,
	}).Info("ending project execution")
	return p
}

// SplitProject returns the function to send all Jobs of a given Project into pipe processing.
func SplitProject(client *gitlab.Client, opts Options) func(p Project, jobs chan<- Job) {
	return func(p Project, in chan<- Job) {
		log := logrus.WithFields(logrus.Fields{
			"project_id":   p.ID,
			"project_path": p.PathWithNamespace,
		})

		jobsOpts := &gitlab.ListJobsOptions{
			ListOptions: gitlab.ListOptions{
				Page:    1,
				PerPage: 100,
			},
			Scope: &[]gitlab.BuildStateValue{"failed", "success"},
		}

		for {
			jobs, _, err := client.Jobs.ListProjectJobs(p.ID, jobsOpts)
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

				// check that the job needs cleanup before sending it
				if j.NeedCleanup(opts.ThresholdSize, opts.ThresholdTime) {
					in <- j
				}
			}
		}
	}
}

// Merge merges all Project's jobs and returns the project.
func (p Project) Merge(out <-chan Job) Project {
	for job := range out {
		if job.Cleaned {
			p.JobsCleaned++
		}
	}
	return p
}

// ReadProjects reads all projects from gitlab api and send them into the output channel.
// The output channel is closed once all projects were sent into it.
func ReadProjects(ctx context.Context, client *gitlab.Client, opts Options) <-chan Project {
	// un-buffered channel to avoid too many pages in memory
	tasks := make(chan Project)

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

				// confirm that project is inside cleanup regexp slice
				_, found := lo.Find(opts.PathRegexps, func(reg *regexp.Regexp) bool {
					return reg.MatchString(p.PathWithNamespace)
				})
				if found {
					tasks <- p
				}
			}
		}
	}()

	return tasks
}
