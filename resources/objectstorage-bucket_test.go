package resources

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/stackitcloud/stackit-sdk-go/core/oapierror"

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

func TestS3EndpointFromURL(t *testing.T) {
	for in, want := range map[string]string{
		"https://object.storage.eu01.onstackit.cloud/my-bucket": "object.storage.eu01.onstackit.cloud",
		"":  "",
		"%": "",
	} {
		if got := s3EndpointFromURL(in); got != want {
			t.Errorf("s3EndpointFromURL(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestIsBucketNotEmpty(t *testing.T) {
	notEmpty := &oapierror.GenericOpenAPIError{
		StatusCode: 422,
		Body:       []byte(`{"detail":[{"key":"bucket.not_empty","msg":"The bucket is not empty"}]}`),
	}
	if !isBucketNotEmpty(fmt.Errorf("delete bucket: %w", error(notEmpty))) {
		t.Error("expected wrapped 422 bucket.not_empty to match")
	}
	for _, err := range []error{
		nil,
		errors.New("boom"),
		&oapierror.GenericOpenAPIError{StatusCode: 422, Body: []byte(`{"detail":[{"key":"other"}]}`)},
		&oapierror.GenericOpenAPIError{StatusCode: 404, Body: []byte(`bucket.not_empty`)},
	} {
		if isBucketNotEmpty(err) {
			t.Errorf("isBucketNotEmpty(%v) = true, want false", err)
		}
	}
}

func TestObjectStorageBucket_EmptyWithoutEndpointErrors(t *testing.T) {
	r := &ObjectStorageBucket{BaseResource: &BaseResource{ProjectID: "p", Region: "eu01"}, Name: "b"}
	if err := r.empty(context.Background()); err == nil || !strings.Contains(err.Error(), "no S3 endpoint") {
		t.Errorf("empty() = %v, want missing endpoint error", err)
	}
}
