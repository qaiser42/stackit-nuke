package stackit

import (
	stackitconfig "github.com/stackitcloud/stackit-sdk-go/core/config"
)

// SDKConfigForRegion returns a copy of the base SDK config with the supplied
// region applied. Service-specific client constructors accept a
// *stackitconfig.Configuration directly (see services/* in stackit-sdk-go).
func (c *Credentials) SDKConfigForRegion(region string) *stackitconfig.Configuration {
	if c == nil || c.SDKConfig == nil {
		return nil
	}
	cfg := *c.SDKConfig
	if region != "" {
		cfg.Region = region
	}
	return &cfg
}
