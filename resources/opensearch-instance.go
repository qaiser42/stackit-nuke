package resources

import (
	"context"
	"fmt"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/qaiser42/stackit-nuke/pkg/stackit"
)

const OpenSearchInstanceResource = "OpenSearchInstance"

func init() {
	registry.Register(&registry.Registration{
		Name:     OpenSearchInstanceResource,
		Scope:    stackit.ProjectScope,
		Resource: &OpenSearchInstance{},
		Lister:   &OpenSearchInstanceLister{},
	})
}

type OpenSearchInstance struct {
	*BaseResource `property:",inline"`
	ID            string
	Name          string
}

func (r *OpenSearchInstance) Remove(_ context.Context) error {
	return fmt.Errorf("OpenSearchInstance.Remove not yet implemented")
}
func (r *OpenSearchInstance) Properties() types.Properties { return PropsFromStruct(r) }
func (r *OpenSearchInstance) String() string               { return r.Name }

type OpenSearchInstanceLister struct{}

func (l *OpenSearchInstanceLister) List(_ context.Context, _ any) ([]resource.Resource, error) {
	return []resource.Resource{}, nil
}
