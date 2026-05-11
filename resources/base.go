// Package resources contains one file per STACKIT resource type that
// stackit-nuke knows how to discover and delete. Each file registers a
// libnuke registration during init() and exposes a Lister that returns
// resources observed in the supplied project + region.
//
// The Lister bodies are placeholders — they return an empty slice so the
// framework wires up correctly end-to-end. Filling them in with real
// STACKIT SDK calls is the per-resource implementation work.
package resources

import (
	"github.com/ekristen/libnuke/pkg/types"
)

// BaseResource embeds the project + region context every STACKIT resource
// shares. Use struct embedding with `property:",inline"` so the fields
// participate in libnuke's filter/property system.
type BaseResource struct {
	OrganizationID string `property:"OrganizationID"`
	ProjectID      string `property:"ProjectID"`
	Region         string `property:"Region"`
}

// PropsFromStruct is a small helper so resource Properties() bodies stay
// one-liners.
func PropsFromStruct(v any) types.Properties {
	return types.NewPropertiesFromStruct(v)
}
