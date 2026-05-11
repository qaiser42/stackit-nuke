# Filtering

Filters mark resources as **ineligible for deletion**. Anything that does not match a filter is deleted.

## Shape

```yaml
accounts:
  "<project-id>":
    filters:
      <ResourceType>:
        - property: <PropertyName>   # e.g. Name, ID, Region, tag:<key>
          value: "<literal-or-pattern>"
          type: <exact|glob|regex|contains|dateOlderThan|...>
          invert: false
```

## Property names

Each resource exposes properties via the libnuke property reflection. `stackit-nuke resource-types` lists all registered types; the property names of each type are documented on the per-resource pages under [Resources](resources/index.md).

## Match types

| Type | Behavior |
|------|----------|
| `exact` (default) | string equality |
| `glob` | shell-style glob: `prod-*`, `*.example.com` |
| `regex` | RE2 |
| `contains` | substring |
| `dateOlderThan` | parse property as time, compare with duration value (`72h`) |
| `NotIn` | property must not appear in `values:` list |

## Invert

`invert: true` flips the match — the resource is **only** considered for deletion if the filter would have excluded it. Useful with `dateOlderThan` to nuke only old resources.
