package objectstorage

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thalassa-cloud/client-go/objectstorage"
)

func TestValidateThalassaPrincipalARN(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		arn     string
		wantErr string
	}{
		{
			name: "valid service account principal",
			arn:  "arn:thalassa:iam:::serviceaccount/o-org123:sa-account456",
		},
		{
			name:    "invalid colon-separated principal",
			arn:     "arn:thalassa:iam:::serviceaccount:o-org123:sa-account456",
			wantErr: "invalid Principal.Thalassa ARN",
		},
		{
			name:    "empty principal",
			arn:     "",
			wantErr: "cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := validateThalassaPrincipalARN(tt.arn)
			if tt.wantErr == "" {
				assert.NoError(t, err)
				return
			}

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func TestValidateBucketPolicyDocument(t *testing.T) {
	t.Parallel()

	err := validateBucketPolicyDocument(objectstorage.PolicyDocument{
		Statement: []objectstorage.Statement{
			{
				Principal: objectstorage.Principal{
					Thalassa: []any{"arn:thalassa:iam:::organisation/o-org123"},
				},
			},
		},
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "policy statement 0")
	assert.Contains(t, err.Error(), "invalid Principal.Thalassa ARN")
}

func TestParseBucketPolicyJSON(t *testing.T) {
	t.Parallel()

	_, err := parseBucketPolicyJSON(`{
		"Version": "2012-10-17",
		"Statement": [{
			"Effect": "Allow",
			"Action": ["s3:ListBucket"],
			"Resource": ["arn:thalassa:s3:::example"],
			"Principal": {
				"Thalassa": ["arn:thalassa:iam:::serviceaccount/o-org123:sa-account456"]
			}
		}]
	}`)
	assert.NoError(t, err)

	_, err = parseBucketPolicyJSON(`{"Statement":[{"Principal":{"Thalassa":["arn:thalassa:iam:::organisation/o-org123"]}}]}`)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid Principal.Thalassa ARN")
}

func TestEquivalentPolicyJSON(t *testing.T) {
	t.Parallel()

	compact := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":["s3:ListBucket"],"Resource":["arn:thalassa:s3:::example"]}]}`
	pretty := `{
		"Version": "2012-10-17",
		"Statement": [{
			"Effect": "Allow",
			"Action": ["s3:ListBucket"],
			"Resource": ["arn:thalassa:s3:::example"]
		}]
	}`

	assert.True(t, equivalentPolicyJSON(compact, pretty))
	assert.False(t, equivalentPolicyJSON(compact, `{"Version":"2012-10-17","Statement":[]}`))
}

func TestBucketPolicyStateValue(t *testing.T) {
	t.Parallel()

	configured := "{\n  \"Version\": \"2012-10-17\"\n}\n"
	apiPolicy := `{"Version":"2012-10-17"}`

	assert.Equal(t, configured, bucketPolicyStateValue(configured, apiPolicy))
	assert.Equal(t, apiPolicy, bucketPolicyStateValue("", apiPolicy))
}

func TestEnrichBucketError(t *testing.T) {
	t.Parallel()

	err := enrichBucketError(fmt.Errorf("bad request: invalid principal ARN"), "create")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create bucket")
	assert.Contains(t, err.Error(), "serviceaccount/<organisation-id>:<service-account-id>")
}
