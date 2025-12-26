package models

import (
	"regexp"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

// Project is a simplified view of a gitlab project with only useful information used during artifacts command.
type Project struct {
	ID                int64
	PathWithNamespace string
	JobsCleaned       int
}

// Matches returns truthy if the project path matches any of the provided regexps.
func (p Project) Matches(regexps ...*regexp.Regexp) bool {
	for _, r := range regexps {
		if r.MatchString(p.PathWithNamespace) {
			return true
		}
	}
	return false
}

// ProjectFromGitLab converts a GitLab project to its simplified view.
func ProjectFromGitLab(project *gitlab.Project) Project {
	return Project{
		ID:                project.ID,
		PathWithNamespace: project.PathWithNamespace,
	}
}
