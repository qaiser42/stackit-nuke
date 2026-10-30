package resources

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/stackitcloud/stackit-sdk-go/core/oapierror"
	objectstoragev2 "github.com/stackitcloud/stackit-sdk-go/services/objectstorage/v2api"

	"github.com/qaiser42/stackit-nuke/pkg/stackit"
)

const ObjectStorageBucketResource = "ObjectStorageBucket"

func init() {
	registry.Register(&registry.Registration{
		Name:     ObjectStorageBucketResource,
		Scope:    stackit.ProjectScope,
		Resource: &ObjectStorageBucket{},
		Lister:   &ObjectStorageBucketLister{},
	})
}

// ObjectStorageBucket is a STACKIT Object Storage bucket.
//
// API: github.com/stackitcloud/stackit-sdk-go/services/objectstorage/v2api
//
// Endpoints used:
//   - GET    /v2/project/{projectId}/regions/{region}/buckets
//   - DELETE /v2/project/{projectId}/regions/{region}/bucket/{bucketName}
//
// Deleting a bucket also removes every object it contains; objects are not
// a separate resource type (the management API has no object endpoints).
// Projects without Object Storage enabled answer the list call with 404,
// which the Lister treats as "no buckets".
type ObjectStorageBucket struct {
	*BaseResource `property:",inline"`

	api objectstoragev2.DefaultAPI

	Name              string
	ObjectLockEnabled bool
	URL               string
}

func (r *ObjectStorageBucket) Remove(ctx context.Context) error {
	if r.api == nil {
		return fmt.Errorf("ObjectStorageBucket.Remove: api client not set")
	}
	if _, err := r.api.DeleteBucket(ctx, r.ProjectID, r.Region, r.Name).Execute(); err != nil {
		return fmt.Errorf("delete bucket: %w", err)
	}
	return nil
}

func (r *ObjectStorageBucket) Properties() types.Properties { return PropsFromStruct(r) }
func (r *ObjectStorageBucket) String() string               { return r.Name }

type ObjectStorageBucketLister struct{}

func (l *ObjectStorageBucketLister) List(ctx context.Context, o any) ([]resource.Resource, error) {
	opts := o.(*stackit.ListerOpts)
	logger := logrus.WithFields(logrus.Fields{
		"resource": ObjectStorageBucketResource,
		"project":  opts.ProjectID,
		"region":   opts.Region,
	})

	if opts.Credentials == nil {
		return nil, fmt.Errorf("missing STACKIT credentials")
	}

	client, err := objectstoragev2.NewAPIClient(stackitConfigOpts(opts)...)
	if err != nil {
		return nil, fmt.Errorf("build objectstorage client: %w", err)
	}

	logger.Trace("listing buckets")
	resp, err := client.DefaultAPI.ListBuckets(ctx, opts.ProjectID, opts.Region).Execute()
	if err != nil {
		var oapiErr *oapierror.GenericOpenAPIError
		if errors.As(err, &oapiErr) && oapiErr.StatusCode == http.StatusNotFound {
			logger.Trace("object storage not enabled for project")
			return []resource.Resource{}, nil
		}
		return nil, fmt.Errorf("list buckets: %w", err)
	}

	buckets := resp.GetBuckets()
	out := make([]resource.Resource, 0, len(buckets))
	for _, b := range buckets {
		if b.GetName() == "" {
			continue
		}
		out = append(out, &ObjectStorageBucket{
			BaseResource: &BaseResource{
				OrganizationID: opts.OrganizationID,
				ProjectID:      opts.ProjectID,
				Region:         opts.Region,
			},
			api:               client.DefaultAPI,
			Name:              b.GetName(),
			ObjectLockEnabled: b.GetObjectLockEnabled(),
			URL:               b.GetUrlVirtualHostedStyle(),
		})
	}

	logger.WithField("count", len(out)).Trace("list complete")
	return out, nil
}
