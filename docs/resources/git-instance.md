# GitInstance

STACKIT Git instance. Listed and deleted via `stackit-sdk-go/services/git/v1betaapi`.

Endpoints used:

- `GET /v1beta/projects/{projectId}/instances`
- `DELETE /v1beta/projects/{projectId}/instances/{instanceId}`

## Behavior

- STACKIT Git is a global service: instances are not bound to a region. The
  lister only returns instances during the first configured region's scan, so
  each instance shows up exactly once per run.
- Deletion is asynchronous: the instance stays listed in state `Deleting`
  until it is gone, which matches libnuke's wait-for-removal polling.
- Repositories, runners, and authentications hosted on the instance are
  deleted together with it; they are not separate resource types.

## Properties

| Name | Description |
|------|-------------|
| `OrganizationID` | parent org |
| `ProjectID` | STACKIT project ID |
| `Region` | region of the scan that listed the instance (always the first configured region) |
| `ID` | instance UUID |
| `Name` | instance display name |
| `URL` | URL for reaching the instance |
| `Flavor` | instance flavor |
| `Version` | deployed STACKIT Git version |
| `State` | instance state, e.g. `Ready` |
