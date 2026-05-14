package resources

import (
	"sort"
	"testing"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"

	"github.com/qaiser42/stackit-nuke/pkg/stackit"
)

// expected enumerates every resource the v0.1 scaffold registers. Adding a
// new resource means adding it here too — this catches accidental drops.
var expected = []string{
	"ComputeServer",
	"ComputeVolume",
	"ComputeSnapshot",
	"ComputeKeypair",
	"Network",
	"Subnet",
	"Router",
	"SecurityGroup",
	"FloatingIP",
	"ObjectStorageBucket",
	"ObjectStorageObject",
	"SKECluster",
	"PostgresFlexInstance",
	"MongoDBFlexInstance",
	"RedisInstance",
	"OpenSearchInstance",
	"RabbitMQInstance",
	"LoadBalancer",
	"DNSZone",
	"NetworkInterface",
}

func TestAllResourcesRegistered(t *testing.T) {
	got := registry.GetNames()
	sort.Strings(got)
	want := append([]string(nil), expected...)
	sort.Strings(want)

	gotSet := map[string]struct{}{}
	for _, n := range got {
		gotSet[n] = struct{}{}
	}
	for _, n := range want {
		if _, ok := gotSet[n]; !ok {
			t.Errorf("resource %q not registered", n)
		}
	}
}

func TestAllResourcesUseProjectScope(t *testing.T) {
	for _, name := range expected {
		names := registry.GetNamesForScope(stackit.ProjectScope)
		found := false
		for _, n := range names {
			if n == name {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("resource %q not in ProjectScope", name)
		}
	}
}

// realImpls are resources whose Listers actually call the STACKIT API.
// They are excluded from the stub-emptiness check below; they get their own
// per-resource tests instead.
var realImpls = map[string]bool{
	"ComputeServer":    true,
	"ComputeVolume":    true,
	"NetworkInterface": true,
	"Network":          true,
	"SecurityGroup":    true,
}

func TestListersReturnEmpty(t *testing.T) {
	// Stub Listers must not panic and must return empty slices. As real
	// implementations land they move into realImpls and are tested per-file.
	listers := registry.GetListers()
	for _, name := range expected {
		if realImpls[name] {
			continue
		}
		l, ok := listers[name]
		if !ok {
			t.Errorf("no lister for %q", name)
			continue
		}
		opts := &stackit.ListerOpts{ProjectID: "test", Region: "eu01"}
		got, err := l.List(t.Context(), opts)
		if err != nil {
			t.Errorf("%s.List error: %v", name, err)
		}
		if len(got) != 0 {
			t.Errorf("%s.List returned %d items, expected stub-empty", name, len(got))
		}
		_ = resource.Resource(nil) // silence unused import
	}
}
