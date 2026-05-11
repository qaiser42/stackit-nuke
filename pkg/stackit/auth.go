package stackit

import (
	"errors"
	"fmt"
	"os"

	stackitconfig "github.com/stackitcloud/stackit-sdk-go/core/config"
)

// Credentials holds the loaded service-account material plus a STACKIT
// SDK Configuration ready to be used with any service client.
type Credentials struct {
	ServiceAccountKeyPath string
	PrivateKeyPath        string
	Token                 string

	// SDKConfig is the base STACKIT SDK configuration. Per-service clients
	// receive a copy with the appropriate region overridden.
	SDKConfig *stackitconfig.Configuration
}

// LoadCredentials resolves auth material in the following order:
//
//  1. explicit service-account-key-path
//  2. STACKIT_SERVICE_ACCOUNT_KEY_PATH env var
//  3. STACKIT_SERVICE_ACCOUNT_TOKEN env var
//
// At least one must produce a valid configuration.
func LoadCredentials(saKeyPath, privateKeyPath, token string) (*Credentials, error) {
	if saKeyPath == "" {
		saKeyPath = os.Getenv("STACKIT_SERVICE_ACCOUNT_KEY_PATH")
	}
	if privateKeyPath == "" {
		privateKeyPath = os.Getenv("STACKIT_PRIVATE_KEY_PATH")
	}
	if token == "" {
		token = os.Getenv("STACKIT_SERVICE_ACCOUNT_TOKEN")
	}

	opts := []stackitconfig.ConfigurationOption{}
	switch {
	case saKeyPath != "":
		if _, err := os.Stat(saKeyPath); err != nil {
			return nil, fmt.Errorf("service-account key path %q: %w", saKeyPath, err)
		}
		opts = append(opts, stackitconfig.WithServiceAccountKeyPath(saKeyPath))
		if privateKeyPath != "" {
			opts = append(opts, stackitconfig.WithPrivateKeyPath(privateKeyPath))
		}
	case token != "":
		opts = append(opts, stackitconfig.WithToken(token))
	default:
		return nil, errors.New("no STACKIT credentials: set --auth-file or STACKIT_SERVICE_ACCOUNT_KEY_PATH or STACKIT_SERVICE_ACCOUNT_TOKEN")
	}

	cfg := &stackitconfig.Configuration{}
	for _, opt := range opts {
		if err := opt(cfg); err != nil {
			return nil, fmt.Errorf("apply STACKIT SDK config option: %w", err)
		}
	}

	return &Credentials{
		ServiceAccountKeyPath: saKeyPath,
		PrivateKeyPath:        privateKeyPath,
		Token:                 token,
		SDKConfig:             cfg,
	}, nil
}
