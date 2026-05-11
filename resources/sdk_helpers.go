package resources

import (
	stackitconfig "github.com/stackitcloud/stackit-sdk-go/core/config"

	"github.com/qaiser42/stackit-nuke/pkg/stackit"
)

// stackitConfigOpts is the bridge from a Lister's incoming ListerOpts to the
// option-functions every STACKIT service-client constructor accepts. Each
// service does its own NewAPIClient(opts...); this keeps the wiring uniform.
func stackitConfigOpts(opts *stackit.ListerOpts) []stackitconfig.ConfigurationOption {
	// Note: region is NOT applied here. The IaaS v2 client (and most newer
	// STACKIT services) take region as a per-call parameter rather than a
	// client-wide setting; passing WithRegion makes NewAPIClient reject the
	// config. Region comes through ListerOpts.Region into each SDK call.
	cfg := opts.Credentials.SDKConfig
	var out []stackitconfig.ConfigurationOption
	if cfg.Token != "" {
		out = append(out, stackitconfig.WithToken(cfg.Token))
	}
	if cfg.ServiceAccountKeyPath != "" {
		out = append(out, stackitconfig.WithServiceAccountKeyPath(cfg.ServiceAccountKeyPath))
	}
	if cfg.PrivateKeyPath != "" {
		out = append(out, stackitconfig.WithPrivateKeyPath(cfg.PrivateKeyPath))
	}
	return out
}

func stringDeref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
