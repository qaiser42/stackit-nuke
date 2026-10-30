package resources

import (
	"context"
	"strings"
	"testing"

	"github.com/ekristen/libnuke/pkg/registry"

	"github.com/qaiser42/stackit-nuke/pkg/stackit"
)

func TestComputeSnapshot_StringAndProperties(t *testing.T) {
	r := &ComputeSnapshot{
		BaseResource: &BaseResource{ProjectID: "proj-1", Region: "eu01"},
		ID:           "snap-abc",
		Name:         "nightly",
		VolumeID:     "vol-xyz",
		SizeGB:       20,
		Status:       "AVAILABLE",
		Labels:       map[string]string{"env": "test"},
	}

	if got := r.String(); got != "nightly" {
		t.Errorf("String() = %q, want nightly", got)
	}

	props := r.Properties()
	rendered := props.String()
	for _, want := range []string{"snap-abc", "nightly", "vol-xyz", "proj-1", "eu01", "AVAILABLE"} {
		if !strings.Contains(rendered, want) {
			t.Errorf("Properties() = %q; missing %q", rendered, want)
		}
	}
}

func TestComputeSnapshot_RemoveWithoutAPIErrors(t *testing.T) {
	r := &ComputeSnapshot{
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

func TestComputeSnapshotLister_RequiresCredentials(t *testing.T) {
	l := &ComputeSnapshotLister{}
	_, err := l.List(context.Background(), &stackit.ListerOpts{ProjectID: "p", Region: "r"})
	if err == nil {
		t.Fatal("expected error when credentials missing")
	}
}

func TestComputeSnapshotRegistration(t *testing.T) {
	regs := registry.GetRegistration(ComputeSnapshotResource)
	if regs == nil {
		t.Fatal("ComputeSnapshot not registered")
	}
	if regs.Scope != stackit.ProjectScope {
		t.Errorf("scope = %q, want %q", regs.Scope, stackit.ProjectScope)
	}
	if regs.Lister == nil || regs.Resource == nil {
		t.Error("lister or resource nil in registration")
	}
}

func TestComputeVolumeDependsOnSnapshot(t *testing.T) {
	regs := registry.GetRegistration(ComputeVolumeResource)
	if regs == nil {
		t.Fatal("ComputeVolume not registered")
	}
	for _, d := range regs.DependsOn {
		if d == ComputeSnapshotResource {
			return
		}
	}
	t.Errorf("ComputeVolume.DependsOn = %v; missing %q", regs.DependsOn, ComputeSnapshotResource)
}
