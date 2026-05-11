package resources

import (
	"context"
	"fmt"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

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

type ObjectStorageBucket struct {
	*BaseResource `property:",inline"`
	Name          string
}

func (r *ObjectStorageBucket) Remove(_ context.Context) error {
	return fmt.Errorf("ObjectStorageBucket.Remove not yet implemented")
}
func (r *ObjectStorageBucket) Properties() types.Properties { return PropsFromStruct(r) }
func (r *ObjectStorageBucket) String() string               { return r.Name }

type ObjectStorageBucketLister struct{}

func (l *ObjectStorageBucketLister) List(_ context.Context, _ any) ([]resource.Resource, error) {
	return []resource.Resource{}, nil
}
