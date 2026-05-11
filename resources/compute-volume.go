package resources

import (
	"context"
	"fmt"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/qaiser42/stackit-nuke/pkg/stackit"
)

const ComputeVolumeResource = "ComputeVolume"

func init() {
	registry.Register(&registry.Registration{
		Name:      ComputeVolumeResource,
		Scope:     stackit.ProjectScope,
		Resource:  &ComputeVolume{},
		Lister:    &ComputeVolumeLister{},
		DependsOn: []string{ComputeServerResource},
	})
}

type ComputeVolume struct {
	*BaseResource `property:",inline"`
	ID            string
	Name          string
}

func (r *ComputeVolume) Remove(_ context.Context) error {
	return fmt.Errorf("ComputeVolume.Remove not yet implemented")
}
func (r *ComputeVolume) Properties() types.Properties { return PropsFromStruct(r) }
func (r *ComputeVolume) String() string               { return r.Name }

type ComputeVolumeLister struct{}

func (l *ComputeVolumeLister) List(_ context.Context, _ any) ([]resource.Resource, error) {
	return []resource.Resource{}, nil
}
