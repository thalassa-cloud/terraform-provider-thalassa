package base

import "time"

type AppUser struct {
	Subject   string    `json:"subject"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

// Organisation represents an organisation in the system
type Organisation struct {
	// Identity is a unique identifier for the object
	Identity string `json:"identity"`
	// Name is the name of the organisation
	Name string `json:"name"`
	// Slug is the slug of the organisation
	Slug string `json:"slug"`
	// Description is the description of the organisation
	Description *string `json:"description,omitempty"`
	// Labels are the labels of the organisation
	Labels map[string]string `json:"labels,omitempty"`
	// Annotations are the annotations of the organisation
	Annotations   map[string]string `json:"annotations,omitempty"`
	CreatedAt     time.Time         `json:"createdAt"`
	UpdatedAt     *time.Time        `json:"updatedAt,omitempty"`
	ObjectVersion int               `json:"objectVersion"`
}

// OrganisationMemberType is a type that represents a role of a member in an organisation
type OrganisationMemberType string

const (
	// OrganisationMemberTypeOwner is a role that indicates that the user is an owner of the organisation
	OrganisationMemberTypeOwner OrganisationMemberType = "OWNER"
	// OrganisationMemberTypeMember is a role that indicates that the user is a member of the organisation
	OrganisationMemberTypeMember OrganisationMemberType = "MEMBER"
)

// OrganisationMember is a type that represents a member of an organisation
type OrganisationMember struct {
	// Identity is a unique identifier for the object
	Identity string `json:"identity"`
	// CreatedAt is the timestamp when the object was created
	CreatedAt time.Time `json:"createdAt"`
	// Organisation is the organisation that the user is a member of
	Organisation *Organisation `json:"organisation,omitempty"`
	// AppUser is the user that is a member of the organisation
	AppUser *AppUser `json:"user,omitempty"`
	// Role is the role of the user in the organisation
	Role OrganisationMemberType `json:"role"`
}
