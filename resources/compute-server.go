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

const ComputeServerResource = "ComputeServer"

func init() {
	registry.Register(&registry.Registration{
		Name:     ComputeServerResource,
		Scope:    stackit.ProjectScope,
		Resource: &ComputeServer{},
		Lister:   &ComputeServerLister{},
	})
}

// ComputeServer is a STACKIT IaaS server (a.k.a. virtual machine).
//
// API: github.com/stackitcloud/stackit-sdk-go/services/iaas/v2api
//
// Endpoints used:
//   - GET    /v2/projects/{projectId}/regions/{region}/servers
//   - DELETE /v2/projects/{projectId}/regions/{region}/servers/{serverId}
type ComputeServer struct {
	*BaseResource `property:",inline"`

	api iaasv2.DefaultAPI

	ID               string
	Name             string
	AvailabilityZone string
	CreatedAt        *time.Time
	Labels           map[string]string
}

func (r *ComputeServer) Remove(ctx context.Context) error {
	if r.api == nil {
		return fmt.Errorf("ComputeServer.Remove: api client not set")
	}
	return r.api.DeleteServer(ctx, r.ProjectID, r.Region, r.ID).Execute()
}

func (r *ComputeServer) Properties() types.Properties { return PropsFromStruct(r) }
func (r *ComputeServer) String() string                { return r.Name }

type ComputeServerLister struct{}

func (l *ComputeServerLister) List(ctx context.Context, o any) ([]resource.Resource, error) {
	opts := o.(*stackit.ListerOpts)
	logger := logrus.WithFields(logrus.Fields{
		"resource": ComputeServerResource,
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

	logger.Trace("listing servers")
	resp, err := client.DefaultAPI.ListServers(ctx, opts.ProjectID, opts.Region).Details(true).Execute()
	if err != nil {
		return nil, fmt.Errorf("list servers: %w", err)
	}

	out := make([]resource.Resource, 0, len(resp.GetItems()))
	for _, s := range resp.GetItems() {
		id, ok := s.GetIdOk()
		if !ok || id == nil {
			continue
		}
		name, _ := s.GetNameOk()
		az, _ := s.GetAvailabilityZoneOk()
		createdAt, _ := s.GetCreatedAtOk()

		labels := map[string]string{}
		if raw, ok := s.GetLabelsOk(); ok {
			for k, v := range raw {
				if vs, ok := v.(string); ok {
					labels[k] = vs
				}
			}
		}

		out = append(out, &ComputeServer{
			BaseResource: &BaseResource{
				OrganizationID: opts.OrganizationID,
				ProjectID:      opts.ProjectID,
				Region:         opts.Region,
			},
			api:              client.DefaultAPI,
			ID:               *id,
			Name:             stringDeref(name),
			AvailabilityZone: stringDeref(az),
			CreatedAt:        createdAt,
			Labels:           labels,
		})
	}

	logger.WithField("count", len(out)).Trace("list complete")
	return out, nil
}
