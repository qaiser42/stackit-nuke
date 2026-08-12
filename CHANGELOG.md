# Changelog

All notable changes to this project will be documented here. Format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/); the project adheres
to [Semantic Versioning](https://semver.org/).

## [Unreleased]

## [0.0.7] - 2026-08-12

### Added
- New `GitInstance` resource type: real list + delete via STACKIT Git
  v1beta SDK. Instances listed once per run (global service); deletion is
  asynchronous and polled by libnuke until gone
- Docs page for `git-instance`

## [0.0.6] - 2026-07-27

### Added
- Real list + delete via STACKIT DNS v1 SDK for `DNSZone`: zones listed once
  per run (global service), soft-deleted zones filtered out
- Real list + delete via STACKIT IaaS v2 public-ips API for `FloatingIP`
- Real docs for `dns-zone` and `floating-ip` resource pages

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
