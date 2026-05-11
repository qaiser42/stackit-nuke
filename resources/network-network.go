package resources

import (
	"context"
	"fmt"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/qaiser42/stackit-nuke/pkg/stackit"
)

const NetworkResource = "Network"

func init() {
	registry.Register(&registry.Registration{
		Name:     NetworkResource,
		Scope:    stackit.ProjectScope,
		Resource: &Network{},
		Lister:   &NetworkLister{},
	})
}

type Network struct {
	*BaseResource `property:",inline"`
	ID            string
	Name          string
}

func (r *Network) Remove(_ context.Context) error {
	return fmt.Errorf("Network.Remove not yet implemented")
}
func (r *Network) Properties() types.Properties { return PropsFromStruct(r) }
func (r *Network) String() string                { return r.Name }

type NetworkLister struct{}

func (l *NetworkLister) List(_ context.Context, _ any) ([]resource.Resource, error) {
	return []resource.Resource{}, nil
}
