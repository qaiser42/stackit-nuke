package stackit

import (
	"fmt"
	"time"

	libnuke "github.com/ekristen/libnuke/pkg/nuke"
	"github.com/ekristen/libnuke/pkg/utils"
)

// Prompt asks the operator to type the project ID before any destructive
// action runs. With --no-prompt/--force the prompt is replaced by a sleep
// matching the value of --prompt-delay.
type Prompt struct {
	Parameters *libnuke.Parameters
	ProjectIDs []string
}

func (p *Prompt) Prompt() error {
	forceSleep := time.Duration(p.Parameters.ForceSleep) * time.Second
	target := p.ProjectIDs[0]

	fmt.Printf("Do you really want to nuke STACKIT project(s) %v?\n", p.ProjectIDs)
	if p.Parameters.Force {
		fmt.Printf("Waiting %v before continuing.\n", forceSleep)
		time.Sleep(forceSleep)
		return nil
	}
	fmt.Printf("Do you want to continue? Enter project ID %q to continue.\n", target)
	return utils.Prompt(target)
}
