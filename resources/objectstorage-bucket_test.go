package resources

import (
	"context"
	"strings"
	"testing"

	"github.com/ekristen/libnuke/pkg/registry"

	"github.com/qaiser42/stackit-nuke/pkg/stackit"
)

func TestObjectStorageBucket_StringAndProperties(t *testing.T) {
	r := &ObjectStorageBucket{
		BaseResource:      &BaseResource{ProjectID: "proj-1", Region: "eu01"},
		Name:              "my-bucket",
		ObjectLockEnabled: true,
		URL:               "https://my-bucket.object.storage.eu01.onstackit.cloud",
	}

	if got := r.String(); got != "my-bucket" {
		t.Errorf("String() = %q, want my-bucket", got)
	}

	props := r.Properties()
	rendered := props.String()
	for _, want := range []string{"my-bucket", "https://my-bucket.object.storage.eu01.onstackit.cloud", "proj-1", "eu01", "true"} {
		if !strings.Contains(rendered, want) {
			t.Errorf("Properties() = %q; missing %q", rendered, want)
		}
	}
}

func TestObjectStorageBucket_RemoveWithoutAPIErrors(t *testing.T) {
	r := &ObjectStorageBucket{
		BaseResource: &BaseResource{ProjectID: "p", Region: "r"},
		Name:         "b",
	}
	err := r.Remove(context.Background())
	if err == nil {
		t.Fatal("expected error when api is unset")
	}
	if !strings.Contains(err.Error(), "api client not set") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestObjectStorageBucketLister_RequiresCredentials(t *testing.T) {
	l := &ObjectStorageBucketLister{}
	_, err := l.List(context.Background(), &stackit.ListerOpts{ProjectID: "p", Region: "r"})
	if err == nil {
		t.Fatal("expected error when credentials missing")
	}
}

func TestObjectStorageBucketRegistration(t *testing.T) {
	regs := registry.GetRegistration(ObjectStorageBucketResource)
	if regs == nil {
		t.Fatal("ObjectStorageBucket not registered")
	}
	if regs.Scope != stackit.ProjectScope {
		t.Errorf("scope = %q, want %q", regs.Scope, stackit.ProjectScope)
	}
	if regs.Lister == nil || regs.Resource == nil {
		t.Error("lister or resource nil in registration")
	}
}
