# Quick Start

## 1. Write a config

Save as `config.yaml`:

```yaml
regions:
  - eu01

project-ids:
  - 00000000-0000-0000-0000-000000000000

auth:
  service-account-key-path: ~/.stackit/sa-key.json

resource-types:
  excludes:
    - DNSZone   # keep DNS zones; delete everything else

accounts:
  "00000000-0000-0000-0000-000000000000":
    filters:
      ComputeServer:
        - property: Name
          value: "keep-*"
          type: glob
```

## 2. Dry run

```bash
stackit-nuke run --config config.yaml
```

You will see the plan — what would be deleted — but nothing happens.

## 3. Delete for real

```bash
stackit-nuke run --config config.yaml --no-dry-run
```

You will be prompted to type the project ID. Type it to proceed; anything else aborts.

## 4. Non-interactive (CI)

```bash
stackit-nuke run \
  --config config.yaml \
  --no-dry-run \
  --no-prompt \
  --prompt-delay 30
```

`--prompt-delay 30` gives operators a 30-second window to Ctrl-C the run before deletion begins.
