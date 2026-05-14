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

const SecurityGroupResource = "SecurityGroup"

func init() {
	registry.Register(&registry.Registration{
		Name:     SecurityGroupResource,
		Scope:    stackit.ProjectScope,
		Resource: &SecurityGroup{},
		Lister:   &SecurityGroupLister{},
		// NICs reference SGs — drop NICs first so SG delete isn't blocked.
		DependsOn: []string{NetworkInterfaceResource},
	})
}

// SecurityGroup is a STACKIT IaaS security group.
//
// Note: STACKIT auto-creates a "default" SG per project. It is generally
// protected and will fail to delete; filter it out in config via
// `filters.SecurityGroup: [{property: Name, value: default}]`.
//
// API: github.com/stackitcloud/stackit-sdk-go/services/iaas/v2api
//
// Endpoints used:
//   - GET    /v2/projects/{projectId}/regions/{region}/security-groups
//   - DELETE /v2/projects/{projectId}/regions/{region}/security-groups/{securityGroupId}
type SecurityGroup struct {
	*BaseResource `property:",inline"`

	api iaasv2.DefaultAPI

	ID          string
	Name        string
	Description string
	Stateful    bool
	CreatedAt   *time.Time
	Labels      map[string]string
}

func (r *SecurityGroup) Remove(ctx context.Context) error {
	if r.api == nil {
		return fmt.Errorf("SecurityGroup.Remove: api client not set")
	}
	return r.api.DeleteSecurityGroup(ctx, r.ProjectID, r.Region, r.ID).Execute()
}

func (r *SecurityGroup) Properties() types.Properties { return PropsFromStruct(r) }
func (r *SecurityGroup) String() string               { return r.Name }

type SecurityGroupLister struct{}

func (l *SecurityGroupLister) List(ctx context.Context, o any) ([]resource.Resource, error) {
	opts := o.(*stackit.ListerOpts)
	logger := logrus.WithFields(logrus.Fields{
		"resource": SecurityGroupResource,
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

	logger.Trace("listing security groups")
	resp, err := client.DefaultAPI.ListSecurityGroups(ctx, opts.ProjectID, opts.Region).Execute()
	if err != nil {
		return nil, fmt.Errorf("list security groups: %w", err)
	}

	out := make([]resource.Resource, 0, len(resp.GetItems()))
	for _, sg := range resp.GetItems() {
		id, ok := sg.GetIdOk()
		if !ok || id == nil {
			continue
		}
		name, _ := sg.GetNameOk()
		desc, _ := sg.GetDescriptionOk()
		stateful, _ := sg.GetStatefulOk()
		createdAt, _ := sg.GetCreatedAtOk()

		labels := map[string]string{}
		if raw, ok := sg.GetLabelsOk(); ok {
			for k, val := range raw {
				if vs, ok := val.(string); ok {
					labels[k] = vs
				}
			}
		}

		var s bool
		if stateful != nil {
			s = *stateful
		}

		out = append(out, &SecurityGroup{
			BaseResource: &BaseResource{
				OrganizationID: opts.OrganizationID,
				ProjectID:      opts.ProjectID,
				Region:         opts.Region,
			},
			api:         client.DefaultAPI,
			ID:          *id,
			Name:        stringDeref(name),
			Description: stringDeref(desc),
			Stateful:    s,
			CreatedAt:   createdAt,
			Labels:      labels,
		})
	}

	logger.WithField("count", len(out)).Trace("list complete")
	return out, nil
}
