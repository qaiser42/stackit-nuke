# Security Policy

## Supported versions

| Version | Supported |
|---------|-----------|
| `0.x`   | ✅ — pre-1.0, security fixes ship in the next minor |
| older   | ❌ |

## Reporting a vulnerability

**Do not open a public GitHub issue for security reports.**

Please report privately via [GitHub Security Advisories](https://github.com/qaiser42/stackit-nuke/security/advisories/new).

Include:

- Affected version(s) (`stackit-nuke --version`)
- Reproduction steps or proof-of-concept
- Impact assessment (data loss, privilege escalation, info disclosure, …)

## What to expect

- Acknowledgement within **3 business days**.
- Triage + initial assessment within **7 business days**.
- A fix or mitigation plan within **30 days** for confirmed high-severity issues.

## Scope

Because `stackit-nuke` deletes infrastructure, we are particularly interested in:

- Auth/credential mishandling (leakage in logs, world-readable temp files, etc.)
- Config-allow-list bypasses (anything that lets the tool delete resources outside the configured `project-ids`)
- Dependency vulnerabilities in our supply chain

Out of scope:

- Issues in upstream STACKIT services
- Vulnerabilities in `libnuke` itself — please report those at <https://github.com/ekristen/libnuke/security>

## Continuous scanning

Every PR + push to `main` runs [Trivy](https://github.com/aquasecurity/trivy) against:

- the working tree (filesystem + Go module deps)
- IaC + GitHub workflow configs
- the published container image (post-release, on tag)

Findings show up under the repo's [Security → Code scanning](https://github.com/qaiser42/stackit-nuke/security/code-scanning) tab. PRs fail on `CRITICAL` filesystem findings; lower severities are reported but non-blocking. Accepted findings are tracked in [`.trivyignore`](.trivyignore) with an expiry and rationale.

## Release signing

Binaries and container images are signed with [Cosign](https://github.com/sigstore/cosign) keyless via GitHub OIDC. Verify before running — see [docs/releases.md](docs/releases.md).
