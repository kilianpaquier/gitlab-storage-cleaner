package artifacts

import (
	"context"

	pooling "github.com/kilianpaquier/pooling/pkg"
	"github.com/samber/lo"
	gitlab "gitlab.com/gitlab-org/api/client-go"

	"github.com/kilianpaquier/gitlab-storage-cleaner/internal/artifacts/engine"
	"github.com/kilianpaquier/gitlab-storage-cleaner/internal/artifacts/models"
)

// ReadProjects reads all projects from gitlab api and send them into the output channel.
//
// The output channel is closed once all projects were sent into it.
func ReadProjects(ctx context.Context, client *gitlab.Client, runOptions engine.RunOptions) <-chan pooling.PoolerFunc {
	logger := engine.GetLogger(ctx)

	// un-buffered channel to avoid too many pages in memory
	tasks := make(chan pooling.PoolerFunc)

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

			// stop infinite loop
			if len(projects) == 0 {
				break
			}
			opts.Page++

			// send all projects for cleanup and iterate to next page
			for _, gitlab := range projects {
				project := models.ProjectFromGitLab(gitlab)
				// confirm that project is inside cleanup slice
				if !project.Matches(runOptions.PathRegexps...) {
					logger.Info("skipping project cleaning",
						"project_id", project.ID,
						"project_path", project.PathWithNamespace)
					continue
				}
				tasks <- ReadJobs(ctx, client, project, runOptions)
			}
		}
	}()

	return tasks
}

// ReadJobs returns the function to clean artifacts a specific project.
//
// This function retrieves all project's jobs and send them into pooling PoolerFunc input channel.
func ReadJobs(ctx context.Context, client *gitlab.Client, project models.Project, runOptions engine.RunOptions) pooling.PoolerFunc {
	return func(funcs chan<- pooling.PoolerFunc) {
		logger := engine.GetLogger(ctx)

		logger.Info("running project cleanup",
			"project_id", project.ID,
			"project_path", project.PathWithNamespace)

		options := &gitlab.ListJobsOptions{
			ListOptions: gitlab.ListOptions{
				Page:    1,
				PerPage: 100,
			},
		}

		for {
			jobs, _, err := client.Jobs.ListProjectJobs(project.ID, options, gitlab.WithContext(ctx))
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
			options.Page++

			// send all jobs for cleanup and iterate to next page
			for _, gitlab := range jobs {
				job := models.JobFromGitLab(project.ID, gitlab)
				// check that the job needs to be cleaned up
				if job.NeedCleanup(runOptions.ThresholdSize, runOptions.ThresholdTime()) {
					funcs <- DeleteArtifacts(ctx, client, job, runOptions)
				}
			}
		}
		logger.Info("ended project cleanup",
			"project_id", project.ID,
			"project_path", project.PathWithNamespace)
	}
}

// DeleteArtifacts returns a pooling PoolerFunc to be executed in a specific pool to delete job's artifacts.
func DeleteArtifacts(ctx context.Context, client *gitlab.Client, job models.Job, runOptions engine.RunOptions) pooling.PoolerFunc {
	return func(chan<- pooling.PoolerFunc) {
		logger := engine.GetLogger(ctx)

		if runOptions.DryRun {
			logger.Info("running in dry run mode, skipping job's artifacts deletion",
				"job_id", job.ID,
				"project_id", job.ProjectID)
			return
		}

		if err := job.DeleteArtifacts(ctx, client); err != nil {
			logger.Warn("failed to delete job's artifacts",
				"error", err,
				"job", job.ID,
				"project", job.ProjectID)
			return
		}

		// job.Cleaned = true // cannot be used with pooling engine unless job is a pointer
	}
}
