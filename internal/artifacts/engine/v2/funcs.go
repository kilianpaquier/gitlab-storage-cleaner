package artifacts

import (
	"context"

	"github.com/fogfactory/pipe"
	"github.com/samber/lo"
	gitlab "gitlab.com/gitlab-org/api/client-go"

	"github.com/kilianpaquier/gitlab-storage-cleaner/internal/artifacts/engine"
	"github.com/kilianpaquier/gitlab-storage-cleaner/internal/artifacts/models"
)

// ReadProjects reads all projects from gitlab api and send them into the output channel.
// The output channel is closed once all projects were sent into it.
func ReadProjects(ctx context.Context, client *gitlab.Client, runOptions engine.RunOptions) <-chan Project {
	logger := engine.GetLogger(ctx)

	// un-buffered channel to avoid too many pages in memory
	tasks := make(chan Project)

	opts := &gitlab.ListProjectsOptions{
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
			projects, _, err := client.Projects.ListProjects(opts, gitlab.WithContext(ctx))
			if err != nil {
				logger.Warn("failed to retrieve projects", "error", err)
				break
			}
			opts.Page++

			// stop infinite loop
			if len(projects) == 0 {
				break
			}

			// send all projects for cleanup and iterate to next page
			for _, gitlab := range projects {
				project := models.ProjectFromGitLab(gitlab)

				if !project.Matches(runOptions.Regexps()...) {
					logger.Info("skipping project cleaning",
						"project_id", project.ID,
						"project_path", project.PathWithNamespace)
					continue
				}
				tasks <- Project{Project: project}
			}
		}
	}()

	return tasks
}

// ReadJobs returns the function to send all Jobs of a given Project into pipe processing.
func ReadJobs(ctx context.Context, client *gitlab.Client, runOptions engine.RunOptions) pipe.Split[Project, models.Job] {
	logger := engine.GetLogger(ctx)
	return func(project Project, in chan<- models.Job) {
		opts := &gitlab.ListJobsOptions{
			ListOptions: gitlab.ListOptions{
				Page:    1,
				PerPage: 100,
			},
			Scope: &[]gitlab.BuildStateValue{"failed", "success"},
		}

		for {
			jobs, _, err := client.Jobs.ListProjectJobs(project.ID, opts)
			if err != nil {
				logger.Warn("failed to retrieve project jobs",
					"error", err,
					"project_id", project.ID,
					"project_path", project.PathWithNamespace)
				break
			}

			// stop infinite loop
			if len(jobs) == 0 {
				break
			}
			opts.Page++

			for _, gitlab := range jobs {
				job := models.JobFromGitLab(project.ID, gitlab)
				// check that the job needs cleanup before sending it
				if job.NeedCleanup(runOptions.ThresholdDuration) {
					in <- job
				}
			}
		}
	}
}

// DeleteArtifacts returns the function to delete a specific job artifacts.
func DeleteArtifacts(ctx context.Context, client *gitlab.Client, opts engine.RunOptions) pipe.Process[models.Job] {
	return func(job models.Job) models.Job {
		logger := engine.GetLogger(ctx)

		if opts.DryRun {
			logger.Info("running in dry run mode, skipping job's artifacts deletion",
				"job_id", job.ID,
				"project_id", job.ProjectID)
			return job
		}

		if err := job.DeleteArtifacts(ctx, client); err != nil {
			logger.Warn("failed to delete job's artifacts",
				"error", err,
				"job_id", job.ID,
				"project_id", job.ProjectID)
			return job
		}

		job.Cleaned = true
		return job
	}
}

// ObserveCleanup merges all Project's jobs and returns the project.
func ObserveCleanup(project Project, out <-chan models.Job) Project {
	for job := range out {
		if job.Cleaned {
			project.JobsCleaned++
		}
	}
	return project
}

var _ pipe.Merge[Project, models.Job] = ObserveCleanup // ensure interface is implemented
