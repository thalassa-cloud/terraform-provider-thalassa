package objectstorage

import (
	"fmt"
	"time"

	"github.com/thalassa-cloud/client-go/iaas"
	"github.com/thalassa-cloud/client-go/pkg/base"
)

type ObjectStorageBucketVersioning string

const (
	ObjectStorageBucketVersioningDisabled  ObjectStorageBucketVersioning = "Disabled"
	ObjectStorageBucketVersioningEnabled   ObjectStorageBucketVersioning = "Enabled"
	ObjectStorageBucketVersioningSuspended ObjectStorageBucketVersioning = "Suspended"
)

type ObjectStorageBucket struct {
	Identity     string             `json:"identity"`
	Organisation *base.Organisation `json:"organisation,omitempty"`
	CreatedAt    time.Time          `json:"createdAt"`
	UpdatedAt    time.Time          `json:"updatedAt"`

	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`

	// Name is the name of the bucket
	Name string `json:"name"`
	// Policy is the policy of the bucket
	Policy PolicyDocument `json:"policy"`
	// Public is a flag that indicates if the bucket is public
	Public bool `json:"public"`
	// Status is the status of the bucket
	Status string `json:"status"`
	// Endpoint for the bucket
	Endpoint string `json:"endpoint"`
	// Usage is the usage of the bucket
	Usage ObjectStorageBucketUsage `json:"usage"`
	// Versioning is the versioning of the bucket
	Versioning ObjectStorageBucketVersioning `json:"versioning"`
	// ObjectLockEnabled is the object lock of the bucket
	ObjectLockEnabled bool `json:"objectLockEnabled"`
	// Region is the region of the bucket
	Region *iaas.Region `json:"cloudRegion,omitempty"`
	// Lifecycle is the bucket lifecycle configuration.
	Lifecycle *BucketLifecycle `json:"lifecycle,omitempty"`
}

type ObjectStorageBucketUsage struct {
	TotalSizeGB  float64 `json:"total_size_gb"`
	TotalObjects int64   `json:"total_objects"`
}

type CreateBucketRequest struct {
	// BucketName is the name of the bucket.
	BucketName string `json:"bucketName"`
	// Public is a flag that indicates if the bucket can be accessed by the public.
	// When set to false, it blocks all public access to the bucket.
	Public bool `json:"public"`
	// Region is the region of the bucket.
	Region string `json:"region"`
	// PolicyDocument is the policy document for the bucket.
	PolicyDocument *PolicyDocument `json:"policy,omitempty"`
	// Labels is the labels of the bucket.
	Labels map[string]string `json:"labels,omitempty"`
	// Annotations is the annotations of the bucket.
	Annotations map[string]string `json:"annotations,omitempty"`
	// Versioning is the versioning of the bucket.
	Versioning ObjectStorageBucketVersioning `json:"versioning"`
	// ObjectLockEnabled is the object lock enabled of the bucket.
	ObjectLockEnabled bool `json:"objectLockEnabled"`
	// Lifecycle optionally sets lifecycle rules at bucket creation.
	Lifecycle *SetBucketLifecycleRequest `json:"lifecycle,omitempty"`
}

type UpdateBucketRequest struct {
	// Public is a flag that indicates if the bucket can be accessed by the public.
	Public bool `json:"public"`
	// PolicyDocument is the policy document for the bucket.
	PolicyDocument *PolicyDocument `json:"policy,omitempty"`
	// Versioning is the versioning of the bucket.
	Versioning ObjectStorageBucketVersioning `json:"versioning"`
	// ObjectLockEnabled is the object lock enabled of the bucket.
	ObjectLockEnabled *bool `json:"objectLockEnabled"`

	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

// PolicyDocument represents a full S3 bucket policy.
type PolicyDocument struct {
	Version   string      `json:"Version"`
	Statement []Statement `json:"Statement"`
}

// Statement defines an individual rule in the policy.
type Statement struct {
	Sid       string      `json:"Sid,omitempty"`
	Effect    string      `json:"Effect"`
	Principal Principal   `json:"Principal"`
	Action    interface{} `json:"Action"` // can be string or []string
	Resource  []string    `json:"Resource"`
	Condition interface{} `json:"Condition,omitempty"`
}

// Principal defines which user(s) the statement applies to.
type Principal struct {
	AWS      interface{} `json:"AWS,omitempty"`      // can be string or []string
	Thalassa interface{} `json:"Thalassa,omitempty"` // can be string or []string
}

type PrincipalARN string

func (p PrincipalARN) Validate() error {
	if p == "" {
		return fmt.Errorf("principal ARN is required")
	}
	return nil
}

type BucketLifecycle struct {
	Rules []BucketLifecycleRule `json:"rules"`
}

type BucketLifecycleRule struct {
	ID                             string                                             `json:"id"`
	Prefix                         string                                             `json:"prefix,omitempty"`
	Filter                         *BucketLifecycleRuleFilter                         `json:"filter,omitempty"`
	Status                         BucketLifecycleRuleStatus                          `json:"status,omitempty"`
	Expiration                     *BucketLifecycleRuleExpiration                     `json:"expiration,omitempty"`
	Transitions                    []BucketLifecycleRuleTransition                    `json:"transitions,omitempty"`
	NoncurrentVersionExpiration    *BucketLifecycleRuleNoncurrentVersionExpiration    `json:"noncurrentVersionExpiration,omitempty"`
	NoncurrentVersionTransitions   []BucketLifecycleRuleNoncurrentVersionTransition   `json:"noncurrentVersionTransitions,omitempty"`
	AbortIncompleteMultipartUpload *BucketLifecycleRuleAbortIncompleteMultipartUpload `json:"abortIncompleteMultipartUpload,omitempty"`
}

type BucketLifecycleRuleStatus string

const (
	BucketLifecycleRuleStatusEnabled  BucketLifecycleRuleStatus = "Enabled"
	BucketLifecycleRuleStatusDisabled BucketLifecycleRuleStatus = "Disabled"
)

type BucketLifecycleRuleFilter struct {
	Prefix                string                          `json:"prefix,omitempty"`
	Tag                   *BucketLifecycleRuleTag         `json:"tag,omitempty"`
	And                   *BucketLifecycleRuleAndOperator `json:"and,omitempty"`
	ObjectSizeGreaterThan *int64                          `json:"objectSizeGreaterThan,omitempty"`
	ObjectSizeLessThan    *int64                          `json:"objectSizeLessThan,omitempty"`
}

type BucketLifecycleRuleTag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type BucketLifecycleRuleAndOperator struct {
	Prefix                string                   `json:"prefix,omitempty"`
	Tags                  []BucketLifecycleRuleTag `json:"tags,omitempty"`
	ObjectSizeGreaterThan *int64                   `json:"objectSizeGreaterThan,omitempty"`
	ObjectSizeLessThan    *int64                   `json:"objectSizeLessThan,omitempty"`
}

type BucketLifecycleRuleExpiration struct {
	Days                      *int64     `json:"days,omitempty"`
	Date                      *time.Time `json:"date,omitempty"`
	ExpiredObjectDeleteMarker *bool      `json:"expiredObjectDeleteMarker,omitempty"`
}

type BucketLifecycleRuleTransition struct {
	Days         *int64     `json:"days,omitempty"`
	Date         *time.Time `json:"date,omitempty"`
	StorageClass string     `json:"storageClass"`
}

type BucketLifecycleRuleNoncurrentVersionExpiration struct {
	NoncurrentDays *int64 `json:"noncurrentDays,omitempty"`
}

type BucketLifecycleRuleNoncurrentVersionTransition struct {
	NoncurrentDays *int64 `json:"noncurrentDays,omitempty"`
	StorageClass   string `json:"storageClass"`
}

type BucketLifecycleRuleAbortIncompleteMultipartUpload struct {
	DaysAfterInitiation *int64 `json:"daysAfterInitiation,omitempty"`
}

type SetBucketLifecycleRequest struct {
	Rules []BucketLifecycleRule `json:"rules"`
}
