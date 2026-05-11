package resources

import (
	"context"
	"fmt"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/qaiser42/stackit-nuke/pkg/stackit"
)

const RabbitMQInstanceResource = "RabbitMQInstance"

func init() {
	registry.Register(&registry.Registration{
		Name:     RabbitMQInstanceResource,
		Scope:    stackit.ProjectScope,
		Resource: &RabbitMQInstance{},
		Lister:   &RabbitMQInstanceLister{},
	})
}

type RabbitMQInstance struct {
	*BaseResource `property:",inline"`
	ID            string
	Name          string
}

func (r *RabbitMQInstance) Remove(_ context.Context) error {
	return fmt.Errorf("RabbitMQInstance.Remove not yet implemented")
}
func (r *RabbitMQInstance) Properties() types.Properties { return PropsFromStruct(r) }
func (r *RabbitMQInstance) String() string               { return r.Name }

type RabbitMQInstanceLister struct{}

func (l *RabbitMQInstanceLister) List(_ context.Context, _ any) ([]resource.Resource, error) {
	return []resource.Resource{}, nil
}
