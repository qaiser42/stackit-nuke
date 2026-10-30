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

const ComputeSnapshotResource = "ComputeSnapshot"

func init() {
	registry.Register(&registry.Registration{
		Name:     ComputeSnapshotResource,
		Scope:    stackit.ProjectScope,
		Resource: &ComputeSnapshot{},
		Lister:   &ComputeSnapshotLister{},
	})
}

// ComputeSnapshot is a STACKIT IaaS volume snapshot.
//
// API: github.com/stackitcloud/stackit-sdk-go/services/iaas/v2api
//
// Endpoints used:
//   - GET    /v2/projects/{projectId}/regions/{region}/snapshots
//   - DELETE /v2/projects/{projectId}/regions/{region}/snapshots/{snapshotId}
//
// Deletion is asynchronous: the snapshot stays listed with status DELETING
// until it is gone, which is exactly the wait behavior libnuke expects.
// Remove skips the DELETE call for snapshots already in that state, and the
// Lister drops snapshots with status DELETED so a finished delete can never
// keep libnuke waiting.
type ComputeSnapshot struct {
	*BaseResource `property:",inline"`

	api iaasv2.DefaultAPI

	ID        string
	Name      string
	VolumeID  string
	SizeGB    int64
	Status    string
	CreatedAt *time.Time
	Labels    map[string]string
}

const (
	computeSnapshotStatusDeleting = "DELETING"
	computeSnapshotStatusDeleted  = "DELETED"
)

func (r *ComputeSnapshot) Remove(ctx context.Context) error {
	if r.api == nil {
		return fmt.Errorf("ComputeSnapshot.Remove: api client not set")
	}
	if r.Status == computeSnapshotStatusDeleting {
		return nil
	}
	if err := r.api.DeleteSnapshot(ctx, r.ProjectID, r.Region, r.ID).Execute(); err != nil {
		return fmt.Errorf("delete snapshot: %w", err)
	}
	return nil
}

func (r *ComputeSnapshot) Properties() types.Properties { return PropsFromStruct(r) }
func (r *ComputeSnapshot) String() string               { return r.Name }

type ComputeSnapshotLister struct{}

func (l *ComputeSnapshotLister) List(ctx context.Context, o any) ([]resource.Resource, error) {
	opts := o.(*stackit.ListerOpts)
	logger := logrus.WithFields(logrus.Fields{
		"resource": ComputeSnapshotResource,
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

	logger.Trace("listing snapshots")
	resp, err := client.DefaultAPI.ListSnapshotsInProject(ctx, opts.ProjectID, opts.Region).Execute()
	if err != nil {
		return nil, fmt.Errorf("list snapshots: %w", err)
	}

	out := make([]resource.Resource, 0, len(resp.GetItems()))
	for _, s := range resp.GetItems() {
		id, ok := s.GetIdOk()
		if !ok || id == nil {
			continue
		}
		status, _ := s.GetStatusOk()
		if stringDeref(status) == computeSnapshotStatusDeleted {
			continue
		}
		name, _ := s.GetNameOk()
		size, _ := s.GetSizeOk()
		createdAt, _ := s.GetCreatedAtOk()

		labels := map[string]string{}
		if raw, ok := s.GetLabelsOk(); ok {
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

		out = append(out, &ComputeSnapshot{
			BaseResource: &BaseResource{
				OrganizationID: opts.OrganizationID,
				ProjectID:      opts.ProjectID,
				Region:         opts.Region,
			},
			api:       client.DefaultAPI,
			ID:        *id,
			Name:      stringDeref(name),
			VolumeID:  s.GetVolumeId(),
			SizeGB:    sz,
			Status:    stringDeref(status),
			CreatedAt: createdAt,
			Labels:    labels,
		})
	}

	logger.WithField("count", len(out)).Trace("list complete")
	return out, nil
}
