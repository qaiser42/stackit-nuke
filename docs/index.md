# stackit-nuke

`stackit-nuke` removes **all** resources from a [STACKIT](https://stackit.de) project.

It is the STACKIT equivalent of [`aws-nuke`](https://github.com/ekristen/aws-nuke), [`azure-nuke`](https://github.com/ekristen/azure-nuke) and [`gcp-nuke`](https://github.com/ekristen/gcp-nuke), built on the same engine: [`libnuke`](https://github.com/ekristen/libnuke).

## Features

- Project-scoped destruction with explicit allow-list
- Dry-run by default; real deletion requires `--no-dry-run`
- YAML config with includes/excludes, filters, presets, blocklists (libnuke schema)
- Service-account key authentication (STACKIT standard)
- Multi-region in a single run
- Dependency-aware deletion order
- Distroless multi-arch container image
- Signed binaries + Cosign-signed images

## Status

Early — resource Listers/Removers are placeholders. The CLI, config, and engine are wired up. Resource implementations land incrementally; see [Contributing](contributing.md).

## Read next

- [Warning](warning.md) — read this first
- [Install](install.md)
- [Authentication](auth.md)
- [Quick Start](quick-start.md)
