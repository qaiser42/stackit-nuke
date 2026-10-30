# ComputeSnapshot

STACKIT IaaS volume snapshot. Listed and deleted via `stackit-sdk-go/services/iaas/v2api`.

Endpoints used:

- `GET /v2/projects/{projectId}/regions/{region}/snapshots`
- `DELETE /v2/projects/{projectId}/regions/{region}/snapshots/{snapshotId}`

## Behavior

- Snapshots are regional: each configured region is scanned separately.
- Deletion is asynchronous: the snapshot stays listed with status `DELETING`
  until it is gone, which matches libnuke's wait-for-removal polling.
  Snapshots already in `DELETING` are not deleted a second time.
- Snapshots with status `DELETED` are skipped by the lister.
- `ComputeVolume` depends on `ComputeSnapshot`, so snapshots are removed
  before the volumes they were taken from.

## Properties

| Name | Description |
|------|-------------|
| `OrganizationID` | parent org |
| `ProjectID` | STACKIT project ID |
| `Region` | STACKIT region |
| `ID` | snapshot UUID |
| `Name` | snapshot name |
| `VolumeID` | UUID of the source volume |
| `SizeGB` | snapshot size in GB |
| `Status` | snapshot status, e.g. `AVAILABLE` |
| `CreatedAt` | creation timestamp |
| `Labels` | snapshot labels |
