package resources

import (
	"context"
	"fmt"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/qaiser42/stackit-nuke/pkg/stackit"
)

const SecurityGroupResource = "SecurityGroup"

func init() {
	registry.Register(&registry.Registration{
		Name:     SecurityGroupResource,
		Scope:    stackit.ProjectScope,
		Resource: &SecurityGroup{},
		Lister:   &SecurityGroupLister{},
	})
}

type SecurityGroup struct {
	*BaseResource `property:",inline"`
	ID            string
	Name          string
}

func (r *SecurityGroup) Remove(_ context.Context) error {
	return fmt.Errorf("SecurityGroup.Remove not yet implemented")
}
func (r *SecurityGroup) Properties() types.Properties { return PropsFromStruct(r) }
func (r *SecurityGroup) String() string                { return r.Name }

type SecurityGroupLister struct{}

func (l *SecurityGroupLister) List(_ context.Context, _ any) ([]resource.Resource, error) {
	return []resource.Resource{}, nil
}
