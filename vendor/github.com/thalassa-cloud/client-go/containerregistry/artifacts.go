package containerregistry

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/thalassa-cloud/client-go/pkg/client"
)

// PagedResult is a paginated list response.
type PagedResult[T any] struct {
	Items        []T `json:"items"`
	TotalItems   int `json:"totalItems"`
	Page         int `json:"page"`
	ItemsPerPage int `json:"itemsPerPage"`
	TotalPages   int `json:"totalPages"`
	TotalCount   int `json:"totalCount"`
}

// ContainerRegistryArtifactResponse is the organisation API DTO for an artifact.
type ContainerRegistryArtifactResponse struct {
	CreatedAt    time.Time              `json:"created_at"`
	LastPulledAt *time.Time             `json:"last_pulled_at,omitempty"`
	MediaType    string                 `json:"media_type"`
	Size         float64                `json:"size"`
	Digest       string                 `json:"digest"`
	Tags         []string               `json:"tags"`
	Layers       map[string]interface{} `json:"layers,omitempty"`
	RawData      map[string]interface{} `json:"raw_data,omitempty"`
}

// ContainerRegistryTagResponse is the organisation API DTO for a repository tag.
type ContainerRegistryTagResponse struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Tag       string    `json:"tag"`
	Sha256    string    `json:"sha256"`
	Digest    string    `json:"digest"`
	SizeMb    float64   `json:"size_mb"`
}

// ListContainerRegistryArtifactsRequest is the request for listing repository artifacts.
type ListContainerRegistryArtifactsRequest struct {
	Page      int
	Limit     int
	Digest    string
	MediaType string
	Tag       string
	Untagged  *bool
}

// ListContainerRegistryTagsRequest is the request for listing repository tags.
type ListContainerRegistryTagsRequest struct {
	Page   int
	Limit  int
	Tag    string
	Digest string
}

// ListContainerRegistryArtifacts lists artifacts for a repository.
func (c *Client) ListContainerRegistryArtifacts(ctx context.Context, namespaceIdentity string, repositoryIdentity string, listRequest *ListContainerRegistryArtifactsRequest) (*PagedResult[ContainerRegistryArtifactResponse], error) {
	if namespaceIdentity == "" {
		return nil, fmt.Errorf("namespace identity is required")
	}
	if repositoryIdentity == "" {
		return nil, fmt.Errorf("repository identity is required")
	}

	result := &PagedResult[ContainerRegistryArtifactResponse]{}
	req := c.R().SetResult(result)
	if listRequest != nil {
		if listRequest.Page > 0 {
			req = req.SetQueryParam("page", strconv.Itoa(listRequest.Page))
		}
		if listRequest.Limit > 0 {
			req = req.SetQueryParam("limit", strconv.Itoa(listRequest.Limit))
		}
		if listRequest.Digest != "" {
			req = req.SetQueryParam("digest", listRequest.Digest)
		}
		if listRequest.MediaType != "" {
			req = req.SetQueryParam("mediaType", listRequest.MediaType)
		}
		if listRequest.Tag != "" {
			req = req.SetQueryParam("tag", listRequest.Tag)
		}
		if listRequest.Untagged != nil {
			req = req.SetQueryParam("untagged", strconv.FormatBool(*listRequest.Untagged))
		}
	}

	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s/repositories/%s/artifacts", ContainerRegistryEndpoint, namespaceIdentity, repositoryIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return result, err
	}
	return result, nil
}

// GetContainerRegistryArtifact retrieves a repository artifact by digest.
func (c *Client) GetContainerRegistryArtifact(ctx context.Context, namespaceIdentity string, repositoryIdentity string, digest string) (*ContainerRegistryArtifactResponse, error) {
	if namespaceIdentity == "" {
		return nil, fmt.Errorf("namespace identity is required")
	}
	if repositoryIdentity == "" {
		return nil, fmt.Errorf("repository identity is required")
	}
	if digest == "" {
		return nil, fmt.Errorf("digest is required")
	}

	var artifact *ContainerRegistryArtifactResponse
	req := c.R().SetResult(&artifact)
	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s/repositories/%s/artifacts/%s", ContainerRegistryEndpoint, namespaceIdentity, repositoryIdentity, url.PathEscape(digest)))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return artifact, err
	}
	return artifact, nil
}

// DeleteContainerRegistryArtifact deletes a repository artifact by digest.
func (c *Client) DeleteContainerRegistryArtifact(ctx context.Context, namespaceIdentity string, repositoryIdentity string, digest string) error {
	if namespaceIdentity == "" {
		return fmt.Errorf("namespace identity is required")
	}
	if repositoryIdentity == "" {
		return fmt.Errorf("repository identity is required")
	}
	if digest == "" {
		return fmt.Errorf("digest is required")
	}

	req := c.R()
	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s/repositories/%s/artifacts/%s", ContainerRegistryEndpoint, namespaceIdentity, repositoryIdentity, url.PathEscape(digest)))
	if err != nil {
		return err
	}
	return c.Check(resp)
}

// ListContainerRegistryTags lists tags for a repository.
func (c *Client) ListContainerRegistryTags(ctx context.Context, namespaceIdentity string, repositoryIdentity string, listRequest *ListContainerRegistryTagsRequest) (*PagedResult[ContainerRegistryTagResponse], error) {
	if namespaceIdentity == "" {
		return nil, fmt.Errorf("namespace identity is required")
	}
	if repositoryIdentity == "" {
		return nil, fmt.Errorf("repository identity is required")
	}

	result := &PagedResult[ContainerRegistryTagResponse]{}
	req := c.R().SetResult(result)
	if listRequest != nil {
		if listRequest.Page > 0 {
			req = req.SetQueryParam("page", strconv.Itoa(listRequest.Page))
		}
		if listRequest.Limit > 0 {
			req = req.SetQueryParam("limit", strconv.Itoa(listRequest.Limit))
		}
		if listRequest.Tag != "" {
			req = req.SetQueryParam("tag", listRequest.Tag)
		}
		if listRequest.Digest != "" {
			req = req.SetQueryParam("digest", listRequest.Digest)
		}
	}

	resp, err := c.Do(ctx, req, client.GET, fmt.Sprintf("%s/%s/repositories/%s/tags", ContainerRegistryEndpoint, namespaceIdentity, repositoryIdentity))
	if err != nil {
		return nil, err
	}
	if err := c.Check(resp); err != nil {
		return result, err
	}
	return result, nil
}

// DeleteContainerRegistryTag deletes a repository tag.
func (c *Client) DeleteContainerRegistryTag(ctx context.Context, namespaceIdentity string, repositoryIdentity string, tag string) error {
	if namespaceIdentity == "" {
		return fmt.Errorf("namespace identity is required")
	}
	if repositoryIdentity == "" {
		return fmt.Errorf("repository identity is required")
	}
	if tag == "" {
		return fmt.Errorf("tag is required")
	}

	req := c.R()
	resp, err := c.Do(ctx, req, client.DELETE, fmt.Sprintf("%s/%s/repositories/%s/tags/%s", ContainerRegistryEndpoint, namespaceIdentity, repositoryIdentity, url.PathEscape(tag)))
	if err != nil {
		return err
	}
	return c.Check(resp)
}
