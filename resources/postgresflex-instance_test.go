package resources

import (
	"context"
	"strings"
	"testing"

	"github.com/ekristen/libnuke/pkg/registry"

	"github.com/qaiser42/stackit-nuke/pkg/stackit"
)

func TestPostgresFlexInstance_StringAndProperties(t *testing.T) {
	r := &PostgresFlexInstance{
		BaseResource: &BaseResource{ProjectID: "proj-1", Region: "eu01"},
		ID:           "pg-abc",
		Name:         "orders-db",
		Status:       "Ready",
	}

	if got := r.String(); got != "orders-db" {
		t.Errorf("String() = %q, want orders-db", got)
	}

	props := r.Properties()
	rendered := props.String()
	for _, want := range []string{"orders-db", "pg-abc", "proj-1", "eu01", "Ready"} {
		if !strings.Contains(rendered, want) {
			t.Errorf("Properties() = %q; missing %q", rendered, want)
		}
	}
}

func TestPostgresFlexInstance_RemoveWithoutAPIErrors(t *testing.T) {
	r := &PostgresFlexInstance{
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

func TestPostgresFlexInstanceLister_RequiresCredentials(t *testing.T) {
	l := &PostgresFlexInstanceLister{}
	_, err := l.List(context.Background(), &stackit.ListerOpts{ProjectID: "p", Region: "r"})
	if err == nil {
		t.Fatal("expected error when credentials missing")
	}
}

func TestPostgresFlexInstanceRegistration(t *testing.T) {
	regs := registry.GetRegistration(PostgresFlexInstanceResource)
	if regs == nil {
		t.Fatal("PostgresFlexInstance not registered")
	}
	if regs.Scope != stackit.ProjectScope {
		t.Errorf("scope = %q, want %q", regs.Scope, stackit.ProjectScope)
	}
	if regs.Lister == nil || regs.Resource == nil {
		t.Error("lister or resource nil in registration")
	}
}
