# Changelog

All notable changes to this project will be documented here. Format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/); the project adheres
to [Semantic Versioning](https://semver.org/).

## [Unreleased]

### Added
- Initial scaffold: CLI, libnuke wiring, STACKIT auth, config loader
- Resource registrations (placeholder Listers): ComputeServer, ComputeVolume,
  ComputeSnapshot, ComputeKeypair, Network, Subnet, Router, SecurityGroup,
  FloatingIP, ObjectStorageBucket, ObjectStorageObject, SKECluster,
  PostgresFlexInstance, MongoDBFlexInstance, RedisInstance,
  OpenSearchInstance, RabbitMQInstance, LoadBalancer, DNSZone
- CI: lint, test, build (GitHub Actions, Go 1.25)
- Release: GoReleaser cross-compile, ko distroless multi-arch, Cosign signing
- Docs: MkDocs Material site published to GitHub Pages
