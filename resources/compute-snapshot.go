package resources

import (
	"context"
	"fmt"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/qaiser42/stackit-nuke/pkg/stackit"
)

const ComputeSnapshotResource = "ComputeSnapshot"

func init() {
	registry.Register(&registry.Registration{
		Name:     ComputeSnapshotResource,
		Scope:    stackit.ProjectScope,
		Resource: &ComputeSnapshot{},
		Lister:   &ComputeSnapshotLister{},
	})
}

type ComputeSnapshot struct {
	*BaseResource `property:",inline"`
	ID            string
	Name          string
}

func (r *ComputeSnapshot) Remove(_ context.Context) error {
	return fmt.Errorf("ComputeSnapshot.Remove not yet implemented")
}
func (r *ComputeSnapshot) Properties() types.Properties { return PropsFromStruct(r) }
func (r *ComputeSnapshot) String() string               { return r.Name }

type ComputeSnapshotLister struct{}

func (l *ComputeSnapshotLister) List(_ context.Context, _ any) ([]resource.Resource, error) {
	return []resource.Resource{}, nil
}
