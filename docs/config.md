# Configuration

`stackit-nuke` reuses the [libnuke](https://github.com/ekristen/libnuke) config schema and adds a few STACKIT-specific top-level keys.

## Top-level keys

| Key | Type | Required | Description |
|-----|------|----------|-------------|
| `regions` | `[]string` | yes | STACKIT regions to scan (`eu01`, `eu02`, ...). `all` disables the region filter. |
| `project-ids` | `[]string` | yes | Allow-list of STACKIT project IDs that may be nuked. CLI cannot widen this. |
| `organization-id` | `string` | no | Parent organization. Used by future organization-scoped resources. |
| `auth` | `object` | no | Credential locations. See [auth](auth.md). |
| `accounts` | `map[string]Account` | no | Per-project filters/presets/resource-types. Keyed by project ID. |
| `resource-types` | `object` | no | Global `includes` / `excludes` of resource type names. |
| `presets` | `map[string]Preset` | no | Reusable filter sets, referenced by `accounts.<id>.presets`. |
| `blocklist` | `[]string` | no | Project IDs that may **never** be nuked. Refused with an error. |
| `settings` | `map[string]map[string]any` | no | Per-resource feature toggles. |
| `feature-flags` | `map[string]any` | no | Engine flags. |

## Example

```yaml
regions: [eu01, eu02]

project-ids:
  - 11111111-1111-1111-1111-111111111111

blocklist:
  - 99999999-9999-9999-9999-999999999999  # production: never touch

auth:
  service-account-key-path: ~/.stackit/sa-key.json

resource-types:
  excludes: [DNSZone]

presets:
  protect-prod:
    filters:
      ComputeServer:
        - property: "tag:env"
          value: prod

accounts:
  "11111111-1111-1111-1111-111111111111":
    presets: [protect-prod]
    filters:
      ObjectStorageBucket:
        - property: Name
          value: "logs-*"
          type: glob
```

See [Filtering](config-filtering.md) and [Presets](config-presets.md) for the filter syntax.
