package list

import (
	"fmt"
	"sort"

	"github.com/urfave/cli/v2"

	"github.com/ekristen/libnuke/pkg/registry"

	"github.com/qaiser42/stackit-nuke/pkg/commands/global"
	"github.com/qaiser42/stackit-nuke/pkg/common"
)

func execute(_ *cli.Context) error {
	names := registry.GetNames()
	sort.Strings(names)
	for _, n := range names {
		fmt.Println(n)
	}
	return nil
}

func init() {
	cmd := &cli.Command{
		Name:   "resource-types",
		Usage:  "list all registered resource types",
		Flags:  global.Flags(),
		Before: global.Before,
		Action: execute,
	}
	common.RegisterCommand(cmd)
}
