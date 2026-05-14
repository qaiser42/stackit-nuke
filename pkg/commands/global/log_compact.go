package global

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

// compactHook strips libnuke's per-resource property dump from log entries.
// libnuke emits one INFO line per resource per state with every property as
// prop:<Field>=<value>, which gets unreadable fast. We keep only the fields
// that identify the resource (type, state, name, short ID) and drop the rest.
type compactHook struct{}

func (compactHook) Levels() []logrus.Level { return logrus.AllLevels }

func (compactHook) Fire(e *logrus.Entry) error {
	if len(e.Data) == 0 {
		return nil
	}
	out := logrus.Fields{}
	for k, v := range e.Data {
		switch k {
		case "_handler", "component", "owner", "state_code":
			continue
		case "type", "state", "name":
			out[k] = v
			continue
		}
		if !strings.HasPrefix(k, "prop:") {
			out[k] = v
			continue
		}
		switch strings.TrimPrefix(k, "prop:") {
		case "ID":
			out["id"] = shortID(fmt.Sprint(v))
		case "Name":
			if s := fmt.Sprint(v); s != "" {
				out["name"] = s
			}
		}
	}
	e.Data = out
	return nil
}

func shortID(s string) string {
	if i := strings.Index(s, "-"); i > 0 {
		return s[:i]
	}
	return s
}
