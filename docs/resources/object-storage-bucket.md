# ObjectStorageBucket

STACKIT Object Storage bucket. Listed and deleted via `stackit-sdk-go/services/objectstorage/v2api`.

Endpoints used:

- `GET /v2/project/{projectId}/regions/{region}/buckets`
- `DELETE /v2/project/{projectId}/regions/{region}/bucket/{bucketName}`

## Behavior

- Buckets are regional: each configured region is scanned separately.
- Deleting a bucket removes every object it contains. Objects are not a
  separate resource type; the management API has no object endpoints.
- Projects without Object Storage enabled answer the list call with `404`;
  the lister treats that as "no buckets".
- Access keys and credentials groups are not deleted.

## Properties

| Name | Description |
|------|-------------|
| `OrganizationID` | parent org |
| `ProjectID` | STACKIT project ID |
| `Region` | STACKIT region |
| `Name` | bucket name |
| `ObjectLockEnabled` | whether S3 Object Lock is enabled for the bucket |
| `URL` | virtual-hosted-style bucket URL |
