package config

import (
	"os"
	"path/filepath"
	"testing"

	libconfig "github.com/ekristen/libnuke/pkg/config"
)

func TestNew_SeedsAccountsFromProjectIDs(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	body := `
regions: [eu01]
project-ids:
  - 11111111-1111-1111-1111-111111111111
auth:
  service-account-key-path: /tmp/sa.json
`
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}

	cfg, err := New(libconfig.Options{Path: path})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	acc, ok := cfg.Accounts["11111111-1111-1111-1111-111111111111"]
	if !ok {
		t.Fatal("project-id not seeded into Accounts")
	}
	if acc == nil {
		t.Fatal("seeded Account is nil")
	}

	// Filters() must now succeed for this project id without an explicit
	// accounts: block in the YAML.
	if _, err := cfg.Filters("11111111-1111-1111-1111-111111111111"); err != nil {
		t.Errorf("Filters returned error after seed: %v", err)
	}
}

func TestNew_DoesNotOverwriteExistingAccount(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	body := `
regions: [eu01]
project-ids:
  - 11111111-1111-1111-1111-111111111111
accounts:
  "11111111-1111-1111-1111-111111111111":
    presets: [keep-stuff]
presets:
  keep-stuff:
    filters:
      ComputeServer:
        - property: Name
          value: "keep-*"
          type: glob
auth:
  service-account-key-path: /tmp/sa.json
`
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}

	cfg, err := New(libconfig.Options{Path: path})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	acc := cfg.Accounts["11111111-1111-1111-1111-111111111111"]
	if len(acc.Presets) != 1 || acc.Presets[0] != "keep-stuff" {
		t.Errorf("existing account overwritten: presets = %v", acc.Presets)
	}
}
