package resources

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	postgresflexv2 "github.com/stackitcloud/stackit-sdk-go/services/postgresflex/v2api"

	"github.com/qaiser42/stackit-nuke/pkg/stackit"
)

const PostgresFlexInstanceResource = "PostgresFlexInstance"

func init() {
	registry.Register(&registry.Registration{
		Name:     PostgresFlexInstanceResource,
		Scope:    stackit.ProjectScope,
		Resource: &PostgresFlexInstance{},
		Lister:   &PostgresFlexInstanceLister{},
	})
}

// PostgresFlexInstance is a STACKIT PostgreSQL Flex instance.
//
// API: github.com/stackitcloud/stackit-sdk-go/services/postgresflex/v2api
//
// Endpoints used:
//   - GET    /v2/projects/{projectId}/regions/{region}/instances
//   - DELETE /v2/projects/{projectId}/regions/{region}/instances/{instanceId}
//   - DELETE /v2/projects/{projectId}/regions/{region}/instances/{instanceId}/force
//
// Deleting an instance is a soft delete: it stays listed with status
// "Deleted" during the retention window. Remove therefore force-deletes
// instances already in that state so a nuke run converges instead of
// re-deleting the same instance forever.
type PostgresFlexInstance struct {
	*BaseResource `property:",inline"`

	api postgresflexv2.DefaultAPI

	ID     string
	Name   string
	Status string
}

func (r *PostgresFlexInstance) Remove(ctx context.Context) error {
	if r.api == nil {
		return fmt.Errorf("PostgresFlexInstance.Remove: api client not set")
	}
	if r.Status == "Deleted" {
		return r.api.ForceDeleteInstance(ctx, r.ProjectID, r.Region, r.ID).Execute()
	}
	return r.api.DeleteInstance(ctx, r.ProjectID, r.Region, r.ID).Execute()
}

func (r *PostgresFlexInstance) Properties() types.Properties { return PropsFromStruct(r) }
func (r *PostgresFlexInstance) String() string               { return r.Name }

type PostgresFlexInstanceLister struct{}

func (l *PostgresFlexInstanceLister) List(ctx context.Context, o any) ([]resource.Resource, error) {
	opts := o.(*stackit.ListerOpts)
	logger := logrus.WithFields(logrus.Fields{
		"resource": PostgresFlexInstanceResource,
		"project":  opts.ProjectID,
		"region":   opts.Region,
	})

	if opts.Credentials == nil {
		return nil, fmt.Errorf("missing STACKIT credentials")
	}

	client, err := postgresflexv2.NewAPIClient(stackitConfigOpts(opts)...)
	if err != nil {
		return nil, fmt.Errorf("build postgresflex client: %w", err)
	}

	logger.Trace("listing postgresflex instances")
	resp, err := client.DefaultAPI.ListInstances(ctx, opts.ProjectID, opts.Region).Execute()
	if err != nil {
		return nil, fmt.Errorf("list postgresflex instances: %w", err)
	}

	out := make([]resource.Resource, 0, len(resp.GetItems()))
	for _, i := range resp.GetItems() {
		id, ok := i.GetIdOk()
		if !ok || id == nil {
			continue
		}
		name, _ := i.GetNameOk()
		status, _ := i.GetStatusOk()

		out = append(out, &PostgresFlexInstance{
			BaseResource: &BaseResource{
				OrganizationID: opts.OrganizationID,
				ProjectID:      opts.ProjectID,
				Region:         opts.Region,
			},
			api:    client.DefaultAPI,
			ID:     *id,
			Name:   stringDeref(name),
			Status: stringDeref(status),
		})
	}

	logger.WithField("count", len(out)).Trace("list complete")
	return out, nil
}
