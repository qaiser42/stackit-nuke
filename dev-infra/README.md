# dev-infra

Throwaway STACKIT infrastructure for testing `stackit-nuke` locally.

Uses the **official Pulumi STACKIT provider** ([`pulumi-stackit`](https://www.pulumi.com/registry/packages/stackit/)) via its Go SDK.

## What it deploys

- 1 non-routed network (`10.42.0.0/24`)
- N network interfaces (one per server)
- N servers running Ubuntu 24.04 on a 16 GB boot volume (default N = 2, `g2i.1`)
- 1 standalone 8 GB volume
- 1 SSH keypair `stackit-nuke-dev` (owned by the service account, not the project)
- 1 MariaDB instance (`stackit-mariadb-1.4.10-single`, 10.11)
- 1 Object Storage bucket, plus a credentials group + S3 access key so you can put an object in it

IaaS resources are labelled `managed-by=pulumi`, `purpose=stackit-nuke-dev` so you can filter on those.

The provider has no volume snapshot resource; create one by hand (see below) to exercise `ComputeSnapshot`.

## Prerequisites

- Pulumi CLI (`brew install pulumi`)
- Go 1.23+
- A STACKIT service-account key — same one `stackit-nuke` uses

## One-time setup

```bash
cd dev-infra
go mod download
pulumi login --local                      # or your usual backend
pulumi stack init dev
cp Pulumi.dev.yaml.example Pulumi.dev.yaml
$EDITOR Pulumi.dev.yaml                   # paste your project ID
cp nuke.yaml.example nuke.yaml
$EDITOR nuke.yaml                         # paste same project ID + sa-key path
export STACKIT_SERVICE_ACCOUNT_KEY_PATH=~/.stackit/sa.json
```

## Round trip

```bash
# Deploy test infra
pulumi up

# Optional extras the provider can't create: a snapshot of the standalone
# volume, and an object in the bucket (bucket deletion must remove it)
stackit volume snapshot create -y --project-id <project-id> \
  --volume-id "$(pulumi stack output volumeId)" --name stackit-nuke-dev-snap
AWS_ACCESS_KEY_ID="$(pulumi stack output --show-secrets s3AccessKey)" \
AWS_SECRET_ACCESS_KEY="$(pulumi stack output --show-secrets s3SecretAccessKey)" \
  aws s3 cp README.md "s3://$(pulumi stack output bucketName)/README.md" \
  --endpoint-url https://object.storage.eu01.onstackit.cloud --region eu01

# Confirm stackit-nuke sees it (dry run)
cd .. && ./stackit-nuke run --config dev-infra/nuke.yaml

# Nuke it
./stackit-nuke run --config dev-infra/nuke.yaml --no-dry-run

# Verify the nuker did its job — Pulumi will show resources as "missing"
cd dev-infra && pulumi refresh

# Clear Pulumi state — the cloud resources are already gone
pulumi destroy
```

## Notes

- The provider auto-detects auth from `STACKIT_SERVICE_ACCOUNT_KEY_PATH`. No provider block needed.
- The default image ID is Ubuntu 24.04 in `eu01`. Change `imageId` in `Pulumi.dev.yaml` for other regions.
- `pulumi refresh` after a nuke run is the easiest way to confirm `stackit-nuke` deleted what you wanted — Pulumi will show the deleted servers as gone from the cloud but still in state, then `destroy` cleans up.
- `ComputeKeypair` deletes every keypair of the service account. `nuke.yaml.example` filters it down to `stackit-nuke-dev*`; keep that filter if the account owns other keypairs.
- The Object Storage credentials group and access key are not nuked (no resource type); `pulumi destroy` removes them.
- Cost: `g2i.1` × 2 + 32 GB block + 8 GB block + the smallest MariaDB plan. Don't leave it running.
