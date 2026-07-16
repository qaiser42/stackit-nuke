# Releases

Releases are produced by the [`release`](https://github.com/qaiser42/stackit-nuke/actions/workflows/release.yml) workflow when a `v*` tag is pushed.

## Artifacts

- **Binaries**: linux/darwin/windows × amd64/arm64 on the GitHub Release page
- **Checksums**: `checksums.txt` signed with Cosign keyless (sigstore)
- **SBOMs**: per-archive SPDX JSON
- **Container image**: `ghcr.io/qaiser42/stackit-nuke:<tag>` and `:latest`, multi-arch (amd64, arm64), built with `ko` on a Chainguard distroless base, Cosign-signed

## Verifying a release

```bash
# checksum signature (sigstore bundle contains signature + certificate)
cosign verify-blob \
  --bundle checksums.txt.sigstore.json \
  --certificate-identity-regexp 'github.com/qaiser42/stackit-nuke' \
  --certificate-oidc-issuer 'https://token.actions.githubusercontent.com' \
  checksums.txt

# image
cosign verify ghcr.io/qaiser42/stackit-nuke:v0.1.0 \
  --certificate-identity-regexp 'github.com/qaiser42/stackit-nuke' \
  --certificate-oidc-issuer 'https://token.actions.githubusercontent.com'
```

## Cutting a release

```bash
git tag -a v0.1.0 -m "v0.1.0"
git push origin v0.1.0
```

The workflow handles the rest.
