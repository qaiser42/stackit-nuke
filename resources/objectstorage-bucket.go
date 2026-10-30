package resources

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
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
//   - POST   /v2/project/{projectId}/regions/{region}/credentials-group
//   - DELETE /v2/project/{projectId}/regions/{region}/credentials-group/{groupId}
//   - POST   /v2/project/{projectId}/regions/{region}/access-key
//   - DELETE /v2/project/{projectId}/regions/{region}/access-key/{keyId}
//   - S3 ListObjectVersions + DeleteObjects on the bucket
//
// The API refuses to delete a bucket that still holds objects (422
// bucket.not_empty), and it has no object endpoints. On that error Remove
// creates a temporary credentials group + access key, deletes every object
// version and delete marker over S3, deletes the bucket, and then removes
// the temporary credentials again. Objects are therefore not a separate
// resource type.
//
// Projects without Object Storage enabled answer the list call with 404,
// which the Lister treats as "no buckets".
type ObjectStorageBucket struct {
	*BaseResource `property:",inline"`

	api        objectstoragev2.DefaultAPI
	s3Endpoint string

	Name              string
	ObjectLockEnabled bool
	URL               string
}

func (r *ObjectStorageBucket) Remove(ctx context.Context) error {
	if r.api == nil {
		return fmt.Errorf("ObjectStorageBucket.Remove: api client not set")
	}
	err := r.deleteBucket(ctx)
	if !isBucketNotEmpty(err) {
		return err
	}
	if err := r.empty(ctx); err != nil {
		return fmt.Errorf("empty bucket: %w", err)
	}
	return r.deleteBucket(ctx)
}

func (r *ObjectStorageBucket) deleteBucket(ctx context.Context) error {
	if _, err := r.api.DeleteBucket(ctx, r.ProjectID, r.Region, r.Name).Execute(); err != nil {
		return fmt.Errorf("delete bucket: %w", err)
	}
	return nil
}

func isBucketNotEmpty(err error) bool {
	var oapiErr *oapierror.GenericOpenAPIError
	return errors.As(err, &oapiErr) &&
		oapiErr.StatusCode == http.StatusUnprocessableEntity &&
		strings.Contains(string(oapiErr.Body), "bucket.not_empty")
}

// empty deletes every object version and delete marker in the bucket using
// a temporary S3 access key that is removed again afterwards.
func (r *ObjectStorageBucket) empty(ctx context.Context) (err error) {
	if r.s3Endpoint == "" {
		return fmt.Errorf("no S3 endpoint known for bucket %s", r.Name)
	}

	group, err := r.api.CreateCredentialsGroup(ctx, r.ProjectID, r.Region).
		CreateCredentialsGroupPayload(*objectstoragev2.NewCreateCredentialsGroupPayload(
			fmt.Sprintf("stackit-nuke-%d", time.Now().UnixNano()))).
		Execute()
	if err != nil {
		return fmt.Errorf("create temporary credentials group: %w", err)
	}
	groupID := group.CredentialsGroup.GetCredentialsGroupId()
	defer func() {
		cctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), time.Minute)
		defer cancel()
		if _, derr := r.api.DeleteCredentialsGroup(cctx, r.ProjectID, r.Region, groupID).Execute(); derr != nil && err == nil {
			err = fmt.Errorf("delete temporary credentials group %s: %w", groupID, derr)
		}
	}()

	// The expiry is a safety net in case the deferred cleanup never runs.
	payload := objectstoragev2.NewCreateAccessKeyPayload()
	payload.SetExpires(time.Now().Add(time.Hour))
	key, err := r.api.CreateAccessKey(ctx, r.ProjectID, r.Region).
		CredentialsGroup(groupID).
		CreateAccessKeyPayload(*payload).
		Execute()
	if err != nil {
		return fmt.Errorf("create temporary access key: %w", err)
	}
	defer func() {
		cctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), time.Minute)
		defer cancel()
		if _, derr := r.api.DeleteAccessKey(cctx, r.ProjectID, r.Region, key.GetKeyId()).CredentialsGroup(groupID).Execute(); derr != nil && err == nil {
			err = fmt.Errorf("delete temporary access key %s: %w", key.GetKeyId(), derr)
		}
	}()

	s3, err := minio.New(r.s3Endpoint, &minio.Options{
		Creds:        credentials.NewStaticV4(key.GetAccessKey(), key.GetSecretAccessKey(), ""),
		Secure:       true,
		Region:       r.Region,
		BucketLookup: minio.BucketLookupPath,
	})
	if err != nil {
		return fmt.Errorf("build s3 client: %w", err)
	}
	if err := waitForS3Access(ctx, s3, r.Name); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	objects := make(chan minio.ObjectInfo)
	var listErr error
	go func() {
		defer close(objects)
		for o := range s3.ListObjects(ctx, r.Name, minio.ListObjectsOptions{Recursive: true, WithVersions: true}) {
			if o.Err != nil {
				listErr = o.Err
				return
			}
			select {
			case objects <- o:
			case <-ctx.Done():
				return
			}
		}
	}()

	var removeErr error
	for e := range s3.RemoveObjects(ctx, r.Name, objects, minio.RemoveObjectsOptions{}) {
		if removeErr == nil {
			removeErr = fmt.Errorf("remove object %s (version %s): %w", e.ObjectName, e.VersionID, e.Err)
		}
	}
	if listErr != nil {
		return fmt.Errorf("list objects: %w", listErr)
	}
	return removeErr
}

// waitForS3Access polls until a freshly created access key is accepted by
// the S3 gateway, which takes a few seconds.
func waitForS3Access(ctx context.Context, s3 *minio.Client, bucket string) error {
	deadline := time.Now().Add(time.Minute)
	for {
		_, err := s3.BucketExists(ctx, bucket)
		if err == nil {
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("s3 access with temporary key: %w", err)
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(2 * time.Second):
		}
	}
}

// s3EndpointFromURL returns the host of a path-style bucket URL, e.g.
// "object.storage.eu01.onstackit.cloud".
func s3EndpointFromURL(pathStyleURL string) string {
	u, err := url.Parse(pathStyleURL)
	if err != nil {
		return ""
	}
	return u.Host
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
			s3Endpoint:        s3EndpointFromURL(b.GetUrlPathStyle()),
			Name:              b.GetName(),
			ObjectLockEnabled: b.GetObjectLockEnabled(),
			URL:               b.GetUrlVirtualHostedStyle(),
		})
	}

	logger.WithField("count", len(out)).Trace("list complete")
	return out, nil
}
