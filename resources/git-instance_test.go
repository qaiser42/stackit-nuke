package resources

import (
	"context"
	"strings"
	"testing"

	"github.com/ekristen/libnuke/pkg/registry"

	"github.com/qaiser42/stackit-nuke/pkg/stackit"
)

func TestGitInstance_StringAndProperties(t *testing.T) {
	r := &GitInstance{
		BaseResource: &BaseResource{ProjectID: "proj-1", Region: "eu01"},
		ID:           "git-abc",
		Name:         "my-git",
		URL:          "https://my-git.git.onstackit.cloud",
		Flavor:       "git-100",
		Version:      "v1.5.2",
		State:        "Ready",
	}

	if got := r.String(); got != "my-git" {
		t.Errorf("String() = %q, want my-git", got)
	}

	props := r.Properties()
	rendered := props.String()
	for _, want := range []string{"git-abc", "my-git", "https://my-git.git.onstackit.cloud", "git-100", "v1.5.2", "proj-1", "eu01", "Ready"} {
		if !strings.Contains(rendered, want) {
			t.Errorf("Properties() = %q; missing %q", rendered, want)
		}
	}
}

func TestGitInstance_RemoveWithoutAPIErrors(t *testing.T) {
	r := &GitInstance{
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

func TestGitInstanceLister_RequiresCredentials(t *testing.T) {
	l := &GitInstanceLister{}
	_, err := l.List(context.Background(), &stackit.ListerOpts{ProjectID: "p", Region: "r"})
	if err == nil {
		t.Fatal("expected error when credentials missing")
	}
}

func TestGitInstanceLister_SkipsNonPrimaryRegions(t *testing.T) {
	l := &GitInstanceLister{}
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
		t.Errorf("expected no instances for non-primary region, got %d", len(got))
	}
}

func TestGitInstanceRegistration(t *testing.T) {
	regs := registry.GetRegistration(GitInstanceResource)
	if regs == nil {
		t.Fatal("GitInstance not registered")
	}
	if regs.Scope != stackit.ProjectScope {
		t.Errorf("scope = %q, want %q", regs.Scope, stackit.ProjectScope)
	}
	if regs.Lister == nil || regs.Resource == nil {
		t.Error("lister or resource nil in registration")
	}
}
