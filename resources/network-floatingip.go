package resources

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	iaasv2 "github.com/stackitcloud/stackit-sdk-go/services/iaas/v2api"

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

// FloatingIP is a STACKIT IaaS public IP.
//
// API: github.com/stackitcloud/stackit-sdk-go/services/iaas/v2api
//
// Endpoints used:
//   - GET    /v2/projects/{projectId}/regions/{region}/public-ips
//   - DELETE /v2/projects/{projectId}/regions/{region}/public-ips/{publicIpId}
type FloatingIP struct {
	*BaseResource `property:",inline"`

	api iaasv2.DefaultAPI

	ID               string
	IP               string
	NetworkInterface string
	Labels           map[string]string
}

func (r *FloatingIP) Remove(ctx context.Context) error {
	if r.api == nil {
		return fmt.Errorf("FloatingIP.Remove: api client not set")
	}
	return r.api.DeletePublicIP(ctx, r.ProjectID, r.Region, r.ID).Execute()
}

func (r *FloatingIP) Properties() types.Properties { return PropsFromStruct(r) }
func (r *FloatingIP) String() string               { return r.IP }

type FloatingIPLister struct{}

func (l *FloatingIPLister) List(ctx context.Context, o any) ([]resource.Resource, error) {
	opts := o.(*stackit.ListerOpts)
	logger := logrus.WithFields(logrus.Fields{
		"resource": FloatingIPResource,
		"project":  opts.ProjectID,
		"region":   opts.Region,
	})

	if opts.Credentials == nil {
		return nil, fmt.Errorf("missing STACKIT credentials")
	}

	client, err := iaasv2.NewAPIClient(stackitConfigOpts(opts)...)
	if err != nil {
		return nil, fmt.Errorf("build iaas client: %w", err)
	}

	logger.Trace("listing public ips")
	resp, err := client.DefaultAPI.ListPublicIPs(ctx, opts.ProjectID, opts.Region).Execute()
	if err != nil {
		return nil, fmt.Errorf("list public ips: %w", err)
	}

	out := make([]resource.Resource, 0, len(resp.GetItems()))
	for _, p := range resp.GetItems() {
		id, ok := p.GetIdOk()
		if !ok || id == nil {
			continue
		}
		ip, _ := p.GetIpOk()
		nic, _ := p.GetNetworkInterfaceOk()

		labels := map[string]string{}
		if raw, ok := p.GetLabelsOk(); ok {
			for k, v := range raw {
				if vs, ok := v.(string); ok {
					labels[k] = vs
				}
			}
		}

		out = append(out, &FloatingIP{
			BaseResource: &BaseResource{
				OrganizationID: opts.OrganizationID,
				ProjectID:      opts.ProjectID,
				Region:         opts.Region,
			},
			api:              client.DefaultAPI,
			ID:               *id,
			IP:               stringDeref(ip),
			NetworkInterface: stringDeref(nic),
			Labels:           labels,
		})
	}

	logger.WithField("count", len(out)).Trace("list complete")
	return out, nil
}
