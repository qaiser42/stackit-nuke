package resources

import (
	"context"
	"fmt"
	"time"

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
//   - GET    /v2/projects/{projectId}/regions/{region}/instances/{instanceId}
//   - DELETE /v2/projects/{projectId}/regions/{region}/instances/{instanceId}
//   - DELETE /v2/projects/{projectId}/regions/{region}/instances/{instanceId}/force
//
// Deleting an instance is a soft delete: it stays listed with status
// "Deleted" during the retention window, so libnuke would wait on it
// forever (Remove is only invoked once per item per run). Remove therefore
// soft-deletes, polls until the instance reaches "Deleted", then
// force-deletes to purge it within the same run.
type PostgresFlexInstance struct {
	*BaseResource `property:",inline"`

	api postgresflexv2.DefaultAPI

	ID     string
	Name   string
	Status string
}

const postgresFlexStatusDeleted = "Deleted"

func (r *PostgresFlexInstance) Remove(ctx context.Context) error {
	if r.api == nil {
		return fmt.Errorf("PostgresFlexInstance.Remove: api client not set")
	}
	if r.Status != postgresFlexStatusDeleted {
		if err := r.api.DeleteInstance(ctx, r.ProjectID, r.Region, r.ID).Execute(); err != nil {
			return fmt.Errorf("delete instance: %w", err)
		}
		if err := r.waitUntilDeleted(ctx); err != nil {
			return err
		}
	}
	if err := r.api.ForceDeleteInstance(ctx, r.ProjectID, r.Region, r.ID).Execute(); err != nil {
		return fmt.Errorf("force delete instance: %w", err)
	}
	return nil
}

// waitUntilDeleted polls the instance until the soft delete has landed
// (status "Deleted"); force delete is rejected before that.
func (r *PostgresFlexInstance) waitUntilDeleted(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		resp, err := r.api.GetInstance(ctx, r.ProjectID, r.Region, r.ID).Execute()
		if err != nil {
			return fmt.Errorf("get instance while waiting for soft delete: %w", err)
		}
		if item, ok := resp.GetItemOk(); ok && item.GetStatus() == postgresFlexStatusDeleted {
			return nil
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("timed out waiting for instance %s to reach %s state: %w", r.ID, postgresFlexStatusDeleted, ctx.Err())
		case <-ticker.C:
		}
	}
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
