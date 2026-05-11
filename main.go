package main

import (
	"os"
	"path"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

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

	app := cli.NewApp()
	app.Name = path.Base(os.Args[0])
	app.Usage = "remove all resources from a STACKIT project"
	app.Version = common.AppVersion.Summary
	app.Authors = []*cli.Author{
		{Name: "stackit-nuke contributors"},
	}

	app.Commands = common.GetCommands()
	app.CommandNotFound = func(_ *cli.Context, command string) {
		logrus.Fatalf("Command %s not found.", command)
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
