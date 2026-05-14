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

const NetworkInterfaceResource = "NetworkInterface"

func init() {
	registry.Register(&registry.Registration{
		Name:     NetworkInterfaceResource,
		Scope:    stackit.ProjectScope,
		Resource: &NetworkInterface{},
		Lister:   &NetworkInterfaceLister{},
		// Server delete detaches NICs; deleting Network requires NICs gone first.
		DependsOn: []string{ComputeServerResource},
	})
}

// NetworkInterface is a STACKIT IaaS NIC. NICs live under a Network, so the
// API path includes both networkId and nicId.
//
// API: github.com/stackitcloud/stackit-sdk-go/services/iaas/v2api
//
// Endpoints used:
//   - GET    /v2/projects/{projectId}/regions/{region}/networks/{networkId}/nics
//   - DELETE /v2/projects/{projectId}/regions/{region}/networks/{networkId}/nics/{nicId}
type NetworkInterface struct {
	*BaseResource `property:",inline"`

	api iaasv2.DefaultAPI

	ID        string
	Name      string
	NetworkID string
	Device    string
	IPv4      string
	MAC       string
	Status    string
	Labels    map[string]string
}

func (r *NetworkInterface) Remove(ctx context.Context) error {
	if r.api == nil {
		return fmt.Errorf("NetworkInterface.Remove: api client not set")
	}
	return r.api.DeleteNic(ctx, r.ProjectID, r.Region, r.NetworkID, r.ID).Execute()
}

func (r *NetworkInterface) Properties() types.Properties { return PropsFromStruct(r) }
func (r *NetworkInterface) String() string               { return r.Name }

type NetworkInterfaceLister struct{}

func (l *NetworkInterfaceLister) List(ctx context.Context, o any) ([]resource.Resource, error) {
	opts := o.(*stackit.ListerOpts)
	logger := logrus.WithFields(logrus.Fields{
		"resource": NetworkInterfaceResource,
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

	logger.Trace("listing networks for NIC discovery")
	netResp, err := client.DefaultAPI.ListNetworks(ctx, opts.ProjectID, opts.Region).Execute()
	if err != nil {
		return nil, fmt.Errorf("list networks: %w", err)
	}

	var out []resource.Resource
	for _, n := range netResp.GetItems() {
		netID, ok := n.GetIdOk()
		if !ok || netID == nil {
			continue
		}

		nicResp, err := client.DefaultAPI.ListNics(ctx, opts.ProjectID, opts.Region, *netID).Execute()
		if err != nil {
			return nil, fmt.Errorf("list nics for network %s: %w", *netID, err)
		}

		for _, nic := range nicResp.GetItems() {
			id, ok := nic.GetIdOk()
			if !ok || id == nil {
				continue
			}
			name, _ := nic.GetNameOk()
			dev, _ := nic.GetDeviceOk()
			ipv4, _ := nic.GetIpv4Ok()
			mac, _ := nic.GetMacOk()
			status, _ := nic.GetStatusOk()

			labels := map[string]string{}
			if raw, ok := nic.GetLabelsOk(); ok {
				for k, val := range raw {
					if vs, ok := val.(string); ok {
						labels[k] = vs
					}
				}
			}

			out = append(out, &NetworkInterface{
				BaseResource: &BaseResource{
					OrganizationID: opts.OrganizationID,
					ProjectID:      opts.ProjectID,
					Region:         opts.Region,
				},
				api:       client.DefaultAPI,
				ID:        *id,
				Name:      stringDeref(name),
				NetworkID: *netID,
				Device:    stringDeref(dev),
				IPv4:      stringDeref(ipv4),
				MAC:       stringDeref(mac),
				Status:    stringDeref(status),
				Labels:    labels,
			})
		}
	}

	logger.WithField("count", len(out)).Trace("list complete")
	return out, nil
}
