# Concepts

- **Project allow-list** — `project-ids` is the universe of what may be nuked. The CLI may *narrow* it (`--project-id`); it cannot widen. Without it, the tool refuses to start.
- **Blocklist** — `blocklist` is a hard veto. If any blocked ID also appears in `project-ids`, startup fails. Belt-and-suspenders against fat-fingering.
- **Dry-run default** — `run` lists what would be deleted. Real deletion requires `--no-dry-run` *and* (unless `--force`) typing the project ID back at a prompt.
- **Filters mark resources as ineligible** — counter-intuitive: a filter is a *keep-list*, not a kill-list. Anything matched by a filter survives; everything else gets deleted. Filter by `Name`, `tag:*`, etc.
- **Presets** — reusable named filter sets, attached to accounts via `presets: [...]`. Same shape as filters; just deduped at one site.
- **Resource type include/exclude** — `resource-types.includes` narrows the registered set; `excludes` always wins. Empty `includes` means all registered types.
- **Scopes** — every resource is `ProjectScope` today (region-aware, scoped to one STACKIT project). Engine also supports an `OrganizationScope` we have not yet used.
- **Dependency order** — each `Resource` declares `DependsOn`. libnuke topologically sorts and only deletes a resource once everything it depends on has finished. E.g. `Network` depends on `NetworkInterface` so NICs detach first.
- **Settings & feature-flags** — per-resource toggles (e.g. `EmptyBeforeDelete: true` for buckets) live under `settings:`. Engine behaviors (`wait-on-dependencies`, `filter-groups`) live under `feature-flags:`.
