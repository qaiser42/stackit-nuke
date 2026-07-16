# stackit-nuke

<p align="center">
  <img src="docs/assets/nuke.jpg" alt="stackit-nuke" width="420">
</p>

> Remove **all** resources from a [STACKIT](https://stackit.de) project. The STACKIT counterpart of [`aws-nuke`](https://github.com/ekristen/aws-nuke), built on the same engine: [`libnuke`](https://github.com/ekristen/libnuke).

[![ci](https://github.com/qaiser42/stackit-nuke/actions/workflows/ci.yml/badge.svg)](https://github.com/qaiser42/stackit-nuke/actions/workflows/ci.yml)
[![release](https://github.com/qaiser42/stackit-nuke/actions/workflows/release.yml/badge.svg)](https://github.com/qaiser42/stackit-nuke/actions/workflows/release.yml)
[![docs](https://github.com/qaiser42/stackit-nuke/actions/workflows/docs.yml/badge.svg)](https://qaiser42.io/stackit-nuke)
[![security](https://github.com/qaiser42/stackit-nuke/actions/workflows/security.yml/badge.svg)](https://github.com/qaiser42/stackit-nuke/actions/workflows/security.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/qaiser42/stackit-nuke.svg)](https://pkg.go.dev/github.com/qaiser42/stackit-nuke)

## ⚠️ Read this first

`stackit-nuke` is **destructive by design**. It deletes everything in the targeted STACKIT project that matches the configured resource types. Deletion is permanent.

Read [docs/warning.md](docs/warning.md) before running.

## Resource coverage

Legend: ✅ list + delete via real STACKIT SDK · 🟡 registered, lister returns empty · ⬜ not yet registered

| Service | Resource | Status | SDK package |
|---|---|---|---|
| IaaS / compute | `ComputeServer` | ✅ | `stackit-sdk-go/services/iaas/v2api` |
| IaaS / compute | `ComputeVolume` | ✅ | `iaas/v2api` |
| IaaS / compute | `ComputeSnapshot` | 🟡 | `iaas/v2api` |
| IaaS / compute | `ComputeKeypair` | 🟡 | `iaas/v2api` |
| IaaS / network | `Network` | ✅ | `iaas/v2api` |
| IaaS / network | `NetworkInterface` | ✅ | `iaas/v2api` |
| IaaS / network | `Subnet` | 🟡 | `iaas/v2api` |
| IaaS / network | `Router` | 🟡 | `iaas/v2api` |
| IaaS / network | `SecurityGroup` | ✅ | `iaas/v2api` |
| IaaS / network | `FloatingIP` | 🟡 | `iaas/v2api` |
| Object Storage | `ObjectStorageBucket` | 🟡 | `services/objectstorage` |
| Object Storage | `ObjectStorageObject` | 🟡 | `services/objectstorage` |
| SKE | `SKECluster` | 🟡 | `services/ske` |
| PostgresFlex | `PostgresFlexInstance` | 🟡 | `services/postgresflex` |
| MongoDBFlex | `MongoDBFlexInstance` | 🟡 | `services/mongodbflex` |
| Redis | `RedisInstance` | 🟡 | `services/redis` |
| OpenSearch | `OpenSearchInstance` | 🟡 | `services/opensearch` |
| RabbitMQ | `RabbitMQInstance` | 🟡 | `services/rabbitmq` |
| LoadBalancer | `LoadBalancer` | 🟡 | `services/loadbalancer` |
| DNS | `DNSZone` | 🟡 | `services/dns` |

**5 of 20 resources fully working.** The CLI / config / auth / libnuke engine are functional; the per-resource SDK wiring lands incrementally. Pick one above and follow [`resources/compute-server.go`](resources/compute-server.go) as the reference pattern — see [Contributing](docs/contributing.md).

## Features

- Project-scoped destruction with explicit allow-list
- Dry-run by default; real deletion requires `--no-dry-run`
- libnuke config schema: `includes`/`excludes`/`filters`/`presets`/`blocklist`
- Service-account key auth (STACKIT standard)
- Multi-region in a single run
- Dependency-aware deletion order
- Distroless multi-arch container, Cosign-signed
- Signed binaries + SBOMs

## Install

```bash
# Pre-built binary
VERSION=v0.1.0
curl -L "https://github.com/qaiser42/stackit-nuke/releases/download/${VERSION}/stackit-nuke-${VERSION}-linux-amd64.tar.gz" \
  | tar xz -C /usr/local/bin stackit-nuke

# Container
docker pull ghcr.io/qaiser42/stackit-nuke:latest

# From source
go install github.com/qaiser42/stackit-nuke@latest
```

## Quick start

```yaml
# config.yaml
regions: [eu01]
project-ids:
  - 00000000-0000-0000-0000-000000000000
auth:
  service-account-key-path: ~/.stackit/sa-key.json
```

```bash
stackit-nuke run --config config.yaml                 # dry run
stackit-nuke run --config config.yaml --no-dry-run    # real
```

After a real run, you get a summary of everything deleted:

```
                  _ _.-'`-._
              .-'`           '-.
            .'      _ . - = - . _   '.
           /    .-'             '-.   \
          /   .'                   '.  \
         |   /                       \  |
         |  |          BOOM           | |
         |  |                         | |
          \  '.                     .'  /
           \   '-.                .-'  /
            '.    '-._        _.-'   .'
              '-.     `'----'`    .-'
                 `'-..__________..-'`

9 resource(s) nuked:

  ComputeServer
    - dev-server-0 (2bd1a04d)
    - dev-server-1 (6adb84a9)
  ComputeVolume
    - stackit-nuke-dev-extra (bb668c22)
  Network
    - dev-net (6f452480)
  NetworkInterface
    - dev-nic-0-b7856b3 (605f8edd)
    - dev-nic-1-cf66680 (3958de73)
```

Full docs: <https://qaiser42.io/stackit-nuke>

## Concepts

- **Project allow-list** — `project-ids` is the universe of what may be nuked. The CLI may *narrow* it (`--project-id`); it cannot widen. Without it, the tool refuses to start.
- **Blocklist** — `blocklist` is a hard veto. If any blocked ID also appears in `project-ids`, startup fails. Belt-and-suspenders against fat-fingering.
- **Dry-run default** — `run` lists what would be deleted. Real deletion requires `--no-dry-run` *and* (unless `--force`) typing the project ID back at a prompt.
- **Filters mark resources as ineligible** — counter-intuitive: a filter is a *keep-list*, not a kill-list. Anything matched by a filter survives; everything else gets deleted. Filter by `Name`, `tag:*`, etc.
- **Presets** — reusable named filter sets, attached to accounts via `presets: [...]`. Same shape as filters; just deduped at one site.
- **Resource type include/exclude** — `resource-types.includes` narrows the registered set; `excludes` always wins. Empty `includes` means all registered types.
- **Scopes** — every resource is `ProjectScope` today (region-aware, scoped to one STACKIT project). Engine also supports an `OrganizationScope` we have not yet used.
- **Dependency order** — each `Resource` declares `DependsOn`. libnuke topologically sorts and only deletes a resource once everything it depends on has finished. E.g. `Network` depends on `NetworkInterface` so NICs detach first.
- **Settings & feature-flags** — per-resource toggles (e.g. `EmptyBeforeDelete: true` for buckets) live under `settings:`. Engine behaviors (`wait-on-dependencies`, `filter-groups`) live under `feature-flags:`.

## How it works (101)

`stackit-nuke` is a thin CLI shell over [`libnuke`](https://github.com/ekristen/libnuke). We write the STACKIT-specific bits; libnuke does the engine work.

### Boot

`main.go` blank-imports `pkg/commands/...` and `resources/...`. Their `init()` functions:

- register CLI subcommands (`run`, `resource-types`)
- register **20 resource types** with `libnuke/pkg/registry` — each entry pairs a `Lister` (discover) with a `Resource` (delete) and optional `DependsOn` (ordering)

No reflection, no plugin loader — pure compile-time wiring.

### `run` command — what we wrote

[`pkg/commands/run/command.go`](pkg/commands/run/command.go) is the only real glue:

1. **Load config** — `config.New` parses libnuke schema (filters/presets/blocklist) plus our `project-ids` / `regions` / `auth` extension.
2. **Load credentials** — `stackit.LoadCredentials` returns a STACKIT SDK `*config.Configuration` from key path or token.
3. **Enforce allow-list** — `--project-id` may *narrow* the config list, never widen.
4. **Build engine** — `libnuke.New(params, filters, settings)`; register the typed-confirm prompt.
5. **Resolve resource types** — `types.ResolveResourceTypes` does set arithmetic over registered names ∩ includes \ excludes.
6. **Register one scanner per `(project × region)`** — each scanner carries `*stackit.ListerOpts` (project, region, credentials) which every Lister receives.
7. `n.Run(ctx)` — hand off to libnuke.

That's the whole CLI. Everything below `n.Run` is engine.

### What libnuke does inside `n.Run`

```
Validate → Prompt → Scan → Filter → (dry-run? print : delete-loop)
                                          │
                                          ├─ topological sort by DependsOn
                                          ├─ Resource.Remove(ctx)
                                          ├─ retry "waiting" items (dependency)
                                          └─ surface failures
```

| Concept | Owner |
|---|---|
| Registry, scanner, queue, dependency sort, retries, filters, dry-run | **libnuke** |
| STACKIT auth, project allow-list, per-resource SDK calls | **us** |

### Where our code plugs in

libnuke calls into our code at exactly four interfaces:

1. `registry.Register(...)` — `init()` in each `resources/*.go`
2. `Lister.List(ctx, opts) ([]resource.Resource, error)` — discovery
3. `Resource.Remove(ctx) error` — deletion
4. `RegisterPrompt(fn)` — typed-confirm

That's why the scaffold is small: 19 thin SDK adapters + one CLI wiring file. `aws-nuke`, `azure-nuke`, `gcp-nuke` are built the same way.

### Concrete trace — `ComputeServer`

```
$ stackit-nuke run --config config.yaml --no-dry-run
  │
main.go  →  cli.App.Run
  │
pkg/commands/run/command.go execute()
  ├─ config.New          → libnuke config + STACKIT fields
  ├─ stackit.LoadCredentials
  ├─ libnuke.New(params, filters, settings)
  ├─ n.RegisterScanner(ProjectScope, scanner{Opts: ListerOpts{...}})
  ├─ n.RegisterPrompt(stackit.Prompt.Prompt)
  └─ n.Run(ctx)
       │ (libnuke internals: scan)
       ↓
resources/compute-server.go  ComputeServerLister.List(ctx, opts)
  ├─ iaasv2.NewAPIClient(stackitConfigOpts(opts)...)
  └─ client.DefaultAPI.ListServers(ctx, ProjectID, Region).Execute()
       ↓ for each server → &ComputeServer{...}
       │ (libnuke internals: filter, sort, delete)
       ↓
ComputeServer.Remove(ctx)
  └─ client.DefaultAPI.DeleteServer(ctx, ProjectID, Region, ID).Execute()
```

Implementing the next 18 resources = copy this pattern, swap the SDK package.

## Development

```bash
make build        # builds ./stackit-nuke
make test         # go test -race -cover ./...
make lint         # golangci-lint
make snapshot     # goreleaser --snapshot
make docs-serve   # mkdocs at localhost:8000
```

Requires Go 1.25+.

### Throwaway test infrastructure

[`dev-infra/`](dev-infra/) is a Pulumi project ([`@stackitcloud/pulumi-stackit`](https://www.pulumi.com/registry/packages/stackit/)) that spins up a small STACKIT footprint (network, NICs, servers, volume) you can repeatedly create + nuke + recreate while developing new resource implementations.

```bash
cd dev-infra && go mod download && pulumi up
cd .. && ./stackit-nuke run --config dev-infra/nuke.yaml --no-dry-run
```

See [`dev-infra/README.md`](dev-infra/README.md).

## License

[MIT](LICENSE)
