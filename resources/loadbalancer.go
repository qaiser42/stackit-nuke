package resources

import (
	"context"
	"fmt"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/qaiser42/stackit-nuke/pkg/stackit"
)

const LoadBalancerResource = "LoadBalancer"

func init() {
	registry.Register(&registry.Registration{
		Name:     LoadBalancerResource,
		Scope:    stackit.ProjectScope,
		Resource: &LoadBalancer{},
		Lister:   &LoadBalancerLister{},
	})
}

type LoadBalancer struct {
	*BaseResource `property:",inline"`
	Name          string
}

func (r *LoadBalancer) Remove(_ context.Context) error {
	return fmt.Errorf("LoadBalancer.Remove not yet implemented")
}
func (r *LoadBalancer) Properties() types.Properties { return PropsFromStruct(r) }
func (r *LoadBalancer) String() string               { return r.Name }

type LoadBalancerLister struct{}

func (l *LoadBalancerLister) List(_ context.Context, _ any) ([]resource.Resource, error) {
	return []resource.Resource{}, nil
}
