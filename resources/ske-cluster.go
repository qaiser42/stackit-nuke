package resources

import (
	"context"
	"fmt"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/qaiser42/stackit-nuke/pkg/stackit"
)

const SKEClusterResource = "SKECluster"

func init() {
	registry.Register(&registry.Registration{
		Name:     SKEClusterResource,
		Scope:    stackit.ProjectScope,
		Resource: &SKECluster{},
		Lister:   &SKEClusterLister{},
	})
}

type SKECluster struct {
	*BaseResource `property:",inline"`
	Name          string
}

func (r *SKECluster) Remove(_ context.Context) error {
	return fmt.Errorf("SKECluster.Remove not yet implemented")
}
func (r *SKECluster) Properties() types.Properties { return PropsFromStruct(r) }
func (r *SKECluster) String() string               { return r.Name }

type SKEClusterLister struct{}

func (l *SKEClusterLister) List(_ context.Context, _ any) ([]resource.Resource, error) {
	return []resource.Resource{}, nil
}
