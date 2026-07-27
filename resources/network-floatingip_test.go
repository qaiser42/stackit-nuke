package resources

import (
	"context"
	"strings"
	"testing"

	"github.com/ekristen/libnuke/pkg/registry"

	"github.com/qaiser42/stackit-nuke/pkg/stackit"
)

func TestFloatingIP_StringAndProperties(t *testing.T) {
	r := &FloatingIP{
		BaseResource:     &BaseResource{ProjectID: "proj-1", Region: "eu01"},
		ID:               "ip-abc",
		IP:               "193.148.160.1",
		NetworkInterface: "nic-1",
	}

	if got := r.String(); got != "193.148.160.1" {
		t.Errorf("String() = %q, want 193.148.160.1", got)
	}

	props := r.Properties()
	rendered := props.String()
	for _, want := range []string{"ip-abc", "193.148.160.1", "nic-1", "proj-1", "eu01"} {
		if !strings.Contains(rendered, want) {
			t.Errorf("Properties() = %q; missing %q", rendered, want)
		}
	}
}

func TestFloatingIP_RemoveWithoutAPIErrors(t *testing.T) {
	r := &FloatingIP{
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

func TestFloatingIPLister_RequiresCredentials(t *testing.T) {
	l := &FloatingIPLister{}
	_, err := l.List(context.Background(), &stackit.ListerOpts{ProjectID: "p", Region: "r"})
	if err == nil {
		t.Fatal("expected error when credentials missing")
	}
}

func TestFloatingIPRegistration(t *testing.T) {
	regs := registry.GetRegistration(FloatingIPResource)
	if regs == nil {
		t.Fatal("FloatingIP not registered")
	}
	if regs.Scope != stackit.ProjectScope {
		t.Errorf("scope = %q, want %q", regs.Scope, stackit.ProjectScope)
	}
	if regs.Lister == nil || regs.Resource == nil {
		t.Error("lister or resource nil in registration")
	}
}
