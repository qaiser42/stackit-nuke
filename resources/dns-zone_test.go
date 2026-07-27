package resources

import (
	"context"
	"strings"
	"testing"

	"github.com/ekristen/libnuke/pkg/registry"

	"github.com/qaiser42/stackit-nuke/pkg/stackit"
)

func TestDNSZone_StringAndProperties(t *testing.T) {
	r := &DNSZone{
		BaseResource: &BaseResource{ProjectID: "proj-1", Region: "eu01"},
		ID:           "zone-abc",
		Name:         "my-zone",
		DNSName:      "example.runs.onstackit.cloud",
		State:        "CREATE_SUCCEEDED",
	}

	if got := r.String(); got != "example.runs.onstackit.cloud" {
		t.Errorf("String() = %q, want example.runs.onstackit.cloud", got)
	}

	props := r.Properties()
	rendered := props.String()
	for _, want := range []string{"zone-abc", "my-zone", "example.runs.onstackit.cloud", "proj-1", "eu01", "CREATE_SUCCEEDED"} {
		if !strings.Contains(rendered, want) {
			t.Errorf("Properties() = %q; missing %q", rendered, want)
		}
	}
}

func TestDNSZone_RemoveWithoutAPIErrors(t *testing.T) {
	r := &DNSZone{
		BaseResource: &BaseResource{ProjectID: "p", Region: "r"},
		ID:           "id",
	}
	err := r.Remove(context.Background())
	if err == nil {
		t.Fatal("expected error when api is unset")
	}
	if !strings.Contains(err.Error(), "api client not set") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestDNSZoneLister_RequiresCredentials(t *testing.T) {
	l := &DNSZoneLister{}
	_, err := l.List(context.Background(), &stackit.ListerOpts{ProjectID: "p", Region: "r"})
	if err == nil {
		t.Fatal("expected error when credentials missing")
	}
}

func TestDNSZoneLister_SkipsNonPrimaryRegions(t *testing.T) {
	l := &DNSZoneLister{}
	got, err := l.List(context.Background(), &stackit.ListerOpts{
		ProjectID:   "p",
		Region:      "eu02",
		Regions:     []string{"eu01", "eu02"},
		Credentials: &stackit.Credentials{},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected no zones for non-primary region, got %d", len(got))
	}
}

func TestDNSZoneRegistration(t *testing.T) {
	regs := registry.GetRegistration(DNSZoneResource)
	if regs == nil {
		t.Fatal("DNSZone not registered")
	}
	if regs.Scope != stackit.ProjectScope {
		t.Errorf("scope = %q, want %q", regs.Scope, stackit.ProjectScope)
	}
	if regs.Lister == nil || regs.Resource == nil {
		t.Error("lister or resource nil in registration")
	}
}
