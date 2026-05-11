package resources

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/ekristen/libnuke/pkg/registry"

	"github.com/qaiser42/stackit-nuke/pkg/stackit"
)

func TestComputeServer_StringAndProperties(t *testing.T) {
	ts := time.Date(2026, 5, 7, 12, 0, 0, 0, time.UTC)
	r := &ComputeServer{
		BaseResource:     &BaseResource{ProjectID: "proj-1", Region: "eu01"},
		ID:               "srv-abc",
		Name:             "web-1",
		AvailabilityZone: "eu01-1",
		CreatedAt:        &ts,
		Labels:           map[string]string{"env": "test"},
	}

	if got := r.String(); got != "web-1" {
		t.Errorf("String() = %q, want web-1", got)
	}

	props := r.Properties()
	rendered := props.String()
	for _, want := range []string{"web-1", "srv-abc", "proj-1", "eu01", "eu01-1"} {
		if !strings.Contains(rendered, want) {
			t.Errorf("Properties() = %q; missing %q", rendered, want)
		}
	}
}

func TestComputeServer_RemoveWithoutAPIErrors(t *testing.T) {
	r := &ComputeServer{
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

func TestComputeServerLister_RequiresCredentials(t *testing.T) {
	l := &ComputeServerLister{}
	_, err := l.List(context.Background(), &stackit.ListerOpts{ProjectID: "p", Region: "r"})
	if err == nil {
		t.Fatal("expected error when credentials missing")
	}
}

func TestComputeServerRegistration(t *testing.T) {
	regs := registry.GetRegistration(ComputeServerResource)
	if regs == nil {
		t.Fatal("ComputeServer not registered")
	}
	if regs.Scope != stackit.ProjectScope {
		t.Errorf("scope = %q, want %q", regs.Scope, stackit.ProjectScope)
	}
	if regs.Lister == nil || regs.Resource == nil {
		t.Error("lister or resource nil in registration")
	}
}
