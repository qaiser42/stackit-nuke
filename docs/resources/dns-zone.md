# DNSZone

STACKIT DNS zone. Listed and deleted via `stackit-sdk-go/services/dns/v1api`.

Endpoints used:

- `GET /v1/projects/{projectId}/zones`
- `DELETE /v1/projects/{projectId}/zones/{zoneId}`

## Behavior

- DNS is a global service: zones are not bound to a region. The lister only
  returns zones during the first configured region's scan, so each zone shows
  up exactly once per run.
- Deleting a zone is a soft delete: the zone stays listed with state
  `DELETE_SUCCEEDED` during the retention window. The lister filters those
  out, so the zone counts as removed once the delete has landed.
- Record sets are deleted together with the zone; they are not a separate
  resource type.

## Properties

| Name | Description |
|------|-------------|
| `OrganizationID` | parent org |
| `ProjectID` | STACKIT project ID |
| `Region` | region of the scan that listed the zone (always the first configured region) |
| `ID` | zone UUID |
| `Name` | zone display name |
| `DNSName` | fully qualified domain name of the zone |
| `State` | zone state, e.g. `CREATE_SUCCEEDED` |
