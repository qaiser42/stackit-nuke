package resources

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	mariadbv2 "github.com/stackitcloud/stackit-sdk-go/services/mariadb/v2api"

	"github.com/qaiser42/stackit-nuke/pkg/stackit"
)

const MariaDBInstanceResource = "MariaDBInstance"

func init() {
	registry.Register(&registry.Registration{
		Name:     MariaDBInstanceResource,
		Scope:    stackit.ProjectScope,
		Resource: &MariaDBInstance{},
		Lister:   &MariaDBInstanceLister{},
	})
}

// MariaDBInstance is a STACKIT MariaDB instance.
//
// API: github.com/stackitcloud/stackit-sdk-go/services/mariadb/v2api
//
// Endpoints used:
//   - GET    /v2/projects/{projectId}/regions/{regionId}/instances
//   - DELETE /v2/projects/{projectId}/regions/{regionId}/instances/{instanceId}
//
// Deletion is asynchronous: the instance stays listed with status
// "deleting" until it is gone, which is exactly the wait behavior libnuke
// expects. Remove skips the DELETE call for instances already in that
// state (e.g. left over from an interrupted run). The Lister drops
// instances whose last operation is a succeeded delete, so a finished
// delete can never keep libnuke waiting.
type MariaDBInstance struct {
	*BaseResource `property:",inline"`

	api mariadbv2.DefaultAPI

	ID      string
	Name    string
	Plan    string
	Version string
	Status  string
}

func (r *MariaDBInstance) Remove(ctx context.Context) error {
	if r.api == nil {
		return fmt.Errorf("MariaDBInstance.Remove: api client not set")
	}
	if r.Status == string(mariadbv2.INSTANCESTATUS_DELETING) {
		return nil
	}
	if err := r.api.DeleteInstance(ctx, r.ProjectID, r.Region, r.ID).Execute(); err != nil {
		return fmt.Errorf("delete mariadb instance: %w", err)
	}
	return nil
}

func (r *MariaDBInstance) Properties() types.Properties { return PropsFromStruct(r) }
func (r *MariaDBInstance) String() string               { return r.Name }

type MariaDBInstanceLister struct{}

func (l *MariaDBInstanceLister) List(ctx context.Context, o any) ([]resource.Resource, error) {
	opts := o.(*stackit.ListerOpts)
	logger := logrus.WithFields(logrus.Fields{
		"resource": MariaDBInstanceResource,
		"project":  opts.ProjectID,
		"region":   opts.Region,
	})

	if opts.Credentials == nil {
		return nil, fmt.Errorf("missing STACKIT credentials")
	}

	client, err := mariadbv2.NewAPIClient(stackitConfigOpts(opts)...)
	if err != nil {
		return nil, fmt.Errorf("build mariadb client: %w", err)
	}

	logger.Trace("listing mariadb instances")
	resp, err := client.DefaultAPI.ListInstances(ctx, opts.ProjectID, opts.Region).Execute()
	if err != nil {
		return nil, fmt.Errorf("list mariadb instances: %w", err)
	}

	instances := resp.GetInstances()
	out := make([]resource.Resource, 0, len(instances))
	for _, i := range instances {
		id, ok := i.GetInstanceIdOk()
		if !ok || id == nil {
			continue
		}
		op := i.GetLastOperation()
		if op.GetType() == mariadbv2.INSTANCELASTOPERATIONTYPE_DELETE &&
			op.GetState() == mariadbv2.INSTANCELASTOPERATIONSTATE_SUCCEEDED {
			continue
		}

		out = append(out, &MariaDBInstance{
			BaseResource: &BaseResource{
				OrganizationID: opts.OrganizationID,
				ProjectID:      opts.ProjectID,
				Region:         opts.Region,
			},
			api:     client.DefaultAPI,
			ID:      *id,
			Name:    i.GetName(),
			Plan:    i.GetPlanName(),
			Version: i.GetOfferingVersion(),
			Status:  string(i.GetStatus()),
		})
	}

	logger.WithField("count", len(out)).Trace("list complete")
	return out, nil
}
