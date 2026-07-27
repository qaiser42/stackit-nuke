# FloatingIP

STACKIT public IP (the IaaS API calls these "public IPs"). Listed and deleted
via `stackit-sdk-go/services/iaas/v2api`.

Endpoints used:

- `GET /v2/projects/{projectId}/regions/{region}/public-ips`
- `DELETE /v2/projects/{projectId}/regions/{region}/public-ips/{publicIpId}`

## Behavior

- A public IP still associated with a network interface may fail to delete on
  the first attempt; once the owning NIC or server is removed, the retry
  within the same run succeeds.

## Properties

| Name | Description |
|------|-------------|
| `OrganizationID` | parent org |
| `ProjectID` | STACKIT project ID |
| `Region` | STACKIT region |
| `ID` | public IP UUID |
| `IP` | the IPv4 address |
| `NetworkInterface` | ID of the associated network interface, empty if unattached |
| `Labels` | resource labels |
