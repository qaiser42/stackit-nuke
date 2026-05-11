# Contributing

Thanks for considering a contribution to `stackit-nuke`!

The full contributor guide lives at [docs/contributing.md](docs/contributing.md). This file is the GitHub-discoverable entry point and points you there.

## TL;DR

```bash
git clone https://github.com/qaiser42/stackit-nuke
cd stackit-nuke
make build
make test
```

Requires Go 1.25+.

## Adding a new resource

Use [`resources/compute-server.go`](resources/compute-server.go) as the reference. The pattern:

1. `init()` calls `registry.Register` with name + scope + resource + lister.
2. The struct embeds `*BaseResource` and exposes the deletable identifier plus useful filterable properties.
3. `Remove(ctx) error` calls the appropriate STACKIT SDK delete method.
4. `<Name>Lister.List(ctx, opts) ([]resource.Resource, error)` constructs the SDK client from `opts.(*stackit.ListerOpts).Credentials` and paginates.

When your resource lister becomes real, add the resource name to `realImpls` in `resources/registry_test.go` so the stub-emptiness test skips it, and add a per-resource test mirroring `compute-server_test.go`.

## Code review

PRs need at least one approval. CI must be green. We squash on merge.

## Code of Conduct

This project follows the [Contributor Covenant](CODE_OF_CONDUCT.md). Be kind.

## Reporting security issues

See [SECURITY.md](SECURITY.md) — please **do not** open a public GitHub issue.
