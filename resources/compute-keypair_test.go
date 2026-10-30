package resources

import (
	"context"
	"strings"
	"testing"

	"github.com/ekristen/libnuke/pkg/registry"

	"github.com/qaiser42/stackit-nuke/pkg/stackit"
)

func TestComputeKeypair_StringAndProperties(t *testing.T) {
	r := &ComputeKeypair{
		BaseResource: &BaseResource{ProjectID: "proj-1", Region: "eu01"},
		Name:         "deploy-key",
		Fingerprint:  "aa:bb:cc:dd:ee:ff:00:11:22:33:44:55:66:77:88:99",
		Labels:       map[string]string{"env": "test"},
	}

	if got := r.String(); got != "deploy-key" {
		t.Errorf("String() = %q, want deploy-key", got)
	}

	props := r.Properties()
	rendered := props.String()
	for _, want := range []string{"deploy-key", "aa:bb:cc:dd:ee:ff:00:11:22:33:44:55:66:77:88:99", "proj-1", "eu01"} {
		if !strings.Contains(rendered, want) {
			t.Errorf("Properties() = %q; missing %q", rendered, want)
		}
	}
}

func TestComputeKeypair_RemoveWithoutAPIErrors(t *testing.T) {
	r := &ComputeKeypair{
		BaseResource: &BaseResource{ProjectID: "p", Region: "r"},
		Name:         "k",
	}
	err := r.Remove(context.Background())
	if err == nil {
		t.Fatal("expected error when api is unset")
	}
	if !strings.Contains(err.Error(), "api client not set") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestComputeKeypairLister_RequiresCredentials(t *testing.T) {
	l := &ComputeKeypairLister{}
	_, err := l.List(context.Background(), &stackit.ListerOpts{ProjectID: "p", Region: "r"})
	if err == nil {
		t.Fatal("expected error when credentials missing")
	}
}

func TestComputeKeypairLister_SkipsNonPrimaryScans(t *testing.T) {
	l := &ComputeKeypairLister{}
	for _, opts := range []*stackit.ListerOpts{
		{ProjectID: "p1", ProjectIDs: []string{"p1", "p2"}, Region: "eu02", Regions: []string{"eu01", "eu02"}},
		{ProjectID: "p2", ProjectIDs: []string{"p1", "p2"}, Region: "eu01", Regions: []string{"eu01", "eu02"}},
	} {
		opts.Credentials = &stackit.Credentials{}
		got, err := l.List(context.Background(), opts)
		if err != nil {
			t.Fatalf("project %s region %s: unexpected error: %v", opts.ProjectID, opts.Region, err)
		}
		if len(got) != 0 {
			t.Errorf("project %s region %s: expected no keypairs, got %d", opts.ProjectID, opts.Region, len(got))
		}
	}
}

func TestComputeKeypairRegistration(t *testing.T) {
	regs := registry.GetRegistration(ComputeKeypairResource)
	if regs == nil {
		t.Fatal("ComputeKeypair not registered")
	}
	if regs.Scope != stackit.ProjectScope {
		t.Errorf("scope = %q, want %q", regs.Scope, stackit.ProjectScope)
	}
	if regs.Lister == nil || regs.Resource == nil {
		t.Error("lister or resource nil in registration")
	}
}
