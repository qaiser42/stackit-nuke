package config

import (
	"fmt"
	"os"

	libconfig "github.com/ekristen/libnuke/pkg/config"
	"gopkg.in/yaml.v3"
)

// Auth carries STACKIT-specific credential locations from the config file.
// Any of these can also be supplied via CLI flag or environment variable;
// the CLI takes precedence over the config file.
type Auth struct {
	ServiceAccountKeyPath string `yaml:"service-account-key-path"`
	PrivateKeyPath        string `yaml:"private-key-path"`
	Token                 string `yaml:"token"`
}

// Config extends the libnuke configuration with STACKIT project + auth fields.
// Filtering, presets, blocklists, includes/excludes are all inherited from
// libnuke and behave identically to aws-nuke / azure-nuke.
type Config struct {
	libconfig.Config `yaml:",inline"`

	OrganizationID string   `yaml:"organization-id"`
	ProjectIDs     []string `yaml:"project-ids"`
	Auth           Auth     `yaml:"auth"`
}

// New loads the libnuke base config plus the STACKIT extension fields from
// the same file.
func New(opts libconfig.Options) (*Config, error) {
	base, err := libconfig.New(opts)
	if err != nil {
		return nil, err
	}

	c := &Config{}
	if err := loadYAML(opts.Path, c); err != nil {
		return nil, err
	}
	c.Config = *base

	if len(c.ProjectIDs) == 0 {
		return nil, fmt.Errorf("config %q: project-ids is required (allow-list of STACKIT projects to nuke)", opts.Path)
	}
	if len(c.Regions) == 0 {
		return nil, fmt.Errorf("config %q: regions is required", opts.Path)
	}

	// libnuke's filter/preset machinery looks up filters by account ID. Seed
	// an empty entry for every project-id that lacks one so callers don't
	// have to duplicate the allow-list under `accounts:`.
	if c.Accounts == nil {
		c.Accounts = map[string]*libconfig.Account{}
	}
	for _, pid := range c.ProjectIDs {
		if _, ok := c.Accounts[pid]; !ok {
			c.Accounts[pid] = &libconfig.Account{}
		}
	}

	return c, nil
}

func loadYAML(path string, out any) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read config %q: %w", path, err)
	}
	if err := yaml.Unmarshal(b, out); err != nil {
		return fmt.Errorf("parse config %q: %w", path, err)
	}
	return nil
}
