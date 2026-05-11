# Warning

!!! danger
    `stackit-nuke` deletes **everything** in the targeted STACKIT project that matches the configured resource types. Deletion is permanent and not recoverable.

## Safety defaults

- **Dry-run is the default.** Without `--no-dry-run`, nothing is deleted.
- **Explicit project allow-list required.** The config must list `project-ids:`. The CLI cannot widen this list, only narrow it.
- **Typed-confirmation prompt.** You must type the project ID exactly to proceed, unless you also pass `--no-prompt`.
- **Blocklist support** via libnuke `blocklist:` config — projects on the list cannot be nuked, ever.

## What you should do before running

1. Take a snapshot or backup of any data you may want to keep.
2. Run with `--log-level=debug` and **without** `--no-dry-run` first. Read the plan.
3. Use `--include` to scope an initial run to a single resource type.
4. Add unwanted-to-delete resources to your `filters:` block.
