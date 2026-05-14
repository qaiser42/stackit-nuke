package resources

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	iaasv2 "github.com/stackitcloud/stackit-sdk-go/services/iaas/v2api"

	"github.com/qaiser42/stackit-nuke/pkg/stackit"
)

const NetworkResource = "Network"

func init() {
	registry.Register(&registry.Registration{
		Name:     NetworkResource,
		Scope:    stackit.ProjectScope,
		Resource: &Network{},
		Lister:   &NetworkLister{},
		// STACKIT refuses to delete a network that still has attached NICs.
		DependsOn: []string{NetworkInterfaceResource},
	})
}

// Network is a STACKIT IaaS network.
//
// API: github.com/stackitcloud/stackit-sdk-go/services/iaas/v2api
//
// Endpoints used:
//   - GET    /v2/projects/{projectId}/regions/{region}/networks
//   - DELETE /v2/projects/{projectId}/regions/{region}/networks/{networkId}
type Network struct {
	*BaseResource `property:",inline"`

	api iaasv2.DefaultAPI

	ID        string
	Name      string
	Routed    bool
	Status    string
	CreatedAt *time.Time
	Labels    map[string]string
}

func (r *Network) Remove(ctx context.Context) error {
	if r.api == nil {
		return fmt.Errorf("Network.Remove: api client not set")
	}
	return r.api.DeleteNetwork(ctx, r.ProjectID, r.Region, r.ID).Execute()
}

func (r *Network) Properties() types.Properties { return PropsFromStruct(r) }
func (r *Network) String() string               { return r.Name }

type NetworkLister struct{}

func (l *NetworkLister) List(ctx context.Context, o any) ([]resource.Resource, error) {
	opts := o.(*stackit.ListerOpts)
	logger := logrus.WithFields(logrus.Fields{
		"resource": NetworkResource,
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

	logger.Trace("listing networks")
	resp, err := client.DefaultAPI.ListNetworks(ctx, opts.ProjectID, opts.Region).Execute()
	if err != nil {
		return nil, fmt.Errorf("list networks: %w", err)
	}

	out := make([]resource.Resource, 0, len(resp.GetItems()))
	for _, n := range resp.GetItems() {
		id, ok := n.GetIdOk()
		if !ok || id == nil {
			continue
		}
		name, _ := n.GetNameOk()
		routed, _ := n.GetRoutedOk()
		status, _ := n.GetStatusOk()
		createdAt, _ := n.GetCreatedAtOk()

		labels := map[string]string{}
		if raw, ok := n.GetLabelsOk(); ok {
			for k, val := range raw {
				if vs, ok := val.(string); ok {
					labels[k] = vs
				}
			}
		}

		var r bool
		if routed != nil {
			r = *routed
		}

		out = append(out, &Network{
			BaseResource: &BaseResource{
				OrganizationID: opts.OrganizationID,
				ProjectID:      opts.ProjectID,
				Region:         opts.Region,
			},
			api:       client.DefaultAPI,
			ID:        *id,
			Name:      stringDeref(name),
			Routed:    r,
			Status:    stringDeref(status),
			CreatedAt: createdAt,
			Labels:    labels,
		})
	}

	logger.WithField("count", len(out)).Trace("list complete")
	return out, nil
}
