package artifacts

import (
	"context"
	"time"

	"github.com/kilianpaquier/gitlab-storage-cleaner/internal/artifacts/engine"
	"github.com/kilianpaquier/gitlab-storage-cleaner/internal/artifacts/models"
)

// Project represents a models.Project
// with additional execution information specific
// to features available with pipe (start and stop logs).
type Project struct {
	models.Project

	executionStart    time.Time
	executionDuration time.Duration
}

// StartProject starts the timer for Project execution and logs the Project start execution.
func StartProject(ctx context.Context) func(Project) Project {
	return func(p Project) Project {
		p.executionStart = time.Now()
		engine.GetLogger(ctx).Info("starting project execution",
			"project_id", p.ID,
			"project_path", p.PathWithNamespace)
		return p
	}
}

// StopProject stops the project timer execution and logs the Project execution result.
func StopProject(ctx context.Context) func(Project) Project {
	return func(p Project) Project {
		p.executionDuration = time.Since(p.executionStart)
		engine.GetLogger(ctx).Info("ending project execution",
			"execution_duration", p.executionDuration,
			"jobs_cleaned", p.JobsCleaned,
			"project_id", p.ID,
			"project_path", p.PathWithNamespace)
		return p
	}
}
