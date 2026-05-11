package config

import (
	"path/filepath"
	"testing"

	libconfig "github.com/ekristen/libnuke/pkg/config"
)

// Every shipped example under ../../examples must load cleanly. This guards
// against config-schema drift breaking docs.
func TestExamplesLoad(t *testing.T) {
	matches, err := filepath.Glob("../../examples/*.yaml")
	if err != nil {
		t.Fatalf("glob: %v", err)
	}
	if len(matches) == 0 {
		t.Fatal("no example yaml files found")
	}
	for _, path := range matches {
		t.Run(filepath.Base(path), func(t *testing.T) {
			if _, err := New(libconfig.Options{Path: path}); err != nil {
				t.Errorf("load %s: %v", path, err)
			}
		})
	}
}
