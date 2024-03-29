// Code generated by go-builder-generator (https://github.com/kilianpaquier/go-builder-generator). DO NOT EDIT.

package tests

import (
	"time"

	"github.com/kilianpaquier/gitlab-storage-cleaner/internal/artifacts/v1"
)

// JobBuilder is an alias of Job to build Job with builder-pattern.
type JobBuilder artifacts.Job

// NewJobBuilder creates a new JobBuilder.
func NewJobBuilder() *JobBuilder {
	return &JobBuilder{}
}

// Copy reassigns the builder struct (behind pointer) to a new pointer and returns it.
func (b *JobBuilder) Copy() *JobBuilder {
	c := *b
	return &c
}

// Build returns built Job.
func (b *JobBuilder) Build() *artifacts.Job {
	c := (artifacts.Job)(*b)
	return &c
}

// SetArtifacts sets Job's Artifacts.
func (b *JobBuilder) SetArtifacts(artifacts ...artifacts.Artifact) *JobBuilder {
	b.Artifacts = append(b.Artifacts, artifacts...)
	return b
}

// SetArtifactsExpireAt sets Job's ArtifactsExpireAt.
func (b *JobBuilder) SetArtifactsExpireAt(artifactsExpireAt time.Time) *JobBuilder {
	b.ArtifactsExpireAt = artifactsExpireAt
	return b
}

// SetID sets Job's ID.
func (b *JobBuilder) SetID(id int) *JobBuilder {
	b.ID = id
	return b
}

// SetProjectID sets Job's ProjectID.
func (b *JobBuilder) SetProjectID(projectID int) *JobBuilder {
	b.ProjectID = projectID
	return b
}
