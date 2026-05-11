package common

import (
	"strings"
	"testing"
)

func TestAppVersionPopulated(t *testing.T) {
	if AppVersion.Name != NAME {
		t.Errorf("AppVersion.Name = %q, want %q", AppVersion.Name, NAME)
	}
	if AppVersion.Summary == "" {
		t.Error("AppVersion.Summary empty")
	}
	got := AppVersion.String()
	for _, want := range []string{NAME, AppVersion.Summary, AppVersion.Commit} {
		if !strings.Contains(got, want) {
			t.Errorf("AppVersion.String() = %q, want to contain %q", got, want)
		}
	}
}
