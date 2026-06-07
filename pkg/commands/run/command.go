package run

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"

	libconfig "github.com/ekristen/libnuke/pkg/config"
	"github.com/ekristen/libnuke/pkg/filter"
	libnuke "github.com/ekristen/libnuke/pkg/nuke"
	"github.com/ekristen/libnuke/pkg/registry"
	libscanner "github.com/ekristen/libnuke/pkg/scanner"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/qaiser42/stackit-nuke/pkg/commands/global"
	"github.com/qaiser42/stackit-nuke/pkg/common"
	"github.com/qaiser42/stackit-nuke/pkg/config"
	"github.com/qaiser42/stackit-nuke/pkg/stackit"
)

func execute(ctx context.Context, c *cli.Command) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	logger := logrus.StandardLogger()

	parsedConfig, err := config.New(libconfig.Options{
		Path:         c.String("config"),
		Deprecations: registry.GetDeprecatedResourceTypeMapping(),
	})
	if err != nil {
		return err
	}

	saKey := firstNonEmpty(c.String("auth-file"), parsedConfig.Auth.ServiceAccountKeyPath)
	pkPath := firstNonEmpty(c.String("private-key-file"), parsedConfig.Auth.PrivateKeyPath)
	token := firstNonEmpty(c.String("token"), parsedConfig.Auth.Token)

	creds, err := stackit.LoadCredentials(saKey, pkPath, token)
	if err != nil {
		return err
	}

	// Refuse to run without an explicit project allow-list. config.New already
	// enforces this — the nil check is belt-and-suspenders.
	if len(parsedConfig.ProjectIDs) == 0 {
		return fmt.Errorf("project-ids must be set in config")
	}
	if len(c.StringSlice("project-id")) > 0 {
		// Caller may narrow further but never widen.
		for _, p := range c.StringSlice("project-id") {
			if !slices.Contains(parsedConfig.ProjectIDs, p) {
				return fmt.Errorf("project %q not in config project-ids allow-list", p)
			}
		}
		parsedConfig.ProjectIDs = c.StringSlice("project-id")
	}

	params := &libnuke.Parameters{
		Force:      c.Bool("no-prompt"),
		ForceSleep: c.Int("prompt-delay"),
		Quiet:      c.Bool("quiet"),
		NoDryRun:   c.Bool("no-dry-run"),
		Includes:   c.StringSlice("include"),
		Excludes:   c.StringSlice("exclude"),
	}
	if slices.Contains(c.StringSlice("feature-flag"), "wait-on-dependencies") {
		params.WaitOnDependencies = true
	}
	if slices.Contains(c.StringSlice("feature-flag"), "filter-groups") {
		params.UseFilterGroups = true
	}

	filters, err := parsedConfig.Filters(parsedConfig.ProjectIDs[0])
	if err != nil {
		return err
	}
	if len(filters[filter.Global]) == 0 {
		filters[filter.Global] = []filter.Filter{}
	}
	if !slices.Contains(parsedConfig.Regions, "all") {
		filters[filter.Global] = append(filters[filter.Global], filter.Filter{
			Property: "Region",
			Type:     filter.NotIn,
			Values:   parsedConfig.Regions,
		})
	}

	n := libnuke.New(params, filters, parsedConfig.Settings)
	n.SetRunSleep(5 * time.Second)
	n.SetLogger(logger.WithField("component", "nuke"))
	n.RegisterVersion(fmt.Sprintf("> %s", common.AppVersion.String()))

	// Dry-run is read-only — no confirmation needed. Only require the typed
	// project-id prompt when the operator has opted into real deletion.
	if params.NoDryRun {
		prompt := &stackit.Prompt{Parameters: params, ProjectIDs: parsedConfig.ProjectIDs}
		n.RegisterPrompt(prompt.Prompt)
	}

	projectResourceTypes := types.ResolveResourceTypes(
		registry.GetNamesForScope(stackit.ProjectScope),
		[]types.Collection{
			n.Parameters.Includes,
			parsedConfig.ResourceTypes.GetIncludes(),
		},
		[]types.Collection{
			n.Parameters.Excludes,
			parsedConfig.ResourceTypes.Excludes,
		},
		nil,
		nil,
	)

	for _, projectID := range parsedConfig.ProjectIDs {
		for _, region := range parsedConfig.Regions {
			scope := stackit.ProjectScope
			scannerName := fmt.Sprintf("project/%s/region/%s", projectID, region)
			sc, err := libscanner.New(&libscanner.Config{
				Owner:         scannerName,
				ResourceTypes: projectResourceTypes,
				Opts: &stackit.ListerOpts{
					OrganizationID: parsedConfig.OrganizationID,
					ProjectID:      projectID,
					Region:         region,
					Regions:        parsedConfig.Regions,
					Credentials:    creds,
				},
			})
			if err != nil {
				return err
			}
			if err := n.RegisterScanner(scope, sc); err != nil {
				return err
			}
		}
	}

	hook := &nukedHook{}
	logrus.AddHook(hook)

	logger.Debug("running ...")
	runErr := n.Run(ctx)
	if params.NoDryRun {
		printSummary(logger, hook.entries)
	}
	return runErr
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

func init() {
	flags := []cli.Flag{
		&cli.StringFlag{Name: "config", Usage: "path to config file", Value: "config.yaml"},
		&cli.StringSliceFlag{Name: "include", Usage: "only include this specific resource"},
		&cli.StringSliceFlag{Name: "exclude", Usage: "exclude this specific resource (overrides everything)"},
		&cli.BoolFlag{Name: "quiet", Aliases: []string{"q"}, Usage: "hide filtered messages"},
		&cli.BoolFlag{Name: "no-dry-run", Usage: "actually delete after discovery (default: dry-run)"},
		&cli.BoolFlag{Name: "no-prompt", Aliases: []string{"force"}, Usage: "skip the typed-confirmation prompt"},
		&cli.IntFlag{Name: "prompt-delay", Aliases: []string{"force-sleep"}, Usage: "seconds to wait after prompt before running (min 3)", Value: 10},
		&cli.StringSliceFlag{Name: "feature-flag", Usage: "enable experimental behavior"},
		&cli.StringFlag{Name: "auth-file", Usage: "path to STACKIT service-account key JSON", Sources: cli.EnvVars("STACKIT_SERVICE_ACCOUNT_KEY_PATH")},
		&cli.StringFlag{Name: "private-key-file", Usage: "path to RSA private key (for service-account key auth)", Sources: cli.EnvVars("STACKIT_PRIVATE_KEY_PATH")},
		&cli.StringFlag{Name: "token", Usage: "STACKIT bearer token", Sources: cli.EnvVars("STACKIT_SERVICE_ACCOUNT_TOKEN")},
		&cli.StringSliceFlag{Name: "project-id", Usage: "narrow nuke to one or more project IDs from the config allow-list", Sources: cli.EnvVars("STACKIT_PROJECT_ID")},
	}

	cmd := &cli.Command{
		Name:    "run",
		Aliases: []string{"nuke"},
		Usage:   "run nuke against a STACKIT project to remove configured resources",
		Flags:   append(flags, global.Flags()...),
		Before:  global.Before,
		Action:  execute,
	}
	common.RegisterCommand(cmd)
}
