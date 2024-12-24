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
	Artifacts         []Artifact
	ArtifactsExpireAt time.Time
	Cleaned           bool
	ID                int
	ProjectID         int
}

// Artifact represents a simplified view of a gitlab artifact.
type Artifact struct {
	Size int
}

// NeedCleanup returns truthy if the job needs to be cleaned up.
func (j Job) NeedCleanup(thresholdSize int, thresholdTime time.Time) bool {
	// don't clean job not having artifacts
	if len(j.Artifacts) == 0 {
		return false
	}

	// don't clean job already cleaned by gitlab itself or clean up will be soon
	if !j.ArtifactsExpireAt.IsZero() && !j.ArtifactsExpireAt.After(thresholdTime) {
		return false
	}

	// compute artifacts maxSize size
	var maxSize int
	for _, artifact := range j.Artifacts {
		maxSize += artifact.Size
	}

	// clean job if artifacts size is bigger than threshold size and expiration is zero (no expiration) or after threshold time
	return maxSize >= thresholdSize && (j.ArtifactsExpireAt.IsZero() || j.ArtifactsExpireAt.After(thresholdTime))
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
func JobFromGitLab(projectID int, job *gitlab.Job) Job {
	// convert artifacts
	artifacts := make([]Artifact, len(job.Artifacts))
	for i, artifact := range job.Artifacts {
		artifacts[i] = Artifact{Size: artifact.Size}
	}

	// convert job
	return Job{
		Artifacts:         artifacts,
		ArtifactsExpireAt: lo.FromPtr(job.ArtifactsExpireAt),
		ID:                job.ID,
		ProjectID:         projectID,
	}
}
