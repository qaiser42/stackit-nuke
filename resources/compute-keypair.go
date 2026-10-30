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

const ComputeKeypairResource = "ComputeKeypair"

func init() {
	registry.Register(&registry.Registration{
		Name:     ComputeKeypairResource,
		Scope:    stackit.ProjectScope,
		Resource: &ComputeKeypair{},
		Lister:   &ComputeKeypairLister{},
	})
}

// ComputeKeypair is a STACKIT IaaS SSH keypair.
//
// API: github.com/stackitcloud/stackit-sdk-go/services/iaas/v2api
//
// Endpoints used:
//   - GET    /v2/keypairs
//   - DELETE /v2/keypairs/{keypairName}
//
// Keypairs are neither regional nor project-scoped: they belong to the
// authenticated identity (user or service account). Listers run once per
// (project, region); the Lister only returns keypairs during the first
// configured project's first region scan, so each keypair shows up exactly
// once per run.
type ComputeKeypair struct {
	*BaseResource `property:",inline"`

	api iaasv2.DefaultAPI

	Name        string
	Fingerprint string
	CreatedAt   *time.Time
	Labels      map[string]string
}

func (r *ComputeKeypair) Remove(ctx context.Context) error {
	if r.api == nil {
		return fmt.Errorf("ComputeKeypair.Remove: api client not set")
	}
	if err := r.api.DeleteKeyPair(ctx, r.Name).Execute(); err != nil {
		return fmt.Errorf("delete keypair: %w", err)
	}
	return nil
}

func (r *ComputeKeypair) Properties() types.Properties { return PropsFromStruct(r) }
func (r *ComputeKeypair) String() string               { return r.Name }

type ComputeKeypairLister struct{}

func (l *ComputeKeypairLister) List(ctx context.Context, o any) ([]resource.Resource, error) {
	opts := o.(*stackit.ListerOpts)
	logger := logrus.WithFields(logrus.Fields{
		"resource": ComputeKeypairResource,
		"project":  opts.ProjectID,
		"region":   opts.Region,
	})

	if opts.Credentials == nil {
		return nil, fmt.Errorf("missing STACKIT credentials")
	}

	// Keypairs belong to the identity, not the project or region; only list
	// them in the first project's first region scan.
	if len(opts.ProjectIDs) > 0 && opts.ProjectID != opts.ProjectIDs[0] {
		return []resource.Resource{}, nil
	}
	if len(opts.Regions) > 0 && opts.Region != opts.Regions[0] {
		return []resource.Resource{}, nil
	}

	client, err := iaasv2.NewAPIClient(stackitConfigOpts(opts)...)
	if err != nil {
		return nil, fmt.Errorf("build iaas client: %w", err)
	}

	logger.Trace("listing keypairs")
	resp, err := client.DefaultAPI.ListKeyPairs(ctx).Execute()
	if err != nil {
		return nil, fmt.Errorf("list keypairs: %w", err)
	}

	out := make([]resource.Resource, 0, len(resp.GetItems()))
	for _, k := range resp.GetItems() {
		name, ok := k.GetNameOk()
		if !ok || name == nil || *name == "" {
			continue
		}
		fingerprint, _ := k.GetFingerprintOk()
		createdAt, _ := k.GetCreatedAtOk()

		labels := map[string]string{}
		if raw, ok := k.GetLabelsOk(); ok {
			for key, val := range raw {
				if vs, ok := val.(string); ok {
					labels[key] = vs
				}
			}
		}

		out = append(out, &ComputeKeypair{
			BaseResource: &BaseResource{
				OrganizationID: opts.OrganizationID,
				ProjectID:      opts.ProjectID,
				Region:         opts.Region,
			},
			api:         client.DefaultAPI,
			Name:        *name,
			Fingerprint: stringDeref(fingerprint),
			CreatedAt:   createdAt,
			Labels:      labels,
		})
	}

	logger.WithField("count", len(out)).Trace("list complete")
	return out, nil
}
