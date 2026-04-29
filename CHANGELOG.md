# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Multi-arch Docker builds (amd64 + arm64) via GitHub Actions
- Docker image published to GitHub Container Registry (GHCR)
- Health check endpoint in Docker container
- `.dockerignore` file for cleaner builds
- Non-root user in Docker runtime

### Changed
- Dockerfile optimized: multi-stage build with proper layer caching
- Docker runtime: copies `config.toml.example` instead of requiring local `config.toml`
- `docker-compose.yml` with health check and proper env var defaults

## [1.2.0] - 2026-04-26

### Added
- Fork do projeto original [go-samba4](https://github.com/jniltinho/go-samba4)
- Renamed project to `samba4-manager`
- Docker CI workflow with GHCR publishing
- Release CI workflow for automated binary packaging

### Fixed
- **RBAC:** `isAdminUser()` now supports DN-format group membership (e.g., `CN=Domain Admins,CN=Users,DC=domain,DC=tld`)
- **UI:** Removed `{{if .IsAdmin}}` check that hid the "+ NEW USER" button
- **Security:** StartTLS enabled by default in `config.toml.example`

## [1.1.2] - 2026-03-10

Original release by [@jniltinho](https://github.com/jniltinho).

[Unreleased]: https://github.com/clediomir/samba4-manager/compare/v1.2.0...HEAD
[1.2.0]: https://github.com/clediomir/samba4-manager/releases/tag/v1.2.0
[1.1.2]: https://github.com/jniltinho/go-samba4/releases/tag/v1.1.2
