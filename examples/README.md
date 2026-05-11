# Examples

Drop-in `config.yaml` files for common scenarios. Copy, edit project IDs and auth path, then:

```bash
stackit-nuke run --config examples/<file>.yaml             # dry-run
stackit-nuke run --config examples/<file>.yaml --no-dry-run
```

| File | What |
|------|------|
| [`minimal.yaml`](minimal.yaml) | Smallest valid config — one project, one region |
| [`full.yaml`](full.yaml) | Annotated reference covering every supported key |
| [`compute-only.yaml`](compute-only.yaml) | Nuke only IaaS compute (servers, volumes, snapshots, keypairs); skip everything else |
| [`keep-prod.yaml`](keep-prod.yaml) | Production-protective: blocklist + tag-based filters |
| [`ci.yaml`](ci.yaml) | Non-interactive run, env-driven auth, suitable for scheduled CI cleanup |

## Conventions in these files

- Project ID `11111111-…` is the target.
- Project ID `99999999-…` is on the blocklist (production placeholder).
- Auth points at `~/.stackit/sa-key.json` — override per env.
