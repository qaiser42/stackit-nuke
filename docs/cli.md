# CLI

```
stackit-nuke <command> [flags]
```

## Commands

| Command | Description |
|---------|-------------|
| `run` (alias `nuke`) | scan + delete resources for the configured projects |
| `resource-types` | list all registered resource type names |

## `run` flags

| Flag | Env | Default | Description |
|------|-----|---------|-------------|
| `--config` | | `config.yaml` | path to config file |
| `--no-dry-run` | | `false` | actually perform deletions |
| `--no-prompt` / `--force` | | `false` | skip typed-confirmation |
| `--prompt-delay` / `--force-sleep` | | `10` | seconds to wait after prompt |
| `--include` (repeatable) | | | only consider these resource types |
| `--exclude` (repeatable) | | | always exclude these resource types |
| `--quiet` / `-q` | | `false` | hide filtered messages |
| `--feature-flag` (repeatable) | | | enable experimental engine behaviors |
| `--auth-file` | `STACKIT_SERVICE_ACCOUNT_KEY_PATH` | | service-account key JSON path |
| `--private-key-file` | `STACKIT_PRIVATE_KEY_PATH` | | RSA private key path |
| `--token` | `STACKIT_SERVICE_ACCOUNT_TOKEN` | | bearer token |
| `--project-id` (repeatable) | `STACKIT_PROJECT_ID` | | narrow to subset of `project-ids` |
| `--log-level` / `-l` | `LOGLEVEL` | `info` | `trace`/`debug`/`info`/`warn`/`error` |

## Exit codes

| Code | Meaning |
|------|---------|
| `0` | clean — nothing failed |
| `1` | a resource failed to delete (libnuke surfaces details in the log) |
| any other | configuration / authentication error before scanning began |
