package artifacts

import (
	"context"
	"io"
	"time"

	pooling "github.com/kilianpaquier/pooling/pkg"
	"github.com/sirupsen/logrus"
	"gitlab.com/gitlab-org/api/client-go"
)

// Job is a simplified view of a gitlab job with only useful information for artifacts deletion feature.
type Job struct {
	Artifacts         []Artifact `builder:"append"`
	ArtifactsExpireAt time.Time
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

// DeleteArtifacts returns a pooling PoolerFunc to be executed in a specific pool to delete job's artifacts.
//
// If a specific job doesn't have artifacts, then nothing will be done in the returned function.
func (j Job) DeleteArtifacts(ctx context.Context, client *gitlab.Client) pooling.PoolerFunc {
	log := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"job":     j.ID,
		"project": j.ProjectID,
	})

	return func(_ chan<- pooling.PoolerFunc) {
		// call jobs artifacts deletion
		response, err := client.Jobs.DeleteArtifacts(j.ProjectID, j.ID, gitlab.WithContext(ctx))
		if err != nil {
			log.WithError(err).Warn("failed to delete job's artifacts")
			return
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
}
