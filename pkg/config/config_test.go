package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	libconfig "github.com/ekristen/libnuke/pkg/config"
)

func writeConfig(t *testing.T, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("write tmp config: %v", err)
	}
	return path
}

func TestNew_Valid(t *testing.T) {
	path := writeConfig(t, `
regions: [eu01]
project-ids:
  - 11111111-1111-1111-1111-111111111111
auth:
  service-account-key-path: /tmp/sa-key.json
`)
	cfg, err := New(libconfig.Options{Path: path})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if got := cfg.ProjectIDs[0]; got != "11111111-1111-1111-1111-111111111111" {
		t.Errorf("project id mismatch: %q", got)
	}
	if got := cfg.Regions[0]; got != "eu01" {
		t.Errorf("region mismatch: %q", got)
	}
	if got := cfg.Auth.ServiceAccountKeyPath; got != "/tmp/sa-key.json" {
		t.Errorf("auth path mismatch: %q", got)
	}
}

func TestNew_RejectsMissingProjectIDs(t *testing.T) {
	path := writeConfig(t, `regions: [eu01]`)
	_, err := New(libconfig.Options{Path: path})
	if err == nil {
		t.Fatal("expected error when project-ids missing")
	}
	if !strings.Contains(err.Error(), "project-ids") {
		t.Errorf("error %q should mention project-ids", err)
	}
}

func TestNew_RejectsMissingRegions(t *testing.T) {
	path := writeConfig(t, `
project-ids: [11111111-1111-1111-1111-111111111111]
`)
	_, err := New(libconfig.Options{Path: path})
	if err == nil {
		t.Fatal("expected error when regions missing")
	}
	if !strings.Contains(err.Error(), "regions") {
		t.Errorf("error %q should mention regions", err)
	}
}

func TestNew_FileNotFound(t *testing.T) {
	_, err := New(libconfig.Options{Path: "/does/not/exist.yaml"})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
