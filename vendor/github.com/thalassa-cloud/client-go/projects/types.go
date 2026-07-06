package projects

import "time"

// Project represents a project within an organisation.
type Project struct {
	Identity      string            `json:"identity"`
	Name          string            `json:"name"`
	Slug          string            `json:"slug"`
	Description   string            `json:"description"`
	Labels        map[string]string `json:"labels"`
	Annotations   map[string]string `json:"annotations"`
	ParentProject *ProjectRef       `json:"parentProject,omitempty"`
	CreatedAt     time.Time         `json:"createdAt"`
	UpdatedAt     *time.Time        `json:"updatedAt,omitempty"`
	ObjectVersion int64             `json:"objectVersion"`
}

// ProjectRef is the shallow parent project reference returned in API responses.
type ProjectRef struct {
	Identity string `json:"identity"`
	Name     string `json:"name"`
	Slug     string `json:"slug"`
}

// CreateProjectRequest is the request body for POST /v1/projects.
type CreateProjectRequest struct {
	Name                  string            `json:"name"`
	Description           string            `json:"description"`
	Labels                map[string]string `json:"labels"`
	Annotations           map[string]string `json:"annotations"`
	ParentProjectIdentity *string           `json:"parentProjectIdentity,omitempty"`
}

// UpdateProjectRequest is the request body for PUT /v1/projects/{identity}.
type UpdateProjectRequest struct {
	Name                  string            `json:"name"`
	Description           string            `json:"description"`
	Labels                map[string]string `json:"labels"`
	Annotations           map[string]string `json:"annotations"`
	ParentProjectIdentity *string           `json:"parentProjectIdentity,omitempty"`
}
