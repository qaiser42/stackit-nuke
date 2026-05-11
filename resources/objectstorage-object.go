package resources

import (
	"context"
	"fmt"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/qaiser42/stackit-nuke/pkg/stackit"
)

const ObjectStorageObjectResource = "ObjectStorageObject"

func init() {
	registry.Register(&registry.Registration{
		Name:      ObjectStorageObjectResource,
		Scope:     stackit.ProjectScope,
		Resource:  &ObjectStorageObject{},
		Lister:    &ObjectStorageObjectLister{},
		DependsOn: []string{ObjectStorageBucketResource},
	})
}

type ObjectStorageObject struct {
	*BaseResource `property:",inline"`
	Bucket        string
	Key           string
}

func (r *ObjectStorageObject) Remove(_ context.Context) error {
	return fmt.Errorf("ObjectStorageObject.Remove not yet implemented")
}
func (r *ObjectStorageObject) Properties() types.Properties { return PropsFromStruct(r) }
func (r *ObjectStorageObject) String() string                { return r.Bucket + "/" + r.Key }

type ObjectStorageObjectLister struct{}

func (l *ObjectStorageObjectLister) List(_ context.Context, _ any) ([]resource.Resource, error) {
	return []resource.Resource{}, nil
}
