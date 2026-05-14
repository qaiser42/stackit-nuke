# Changelog

All notable changes to this project will be documented here. Format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/); the project adheres
to [Semantic Versioning](https://semver.org/).

## [Unreleased]

## [0.0.2] - 2026-05-14

### Added
- Real list + delete via STACKIT IaaS v2 SDK for: ComputeVolume,
  Network, NetworkInterface, SecurityGroup
- New `NetworkInterface` resource type
- `dev-infra/` Pulumi (Go) stack for round-trip create + nuke testing,
  plus `dev-infra/nuke.yaml(.example)` companion config
- Post-run summary: grouped list of nuked resources + ASCII art banner
- Compact log formatter (strips libnuke property dump);
  `--log-verbose` restores full output

## [0.0.1]

### Added
- Initial scaffold: CLI, libnuke wiring, STACKIT auth, config loader
- Resource registrations (placeholder Listers): ComputeServer, ComputeVolume,
  ComputeSnapshot, ComputeKeypair, Network, Subnet, Router, SecurityGroup,
  FloatingIP, ObjectStorageBucket, ObjectStorageObject, SKECluster,
  PostgresFlexInstance, MongoDBFlexInstance, RedisInstance,
  OpenSearchInstance, RabbitMQInstance, LoadBalancer, DNSZone
- Real list + delete via STACKIT IaaS v2 SDK for ComputeServer
- CI: lint, test, build (GitHub Actions, Go 1.25)
- Release: GoReleaser cross-compile, ko distroless multi-arch, Cosign signing
- Docs: MkDocs Material site published to GitHub Pages
