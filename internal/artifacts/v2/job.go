package artifacts

import (
	"io"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"
)

//go:generate go-builder-generator generate -f job.go -s Job -d tests

// Job is a simplified view of a gitlab job with only useful information for artifacts deletion feature.
type Job struct {
	Artifacts         []Artifact `builder:"append"`
	ArtifactsExpireAt time.Time
	Cleaned           bool
	ID                int
	ProjectID         int
}

// Artifact represents a simplified view of a gitlab artifact.
type Artifact struct {
	Size uint64
}

// NeedCleanup returns truthy if the job needs to be cleaned up.
func (j Job) NeedCleanup(thresholdSize uint64, thresholdTime time.Time) bool {
	// don't clean job not having artifacts
	if len(j.Artifacts) == 0 {
		return false
	}

	// don't clean job already cleaned by gitlab itself or clean up will be soon
	if !j.ArtifactsExpireAt.IsZero() && !j.ArtifactsExpireAt.After(thresholdTime) {
		return false
	}

	// compute artifacts max size
	var max uint64
	for _, artifact := range j.Artifacts {
		max += artifact.Size
	}

	// clean job if artifacts size is bigger than threshold size and expiration is zero (no expiration) or after threshold time
	return max >= thresholdSize && (j.ArtifactsExpireAt.IsZero() || j.ArtifactsExpireAt.After(thresholdTime))
}

// DeleteArtifacts returns the function to delete a specific job artifacts.
// Deletion will only occur if input opts has DryRun to false.
func DeleteArtifacts(client *gitlab.Client, opts Options) func(j Job) Job {
	return func(j Job) Job {
		log := logrus.WithFields(logrus.Fields{
			"job_id":     j.ID,
			"project_id": j.ProjectID,
		})

		if !opts.DryRun {
			// call jobs artifacts deletion
			response, err := client.Jobs.DeleteArtifacts(j.ProjectID, j.ID)
			if err != nil {
				log.WithError(err).Warn("failed to delete job's artifacts")
				return j
			}
			defer response.Body.Close()

			// handle http errors
			if response.StatusCode/100 != 2 {
				body, _ := io.ReadAll(response.Body)
				log.WithError(err).
					WithField("response", string(body)).
					Warn("failed to delete job's artifacts")
			}
		}

		j.Cleaned = true
		return j
	}
}
