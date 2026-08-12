package resources

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	gitv1beta "github.com/stackitcloud/stackit-sdk-go/services/git/v1betaapi"

	"github.com/qaiser42/stackit-nuke/pkg/stackit"
)

const GitInstanceResource = "GitInstance"

func init() {
	registry.Register(&registry.Registration{
		Name:     GitInstanceResource,
		Scope:    stackit.ProjectScope,
		Resource: &GitInstance{},
		Lister:   &GitInstanceLister{},
	})
}

// GitInstance is a STACKIT Git instance.
//
// API: github.com/stackitcloud/stackit-sdk-go/services/git/v1betaapi
//
// Endpoints used:
//   - GET    /v1beta/projects/{projectId}/instances
//   - DELETE /v1beta/projects/{projectId}/instances/{instanceId}
//
// STACKIT Git is a global service: instances are not bound to a region
// (the API host's region variable defaults to "global" and no endpoint
// takes a region), but listers run once per (project, region). The Lister
// therefore only returns instances for the first configured region, so
// each instance shows up exactly once per run.
//
// Deletion is asynchronous: the instance stays listed in state "Deleting"
// until it is gone, which is exactly the wait behavior libnuke expects.
type GitInstance struct {
	*BaseResource `property:",inline"`

	api gitv1beta.DefaultAPI

	ID      string
	Name    string
	URL     string
	Flavor  string
	Version string
	State   string
}

func (r *GitInstance) Remove(ctx context.Context) error {
	if r.api == nil {
		return fmt.Errorf("GitInstance.Remove: api client not set")
	}
	if err := r.api.DeleteInstance(ctx, r.ProjectID, r.ID).Execute(); err != nil {
		return fmt.Errorf("delete git instance: %w", err)
	}
	return nil
}

func (r *GitInstance) Properties() types.Properties { return PropsFromStruct(r) }
func (r *GitInstance) String() string               { return r.Name }

type GitInstanceLister struct{}

func (l *GitInstanceLister) List(ctx context.Context, o any) ([]resource.Resource, error) {
	opts := o.(*stackit.ListerOpts)
	logger := logrus.WithFields(logrus.Fields{
		"resource": GitInstanceResource,
		"project":  opts.ProjectID,
		"region":   opts.Region,
	})

	if opts.Credentials == nil {
		return nil, fmt.Errorf("missing STACKIT credentials")
	}

	// Git instances are global; only list them in the first region's scan.
	if len(opts.Regions) > 0 && opts.Region != opts.Regions[0] {
		return []resource.Resource{}, nil
	}

	client, err := gitv1beta.NewAPIClient(stackitConfigOpts(opts)...)
	if err != nil {
		return nil, fmt.Errorf("build git client: %w", err)
	}

	logger.Trace("listing git instances")
	resp, err := client.DefaultAPI.ListInstances(ctx, opts.ProjectID).Execute()
	if err != nil {
		return nil, fmt.Errorf("list git instances: %w", err)
	}

	instances := resp.GetInstances()
	out := make([]resource.Resource, 0, len(instances))
	for _, i := range instances {
		out = append(out, &GitInstance{
			BaseResource: &BaseResource{
				OrganizationID: opts.OrganizationID,
				ProjectID:      opts.ProjectID,
				Region:         opts.Region,
			},
			api:     client.DefaultAPI,
			ID:      i.GetId(),
			Name:    i.GetName(),
			URL:     i.GetUrl(),
			Flavor:  i.GetFlavor(),
			Version: i.GetVersion(),
			State:   string(i.GetState()),
		})
	}

	logger.WithField("count", len(out)).Trace("list complete")
	return out, nil
}
