package resources

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	dnsv1 "github.com/stackitcloud/stackit-sdk-go/services/dns/v1api"

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

// DNSZone is a STACKIT DNS zone.
//
// API: github.com/stackitcloud/stackit-sdk-go/services/dns/v1api
//
// Endpoints used:
//   - GET    /v1/projects/{projectId}/zones
//   - DELETE /v1/projects/{projectId}/zones/{zoneId}
//
// DNS is a global service: zones are not bound to a region, but listers run
// once per (project, region). The Lister therefore only returns zones for
// the first configured region, so each zone shows up exactly once per run.
//
// Deleting a zone is a soft delete: the zone stays listed with state
// DELETE_SUCCEEDED during the retention window. The Lister skips those, so
// libnuke sees the zone as gone once the delete has landed.
type DNSZone struct {
	*BaseResource `property:",inline"`

	api dnsv1.DefaultAPI

	ID      string
	Name    string
	DNSName string
	State   string
}

func (r *DNSZone) Remove(ctx context.Context) error {
	if r.api == nil {
		return fmt.Errorf("DNSZone.Remove: api client not set")
	}
	if _, err := r.api.DeleteZone(ctx, r.ProjectID, r.ID).Execute(); err != nil {
		return fmt.Errorf("delete zone: %w", err)
	}
	return nil
}

func (r *DNSZone) Properties() types.Properties { return PropsFromStruct(r) }
func (r *DNSZone) String() string               { return r.DNSName }

type DNSZoneLister struct{}

func (l *DNSZoneLister) List(ctx context.Context, o any) ([]resource.Resource, error) {
	opts := o.(*stackit.ListerOpts)
	logger := logrus.WithFields(logrus.Fields{
		"resource": DNSZoneResource,
		"project":  opts.ProjectID,
		"region":   opts.Region,
	})

	if opts.Credentials == nil {
		return nil, fmt.Errorf("missing STACKIT credentials")
	}

	// Zones are global; only list them in the first region's scan.
	if len(opts.Regions) > 0 && opts.Region != opts.Regions[0] {
		return []resource.Resource{}, nil
	}

	client, err := dnsv1.NewAPIClient(stackitConfigOpts(opts)...)
	if err != nil {
		return nil, fmt.Errorf("build dns client: %w", err)
	}

	logger.Trace("listing dns zones")
	out := []resource.Resource{}
	for page := int32(1); ; page++ {
		resp, err := client.DefaultAPI.ListZones(ctx, opts.ProjectID).
			Page(page).
			StateNeq(dnsv1.LISTZONESSTATENEQPARAMETER_DELETE_SUCCEEDED).
			Execute()
		if err != nil {
			return nil, fmt.Errorf("list dns zones: %w", err)
		}

		for _, z := range resp.GetZones() {
			out = append(out, &DNSZone{
				BaseResource: &BaseResource{
					OrganizationID: opts.OrganizationID,
					ProjectID:      opts.ProjectID,
					Region:         opts.Region,
				},
				api:     client.DefaultAPI,
				ID:      z.GetId(),
				Name:    z.GetName(),
				DNSName: z.GetDnsName(),
				State:   string(z.GetState()),
			})
		}

		if page >= resp.GetTotalPages() {
			break
		}
	}

	logger.WithField("count", len(out)).Trace("list complete")
	return out, nil
}
