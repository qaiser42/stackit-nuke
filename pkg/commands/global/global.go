package global

import (
	"context"
	"fmt"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
)

func Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "log-level",
			Usage:   "log level",
			Aliases: []string{"l"},
			Sources: cli.EnvVars("LOGLEVEL"),
			Value:   "info",
		},
		&cli.BoolFlag{Name: "log-caller", Usage: "log caller (file:line)"},
		&cli.BoolFlag{Name: "log-disable-color", Usage: "disable log coloring"},
		&cli.BoolFlag{Name: "log-full-timestamp", Usage: "always show full timestamp"},
		&cli.BoolFlag{Name: "log-verbose", Usage: "show every libnuke property in log lines (default: compact)"},
	}
}

func Before(ctx context.Context, c *cli.Command) (context.Context, error) {
	formatter := &logrus.TextFormatter{
		DisableColors: c.Bool("log-disable-color"),
		FullTimestamp: c.Bool("log-full-timestamp"),
	}
	if c.Bool("log-caller") {
		logrus.SetReportCaller(true)
		formatter.CallerPrettyfier = func(f *runtime.Frame) (string, string) {
			return "", fmt.Sprintf("%s:%d", path.Base(f.File), f.Line)
		}
	}
	logrus.SetFormatter(formatter)

	if !c.Bool("log-verbose") {
		logrus.AddHook(compactHook{})
	}

	switch c.String("log-level") {
	case "trace":
		logrus.SetLevel(logrus.TraceLevel)
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	}

	return ctx, nil
}
