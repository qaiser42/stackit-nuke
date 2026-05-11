package resources

import (
	"context"
	"fmt"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/qaiser42/stackit-nuke/pkg/stackit"
)

const PostgresFlexInstanceResource = "PostgresFlexInstance"

func init() {
	registry.Register(&registry.Registration{
		Name:     PostgresFlexInstanceResource,
		Scope:    stackit.ProjectScope,
		Resource: &PostgresFlexInstance{},
		Lister:   &PostgresFlexInstanceLister{},
	})
}

type PostgresFlexInstance struct {
	*BaseResource `property:",inline"`
	ID            string
	Name          string
}

func (r *PostgresFlexInstance) Remove(_ context.Context) error {
	return fmt.Errorf("PostgresFlexInstance.Remove not yet implemented")
}
func (r *PostgresFlexInstance) Properties() types.Properties { return PropsFromStruct(r) }
func (r *PostgresFlexInstance) String() string                { return r.Name }

type PostgresFlexInstanceLister struct{}

func (l *PostgresFlexInstanceLister) List(_ context.Context, _ any) ([]resource.Resource, error) {
	return []resource.Resource{}, nil
}
