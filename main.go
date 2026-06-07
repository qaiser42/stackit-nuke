package main

import (
	"context"
	"os"
	"path"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"

	"github.com/qaiser42/stackit-nuke/pkg/common"

	_ "github.com/qaiser42/stackit-nuke/pkg/commands/list"
	_ "github.com/qaiser42/stackit-nuke/pkg/commands/run"

	_ "github.com/qaiser42/stackit-nuke/resources"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(*logrus.Entry); ok {
				os.Exit(1)
			}
			panic(r)
		}
	}()

	cmd := &cli.Command{
		Name:     path.Base(os.Args[0]),
		Usage:    "remove all resources from a STACKIT project",
		Version:  common.AppVersion.Summary,
		Authors:  []any{"stackit-nuke contributors"},
		Commands: common.GetCommands(),
		CommandNotFound: func(_ context.Context, _ *cli.Command, command string) {
			logrus.Fatalf("Command %s not found.", command)
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		logrus.Fatal(err)
	}
}
