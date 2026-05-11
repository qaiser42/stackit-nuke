package stackit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadCredentials_NoSourcesErrors(t *testing.T) {
	t.Setenv("STACKIT_SERVICE_ACCOUNT_KEY_PATH", "")
	t.Setenv("STACKIT_PRIVATE_KEY_PATH", "")
	t.Setenv("STACKIT_SERVICE_ACCOUNT_TOKEN", "")

	_, err := LoadCredentials("", "", "")
	if err == nil {
		t.Fatal("expected error when no credentials provided")
	}
	if !strings.Contains(err.Error(), "no STACKIT credentials") {
		t.Errorf("error %q should explain missing credentials", err)
	}
}

func TestLoadCredentials_TokenOnly(t *testing.T) {
	t.Setenv("STACKIT_SERVICE_ACCOUNT_KEY_PATH", "")
	t.Setenv("STACKIT_PRIVATE_KEY_PATH", "")

	creds, err := LoadCredentials("", "", "test-token-xyz")
	if err != nil {
		t.Fatalf("LoadCredentials: %v", err)
	}
	if creds.Token != "test-token-xyz" {
		t.Errorf("token mismatch: %q", creds.Token)
	}
	if creds.SDKConfig == nil {
		t.Fatal("SDKConfig nil")
	}
	if creds.SDKConfig.Token != "test-token-xyz" {
		t.Errorf("SDKConfig.Token = %q", creds.SDKConfig.Token)
	}
}

func TestLoadCredentials_KeyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sa.json")
	if err := os.WriteFile(path, []byte(`{"id":"x"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("STACKIT_SERVICE_ACCOUNT_KEY_PATH", "")
	t.Setenv("STACKIT_SERVICE_ACCOUNT_TOKEN", "")

	creds, err := LoadCredentials(path, "", "")
	if err != nil {
		t.Fatalf("LoadCredentials: %v", err)
	}
	if creds.SDKConfig.ServiceAccountKeyPath != path {
		t.Errorf("ServiceAccountKeyPath = %q, want %q", creds.SDKConfig.ServiceAccountKeyPath, path)
	}
}

func TestLoadCredentials_KeyFileMissing(t *testing.T) {
	t.Setenv("STACKIT_SERVICE_ACCOUNT_KEY_PATH", "")
	_, err := LoadCredentials("/does/not/exist.json", "", "")
	if err == nil {
		t.Fatal("expected error for missing key file")
	}
}

func TestSDKConfigForRegion(t *testing.T) {
	creds, err := LoadCredentials("", "", "tok")
	if err != nil {
		t.Fatal(err)
	}
	cfg := creds.SDKConfigForRegion("eu02")
	if cfg.Region != "eu02" {
		t.Errorf("Region = %q, want eu02", cfg.Region)
	}
	// base config must remain unmodified — we returned a copy
	if creds.SDKConfig.Region == "eu02" {
		t.Error("SDKConfigForRegion mutated base config")
	}
}
