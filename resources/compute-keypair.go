package resources

import (
	"context"
	"fmt"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/qaiser42/stackit-nuke/pkg/stackit"
)

const ComputeKeypairResource = "ComputeKeypair"

func init() {
	registry.Register(&registry.Registration{
		Name:     ComputeKeypairResource,
		Scope:    stackit.ProjectScope,
		Resource: &ComputeKeypair{},
		Lister:   &ComputeKeypairLister{},
	})
}

type ComputeKeypair struct {
	*BaseResource `property:",inline"`
	Name          string
}

func (r *ComputeKeypair) Remove(_ context.Context) error {
	return fmt.Errorf("ComputeKeypair.Remove not yet implemented")
}
func (r *ComputeKeypair) Properties() types.Properties { return PropsFromStruct(r) }
func (r *ComputeKeypair) String() string                { return r.Name }

type ComputeKeypairLister struct{}

func (l *ComputeKeypairLister) List(_ context.Context, _ any) ([]resource.Resource, error) {
	return []resource.Resource{}, nil
}
