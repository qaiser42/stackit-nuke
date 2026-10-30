package resources

import (
	"context"
	"strings"
	"testing"

	"github.com/ekristen/libnuke/pkg/registry"

	"github.com/qaiser42/stackit-nuke/pkg/stackit"
)

func TestMariaDBInstance_StringAndProperties(t *testing.T) {
	r := &MariaDBInstance{
		BaseResource: &BaseResource{ProjectID: "proj-1", Region: "eu01"},
		ID:           "mdb-abc",
		Name:         "orders-mariadb",
		Plan:         "stackit-mariadb-1.4.10-single",
		Version:      "10.11",
		Status:       "active",
	}

	if got := r.String(); got != "orders-mariadb" {
		t.Errorf("String() = %q, want orders-mariadb", got)
	}

	props := r.Properties()
	rendered := props.String()
	for _, want := range []string{"mdb-abc", "orders-mariadb", "stackit-mariadb-1.4.10-single", "10.11", "proj-1", "eu01", "active"} {
		if !strings.Contains(rendered, want) {
			t.Errorf("Properties() = %q; missing %q", rendered, want)
		}
	}
}

func TestMariaDBInstance_RemoveWithoutAPIErrors(t *testing.T) {
	r := &MariaDBInstance{
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

func TestMariaDBInstanceLister_RequiresCredentials(t *testing.T) {
	l := &MariaDBInstanceLister{}
	_, err := l.List(context.Background(), &stackit.ListerOpts{ProjectID: "p", Region: "r"})
	if err == nil {
		t.Fatal("expected error when credentials missing")
	}
}

func TestMariaDBInstanceRegistration(t *testing.T) {
	regs := registry.GetRegistration(MariaDBInstanceResource)
	if regs == nil {
		t.Fatal("MariaDBInstance not registered")
	}
	if regs.Scope != stackit.ProjectScope {
		t.Errorf("scope = %q, want %q", regs.Scope, stackit.ProjectScope)
	}
	if regs.Lister == nil || regs.Resource == nil {
		t.Error("lister or resource nil in registration")
	}
}
