# Install

## Pre-built binaries

Each release publishes binaries for `linux`, `darwin`, `windows` × `amd64`/`arm64` to [GitHub Releases](https://github.com/qaiser42/stackit-nuke/releases). Checksums and SBOMs are signed with Cosign.

```bash
VERSION=v0.1.0
OS=linux
ARCH=amd64
curl -L "https://github.com/qaiser42/stackit-nuke/releases/download/${VERSION}/stackit-nuke-${VERSION}-${OS}-${ARCH}.tar.gz" \
  | tar xz -C /usr/local/bin stackit-nuke
stackit-nuke --version
```

## Container image

Multi-arch distroless image, signed with Cosign keyless:

```bash
docker pull ghcr.io/qaiser42/stackit-nuke:latest
docker run --rm \
  -v "${HOME}/.stackit:/stackit:ro" \
  -v "${PWD}/config.yaml:/etc/stackit-nuke/config.yaml:ro" \
  ghcr.io/qaiser42/stackit-nuke:latest \
  run --config /etc/stackit-nuke/config.yaml --auth-file /stackit/sa-key.json
```

## From source

```bash
go install github.com/qaiser42/stackit-nuke@latest
```

Requires Go 1.26+.
