package resources

import (
	"context"
	"fmt"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/qaiser42/stackit-nuke/pkg/stackit"
)

const SubnetResource = "Subnet"

func init() {
	registry.Register(&registry.Registration{
		Name:      SubnetResource,
		Scope:     stackit.ProjectScope,
		Resource:  &Subnet{},
		Lister:    &SubnetLister{},
		DependsOn: []string{NetworkResource},
	})
}

type Subnet struct {
	*BaseResource `property:",inline"`
	ID            string
	Name          string
	NetworkID     string
}

func (r *Subnet) Remove(_ context.Context) error {
	return fmt.Errorf("Subnet.Remove not yet implemented")
}
func (r *Subnet) Properties() types.Properties { return PropsFromStruct(r) }
func (r *Subnet) String() string               { return r.Name }

type SubnetLister struct{}

func (l *SubnetLister) List(_ context.Context, _ any) ([]resource.Resource, error) {
	return []resource.Resource{}, nil
}
