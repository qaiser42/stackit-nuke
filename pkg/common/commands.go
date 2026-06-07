package common

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
)

var commands []*cli.Command

func RegisterCommand(command *cli.Command) {
	logrus.Debugln("Registering", command.Name, "command...")
	commands = append(commands, command)
}

func GetCommands() []*cli.Command {
	return commands
}
