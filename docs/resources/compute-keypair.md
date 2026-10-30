# ComputeKeypair

STACKIT IaaS SSH keypair. Listed and deleted via `stackit-sdk-go/services/iaas/v2api`.

Endpoints used:

- `GET /v2/keypairs`
- `DELETE /v2/keypairs/{keypairName}`

## Behavior

- Keypairs are neither regional nor project-scoped: they belong to the
  authenticated identity (user or service account). The lister only returns
  keypairs during the first configured project's first region scan, so each
  keypair shows up exactly once per run.
- Because of that, a run deletes **every** keypair owned by the credentials
  in use, including keypairs used for servers in other projects. Use a
  project-specific service account, or exclude `ComputeKeypair` / filter by
  `Name` or `Labels`, to keep shared keypairs.
- Deleting a keypair does not affect running servers that were created with it.

## Properties

| Name | Description |
|------|-------------|
| `OrganizationID` | parent org |
| `ProjectID` | first configured project ID (keypairs are not project-bound) |
| `Region` | first configured region (keypairs are not regional) |
| `Name` | keypair name (unique per identity) |
| `Fingerprint` | MD5 fingerprint of the public key |
| `CreatedAt` | creation timestamp |
| `Labels` | keypair labels |
