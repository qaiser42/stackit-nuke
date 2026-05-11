# Contributing

Thanks for helping improve `stackit-nuke`!

## Local setup

```bash
git clone https://github.com/qaiser42/stackit-nuke
cd stackit-nuke
make build
./stackit-nuke --version
```

Requires Go 1.26+.

## Adding a new resource type

1. Create `resources/<service>-<thing>.go`.
2. Implement the libnuke contract — see `resources/compute-server.go` for the reference shape:
   - `init()` calling `registry.Register`
   - struct embedding `*BaseResource`
   - `Remove(ctx) error`, `Properties() types.Properties`, `String() string`
   - sibling `<Name>Lister` with `List(ctx, opts any) ([]resource.Resource, error)`
3. Use the STACKIT SDK service package that owns the resource (`github.com/stackitcloud/stackit-sdk-go/services/<svc>`). Build the client from `opts.(*stackit.ListerOpts).Credentials.SDKConfigForRegion(opts.Region)`.
4. Add a docs page under `docs/resources/` and a nav entry in `mkdocs.yml`.
5. Add a test alongside the file (table-driven, mock the SDK client).

## Tests, lint, release

```bash
make test    # go test -race -cover ./...
make lint    # golangci-lint
make snapshot # goreleaser --snapshot
```

## Commit / PR style

- Conventional Commits (`feat:`, `fix:`, `docs:`, ...).
- One resource per PR keeps reviews tractable.
- No AI attribution in commits or PR descriptions.
