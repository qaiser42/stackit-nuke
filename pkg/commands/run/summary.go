package run

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

// nukedHook captures every "removed" log entry libnuke emits when a resource
// finishes deletion. Collected entries feed printSummary at end of run.
type nukedHook struct {
	mu      sync.Mutex
	entries []nukedEntry
}

type nukedEntry struct {
	Type string
	Name string
	ID   string
}

func (*nukedHook) Levels() []logrus.Level { return []logrus.Level{logrus.InfoLevel} }

func (h *nukedHook) Fire(e *logrus.Entry) error {
	if e.Message != "removed" {
		return nil
	}
	entry := nukedEntry{}
	if v, ok := e.Data["type"]; ok {
		entry.Type = fmt.Sprint(v)
	}
	if v, ok := e.Data["name"]; ok {
		entry.Name = fmt.Sprint(v)
	}
	// id may come through as the short form (compact hook) or full prop:ID.
	if v, ok := e.Data["id"]; ok {
		entry.ID = fmt.Sprint(v)
	} else if v, ok := e.Data["prop:ID"]; ok {
		entry.ID = fmt.Sprint(v)
	}
	h.mu.Lock()
	h.entries = append(h.entries, entry)
	h.mu.Unlock()
	return nil
}

const nukeArt = `
                  _ _.-'` + "`" + `-._
              .-'` + "`" + `           '-.
            .'      _ . - = - . _   '.
           /    .-'             '-.   \
          /   .'                   '.  \
         |   /                       \  |
         |  |          BOOM           | |
         |  |                         | |
          \  '.                     .'  /
           \   '-.                .-'  /
            '.    '-._        _.-'   .'
              '-.     ` + "`" + `'----'` + "`" + `    .-'
                 ` + "`" + `'-..__________..-'` + "`" + `
`

func printSummary(logger *logrus.Logger, entries []nukedEntry) {
	if len(entries) == 0 {
		return
	}

	sort.SliceStable(entries, func(i, j int) bool {
		if entries[i].Type != entries[j].Type {
			return entries[i].Type < entries[j].Type
		}
		return label(entries[i]) < label(entries[j])
	})

	var b strings.Builder
	b.WriteString(nukeArt)
	b.WriteString(fmt.Sprintf("\n%d resource(s) nuked:\n\n", len(entries)))

	currentType := ""
	for _, e := range entries {
		if e.Type != currentType {
			fmt.Fprintf(&b, "  %s\n", e.Type)
			currentType = e.Type
		}
		fmt.Fprintf(&b, "    - %s\n", label(e))
	}

	// The summary is final human-facing report output (ASCII art + a
	// multi-line list), not a structured log record. Routing it through
	// logrus makes the TextFormatter quote the msg field and escape every
	// newline to a literal \n. Write straight to the logger's output so the
	// banner renders with real line breaks while still honouring SetOutput.
	_, _ = fmt.Fprint(logger.Out, b.String())
}

func label(e nukedEntry) string {
	switch {
	case e.Name != "" && e.ID != "":
		return fmt.Sprintf("%s (%s)", e.Name, e.ID)
	case e.Name != "":
		return e.Name
	case e.ID != "":
		return e.ID
	default:
		return "<unknown>"
	}
}
