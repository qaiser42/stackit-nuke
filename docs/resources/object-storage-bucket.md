# ObjectStorageBucket

STACKIT Object Storage bucket. Listed and deleted via `stackit-sdk-go/services/objectstorage/v2api`.

Endpoints used:

- `GET /v2/project/{projectId}/regions/{region}/buckets`
- `DELETE /v2/project/{projectId}/regions/{region}/bucket/{bucketName}`
- `POST`/`DELETE` `/v2/project/{projectId}/regions/{region}/credentials-group[/{groupId}]` (temporary)
- `POST`/`DELETE` `/v2/project/{projectId}/regions/{region}/access-key[/{keyId}]` (temporary)
- S3 `ListObjectVersions` and `DeleteObjects` on the bucket

## Behavior

- Buckets are regional: each configured region is scanned separately.
- The API refuses to delete a bucket that still holds objects
  (`422 bucket.not_empty`) and has no object endpoints. On that error the
  bucket is emptied over S3 first: a temporary credentials group and access
  key (expiring after one hour) are created, every object version and delete
  marker is deleted, the bucket is deleted, and the temporary credentials
  are removed again. Empty buckets are deleted without creating credentials.
- Objects are therefore not a separate resource type.
- Objects protected by Object Lock retention or legal hold cannot be
  deleted; the bucket then fails with the S3 error.
- Projects without Object Storage enabled answer the list call with `404`;
  the lister treats that as "no buckets".
- Pre-existing access keys and credentials groups are not deleted.

## Properties

| Name | Description |
|------|-------------|
| `OrganizationID` | parent org |
| `ProjectID` | STACKIT project ID |
| `Region` | STACKIT region |
| `Name` | bucket name |
| `ObjectLockEnabled` | whether S3 Object Lock is enabled for the bucket |
| `URL` | virtual-hosted-style bucket URL |
