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

const ComputeVolumeResource = "ComputeVolume"

func init() {
	registry.Register(&registry.Registration{
		Name:      ComputeVolumeResource,
		Scope:     stackit.ProjectScope,
		Resource:  &ComputeVolume{},
		Lister:    &ComputeVolumeLister{},
		DependsOn: []string{ComputeServerResource},
	})
}

// ComputeVolume is a STACKIT IaaS block volume.
//
// API: github.com/stackitcloud/stackit-sdk-go/services/iaas/v2api
//
// Endpoints used:
//   - GET    /v2/projects/{projectId}/regions/{region}/volumes
//   - DELETE /v2/projects/{projectId}/regions/{region}/volumes/{volumeId}
type ComputeVolume struct {
	*BaseResource `property:",inline"`

	api iaasv2.DefaultAPI

	ID               string
	Name             string
	AvailabilityZone string
	SizeGB           int64
	Bootable         bool
	Status           string
	ServerID         string
	CreatedAt        *time.Time
	Labels           map[string]string
}

func (r *ComputeVolume) Remove(ctx context.Context) error {
	if r.api == nil {
		return fmt.Errorf("ComputeVolume.Remove: api client not set")
	}
	return r.api.DeleteVolume(ctx, r.ProjectID, r.Region, r.ID).Execute()
}

func (r *ComputeVolume) Properties() types.Properties { return PropsFromStruct(r) }
func (r *ComputeVolume) String() string               { return r.Name }

type ComputeVolumeLister struct{}

func (l *ComputeVolumeLister) List(ctx context.Context, o any) ([]resource.Resource, error) {
	opts := o.(*stackit.ListerOpts)
	logger := logrus.WithFields(logrus.Fields{
		"resource": ComputeVolumeResource,
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

	logger.Trace("listing volumes")
	resp, err := client.DefaultAPI.ListVolumes(ctx, opts.ProjectID, opts.Region).Execute()
	if err != nil {
		return nil, fmt.Errorf("list volumes: %w", err)
	}

	out := make([]resource.Resource, 0, len(resp.GetItems()))
	for _, v := range resp.GetItems() {
		id, ok := v.GetIdOk()
		if !ok || id == nil {
			continue
		}
		name, _ := v.GetNameOk()
		az, _ := v.GetAvailabilityZoneOk()
		size, _ := v.GetSizeOk()
		boot, _ := v.GetBootableOk()
		status, _ := v.GetStatusOk()
		serverID, _ := v.GetServerIdOk()
		createdAt, _ := v.GetCreatedAtOk()

		labels := map[string]string{}
		if raw, ok := v.GetLabelsOk(); ok {
			for k, val := range raw {
				if vs, ok := val.(string); ok {
					labels[k] = vs
				}
			}
		}

		var sz int64
		if size != nil {
			sz = *size
		}
		var b bool
		if boot != nil {
			b = *boot
		}

		out = append(out, &ComputeVolume{
			BaseResource: &BaseResource{
				OrganizationID: opts.OrganizationID,
				ProjectID:      opts.ProjectID,
				Region:         opts.Region,
			},
			api:              client.DefaultAPI,
			ID:               *id,
			Name:             stringDeref(name),
			AvailabilityZone: stringDeref(az),
			SizeGB:           sz,
			Bootable:         b,
			Status:           stringDeref(status),
			ServerID:         stringDeref(serverID),
			CreatedAt:        createdAt,
			Labels:           labels,
		})
	}

	logger.WithField("count", len(out)).Trace("list complete")
	return out, nil
}
