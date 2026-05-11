# Authentication

`stackit-nuke` authenticates with the STACKIT API the same way the official STACKIT CLI does — via a **service-account key**.

## Create a service account

1. STACKIT Portal → your project → **Service Accounts** → create.
2. Grant the service account the roles needed for the resources you want to delete (typically `project.owner` for full destruction).
3. Generate a service-account key (JSON). Store the file securely — anyone with it can delete your resources.

## Provide the key to stackit-nuke

In order of precedence:

1. CLI flag: `--auth-file /path/to/sa-key.json`
2. Environment variable: `STACKIT_SERVICE_ACCOUNT_KEY_PATH=/path/to/sa-key.json`
3. Config file:
   ```yaml
   auth:
     service-account-key-path: ~/.stackit/sa-key.json
   ```

If the key references an external private key, point to it with `--private-key-file` or `STACKIT_PRIVATE_KEY_PATH`.

## Token-based auth (CI shortcut)

For ephemeral environments you can skip the key file and pass a pre-issued bearer token:

```bash
export STACKIT_SERVICE_ACCOUNT_TOKEN=eyJhbGciOi...
stackit-nuke run --config config.yaml
```
