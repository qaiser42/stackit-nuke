package resources

import (
	"context"
	"fmt"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/qaiser42/stackit-nuke/pkg/stackit"
)

const RouterResource = "Router"

func init() {
	registry.Register(&registry.Registration{
		Name:     RouterResource,
		Scope:    stackit.ProjectScope,
		Resource: &Router{},
		Lister:   &RouterLister{},
	})
}

type Router struct {
	*BaseResource `property:",inline"`
	ID            string
	Name          string
}

func (r *Router) Remove(_ context.Context) error {
	return fmt.Errorf("Router.Remove not yet implemented")
}
func (r *Router) Properties() types.Properties { return PropsFromStruct(r) }
func (r *Router) String() string                { return r.Name }

type RouterLister struct{}

func (l *RouterLister) List(_ context.Context, _ any) ([]resource.Resource, error) {
	return []resource.Resource{}, nil
}
