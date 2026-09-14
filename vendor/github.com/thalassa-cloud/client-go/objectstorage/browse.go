package objectstorage

import "time"

// BrowseObjectsRequest carries query parameters for listing objects under a prefix.
type BrowseObjectsRequest struct {
	Prefix            string
	Delimiter         string
	MaxKeys           int32
	ContinuationToken string
}

// BrowseObject is a single object returned by a browse listing.
type BrowseObject struct {
	Key          string    `json:"key"`
	Size         int64     `json:"size"`
	LastModified time.Time `json:"lastModified"`
	ETag         string    `json:"etag,omitempty"`
	StorageClass string    `json:"storageClass,omitempty"`
}

// BrowseResponse is a paginated listing of prefixes and objects.
type BrowseResponse struct {
	Prefix                string         `json:"prefix"`
	Prefixes              []string       `json:"prefixes"`
	Objects               []BrowseObject `json:"objects"`
	IsTruncated           bool           `json:"isTruncated"`
	ContinuationToken     string         `json:"continuationToken,omitempty"`
	NextContinuationToken string         `json:"nextContinuationToken,omitempty"`
}

// ObjectMetadata is metadata for a single object.
type ObjectMetadata struct {
	Key          string    `json:"key"`
	Size         int64     `json:"size"`
	LastModified time.Time `json:"lastModified"`
	ETag         string    `json:"etag,omitempty"`
	ContentType  string    `json:"contentType,omitempty"`
	StorageClass string    `json:"storageClass,omitempty"`
}
