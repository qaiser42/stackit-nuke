package resources

import (
	"context"
	"fmt"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/qaiser42/stackit-nuke/pkg/stackit"
)

const FloatingIPResource = "FloatingIP"

func init() {
	registry.Register(&registry.Registration{
		Name:     FloatingIPResource,
		Scope:    stackit.ProjectScope,
		Resource: &FloatingIP{},
		Lister:   &FloatingIPLister{},
	})
}

type FloatingIP struct {
	*BaseResource `property:",inline"`
	ID            string
	IP            string
}

func (r *FloatingIP) Remove(_ context.Context) error {
	return fmt.Errorf("FloatingIP.Remove not yet implemented")
}
func (r *FloatingIP) Properties() types.Properties { return PropsFromStruct(r) }
func (r *FloatingIP) String() string               { return r.IP }

type FloatingIPLister struct{}

func (l *FloatingIPLister) List(_ context.Context, _ any) ([]resource.Resource, error) {
	return []resource.Resource{}, nil
}
