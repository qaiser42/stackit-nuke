# dev-infra

Throwaway STACKIT infrastructure for testing `stackit-nuke` locally.

Uses the **official Pulumi STACKIT provider** ([`@stackitcloud/pulumi-stackit`](https://www.pulumi.com/registry/packages/stackit/)) — currently alpha, IaaS + Resource Manager only.

## What it deploys

- 1 non-routed network (`10.42.0.0/24`)
- N network interfaces (one per server)
- N servers running Ubuntu 24.04 on a 16 GB boot volume (default N = 2, `g2i.1`)
- 1 standalone 8 GB volume

All resources are labelled `managed-by=pulumi`, `purpose=stackit-nuke-dev` so you can filter on those.

## Prerequisites

- Pulumi CLI (`brew install pulumi`)
- Node.js 20+
- A STACKIT service-account key — same one `stackit-nuke` uses

## One-time setup

```bash
cd dev-infra
npm install
pulumi login --local                      # or your usual backend
pulumi stack init dev
cp Pulumi.dev.yaml.example Pulumi.dev.yaml
$EDITOR Pulumi.dev.yaml                   # paste your project ID
export STACKIT_SERVICE_ACCOUNT_KEY_PATH=~/.stackit/sa.json
```

## Round trip

```bash
# Deploy test infra
pulumi up

# Confirm stackit-nuke sees it (dry run)
cd .. && ./stackit-nuke run --config examples/compute-only.yaml

# Nuke it
./stackit-nuke run --config examples/compute-only.yaml --no-dry-run

# Verify the nuker did its job — Pulumi will show the servers as "missing"
cd dev-infra && pulumi refresh

# Clean residual state (volumes, network — stubs in stackit-nuke don't delete these yet)
pulumi destroy
```

## Notes

- The provider auto-detects auth from `STACKIT_SERVICE_ACCOUNT_KEY_PATH`. No provider block needed.
- The default image ID is Ubuntu 24.04 in `eu01`. Change `imageId` in `Pulumi.dev.yaml` for other regions.
- `pulumi refresh` after a nuke run is the easiest way to confirm `stackit-nuke` deleted what you wanted — Pulumi will show the deleted servers as gone from the cloud but still in state, then `destroy` cleans up.
- Cost: at `g2i.1` × 2 + 32 GB block + 8 GB block, well under €1/day. Still — don't leave it running.
