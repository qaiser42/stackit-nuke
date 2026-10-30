# MariaDBInstance

STACKIT MariaDB instance. Listed and deleted via `stackit-sdk-go/services/mariadb/v2api`.

Endpoints used:

- `GET /v2/projects/{projectId}/regions/{regionId}/instances`
- `DELETE /v2/projects/{projectId}/regions/{regionId}/instances/{instanceId}`

## Behavior

- Instances are regional: each configured region is scanned separately.
- Deletion is asynchronous: the instance stays listed with status `deleting`
  until it is gone, which matches libnuke's wait-for-removal polling.
  Instances already in `deleting` (e.g. from an interrupted run) are not
  deleted a second time.
- Instances whose last operation is a succeeded delete are skipped by the
  lister.
- Credentials and backups belong to the instance and are deleted together
  with it; they are not separate resource types.

## Properties

| Name | Description |
|------|-------------|
| `OrganizationID` | parent org |
| `ProjectID` | STACKIT project ID |
| `Region` | STACKIT region |
| `ID` | instance UUID |
| `Name` | instance display name |
| `Plan` | plan name, e.g. `stackit-mariadb-1.4.10-single` |
| `Version` | MariaDB offering version |
| `Status` | instance status, e.g. `active` |
