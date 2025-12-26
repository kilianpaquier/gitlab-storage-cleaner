package models

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/samber/lo"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

// Job is a simplified view of a gitlab job with only useful information for artifacts deletion feature.
type Job struct {
	ArtifactsCount    int
	ArtifactsExpireAt time.Time
	Cleaned           bool
	CreatedAt         time.Time
	ID                int64
	ProjectID         int64
}

// Artifact represents a simplified view of a gitlab artifact.
type Artifact struct{}

// NeedCleanup returns truthy if the job needs to be cleaned up.
//
// It returns true if (all conditions are met):
//   - the job has artifacts
//   - the job creation date is undefined or before now minus the threshold
//   - the job artifacts expiration date is defined and after now
//
// It returns false if (any condition is met):
//   - the job has no artifacts
//   - the job creation date is defined and after now minus the threshold
//   - the job artifacts expiration date is already passed
func (j Job) NeedCleanup(threshold time.Duration) bool {
	// don't clean job not having artifacts
	if j.ArtifactsCount == 0 {
		return false
	}
	now := time.Now()

	// already cleaned up by GitLab
	expired := !j.ArtifactsExpireAt.IsZero() && j.ArtifactsExpireAt.Before(now)

	// creation issue or before threshold
	old := j.CreatedAt.IsZero() || j.CreatedAt.Before(now.Add(-threshold))
	return old && !expired
}

// DeleteArtifacts deletes the artifacts of the job.
//
// It returns an error if the deletion failed.
func (j Job) DeleteArtifacts(ctx context.Context, client *gitlab.Client) error {
	// call jobs artifacts deletion
	response, err := client.Jobs.DeleteArtifacts(j.ProjectID, j.ID, gitlab.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("delete artifacts: %w", err)
	}
	defer response.Body.Close()

	// handle http errors
	if response.StatusCode/100 != 2 {
		body, _ := io.ReadAll(response.Body)
		return fmt.Errorf("delete artifacts: %s", body)
	}
	return nil
}

// JobFromGitLab converts a GitLab job to its simplified view.
func JobFromGitLab(projectID int64, job *gitlab.Job) Job {
	return Job{
		ArtifactsCount:    len(job.Artifacts),
		ArtifactsExpireAt: lo.FromPtr(job.ArtifactsExpireAt),
		CreatedAt:         lo.FromPtr(job.CreatedAt),
		ID:                job.ID,
		ProjectID:         projectID,
	}
}
