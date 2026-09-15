# Changelog

All notable changes to ZigGuard will be documented in this file.

The project follows semantic versioning once the first public release is published.

## [Unreleased]

### Added

- Go CLI foundation.
- `zigguard init`, `zigguard validate`, and `zigguard compile`.
- Strict YAML parsing with semantic validation.
- Deterministic shared `AGENTS.md` generation for Codex, Cursor, and GitHub Copilot.
- `CLAUDE.md` generation for Claude Code, importing `AGENTS.md` when a shared policy surface exists.
- `.zigguard/manifest.json` target binding metadata.
- Safe generated-file ownership and explicit `--force` override.
- Unit tests and GitHub Actions CI.
- Java + Spring Boot + Oracle example with exact compiler contract fixtures.
- Target capability matrix distinguishing instruction context from enforcement.
- Cross-platform tag release workflow for Linux, macOS, and Windows on amd64/arm64.
- Initial product vision and architecture.
- Draft ZigGuard policy format v0.1.
- JSON Schema for the draft policy.
- Open-core boundary and contribution/security guidance.

### Changed

- CLI version is now injectable at build time for tagged releases.
