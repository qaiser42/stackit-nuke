package resources

import (
	"context"
	"fmt"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/qaiser42/stackit-nuke/pkg/stackit"
)

const DNSZoneResource = "DNSZone"

func init() {
	registry.Register(&registry.Registration{
		Name:     DNSZoneResource,
		Scope:    stackit.ProjectScope,
		Resource: &DNSZone{},
		Lister:   &DNSZoneLister{},
	})
}

type DNSZone struct {
	*BaseResource `property:",inline"`
	ID            string
	Name          string
}

func (r *DNSZone) Remove(_ context.Context) error {
	return fmt.Errorf("DNSZone.Remove not yet implemented")
}
func (r *DNSZone) Properties() types.Properties { return PropsFromStruct(r) }
func (r *DNSZone) String() string                { return r.Name }

type DNSZoneLister struct{}

func (l *DNSZoneLister) List(_ context.Context, _ any) ([]resource.Resource, error) {
	return []resource.Resource{}, nil
}
