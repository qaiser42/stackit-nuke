package stackit

import "github.com/ekristen/libnuke/pkg/registry"

const (
	OrganizationScope registry.Scope = "organization"
	ProjectScope      registry.Scope = "project"
)

// ListerOpts is passed to every resource Lister. It carries the current
// project + region context plus the credential payload needed to build
// per-service clients on demand.
type ListerOpts struct {
	OrganizationID string
	ProjectID      string
	Region         string
	Regions        []string
	Credentials    *Credentials
}
