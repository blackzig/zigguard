# ZigGuard

> Universal governance layer for AI coding agents.

ZigGuard is an open-core project for defining engineering policies once and compiling them into safe, reviewable instruction and enforcement surfaces for AI-assisted software development.

> **Status:** pre-alpha. ZigGuard includes a working Go CLI for policy initialization, validation, deterministic instruction compilation, and the first mechanically enforced drift check.

## Why ZigGuard?

Teams increasingly use more than one AI coding agent. Each tool has its own instruction formats and capabilities. Maintaining separate copies creates duplicated configuration, inconsistent behavior, and policy drift.

ZigGuard provides one source of truth:

```
zigguard.yml
     |
     +--> validate
     |
     +--> compile
           |
           +--> AGENTS.md
           +--> CLAUDE.md
           +--> .zigguard/manifest.json

     +--> check
           |
           +--> missing artifacts
           +--> policy drift
           +--> unmanaged conflicts

Future:
     +--> hooks and additional enforceable controls
```

## MVP 0.1

The first milestone is intentionally narrow:

- versioned `zigguard.yml` policy;
- strict YAML parsing and semantic validation;
- `zigguard init`;
- `zigguard validate`;
- `zigguard compile`;
- `zigguard check` for generated-artifact integrity;
- deterministic generated output;
- managed-section merging that preserves existing human-authored agent instructions;
- shared `AGENTS.md` output for Codex, Cursor, and GitHub Copilot;
- `CLAUDE.md` output for Claude Code;
- generated target manifest;
- automated Go tests and CI.

## Build

Requires Go 1.23+ for development.

```bash
go test ./...
go build -o bin/zigguard ./cmd/zigguard
```

## Quick start

```bash
./bin/zigguard init
./bin/zigguard validate
./bin/zigguard compile --dry-run
./bin/zigguard compile
./bin/zigguard check
```

If `AGENTS.md` or `CLAUDE.md` already contains human-authored instructions, ZigGuard preserves them and manages only the section between `<!-- zigguard:managed:start -->` and `<!-- zigguard:managed:end -->`. Malformed or duplicate managed markers fail safely instead of triggering a destructive rewrite.

See [docs/cli.md](docs/cli.md).

## Example policy

```yaml
version: "0.1"

project:
  name: billing-api
  stack:
    - java
    - spring-boot
    - oracle

targets:
  - claude
  - codex
  - cursor
  - copilot

rules:
  scope:
    refactor_outside_request: block
    public_api_changes: require_approval

  security:
    secrets: block
    dynamic_sql: block

  quality:
    tests_required: true
    lint_required: true

  git:
    force_push_protected_branches: block
```

See [zigguard.example.yml](zigguard.example.yml), [docs/policy-spec.md](docs/policy-spec.md), and the complete [Java + Spring Boot + Oracle fixture](examples/java-spring-oracle/README.md).

## Important: instructions are not enforcement

MVP 0.1 generates **instruction context**. A rule rendered as `BLOCK` is a strong instruction to the agent, but it is not equivalent to branch protection, a CI gate, a hook, or a security scanner.

ZigGuard will add mechanical enforcement paths where they are technically possible. It will not market natural-language guidance as a hard security boundary.

`zigguard check` is the first mechanical control: it fails when required managed sections/files are missing, when ZigGuard-owned content drifts from `zigguard.yml`, or when managed markers are malformed/ambiguous. Human-authored content outside ZigGuard's section is intentionally ignored. The command does not yet enforce the semantic rules inside agent instructions.

See the current [target capability matrix](docs/capability-matrix.md).

## Design principles

- **Policy first** — the canonical policy is tool-neutral.
- **One policy, minimum duplication** — reuse shared standards where possible.
- **Deterministic output** — the same policy and version generate the same result.
- **Least privilege** — do not grant capabilities unnecessarily.
- **No silent weakening** — capability level is explicit.
- **Human reviewable** — generated instructions remain understandable in Git.
- **Safe writes** — unmanaged files are not overwritten silently.
- **Offline by default** — the core does not require a SaaS account.
- **Open core, optional cloud** — local compilation remains useful without hosted services.

## Repository layout

```
cmd/             CLI entry point
internal/        policy parser, compiler, checker, output safety, CLI
adapters/        target behavior and capability documentation
docs/            product, architecture, CLI and policy specification
examples/        example projects and future fixtures
rules/           built-in rule catalog
schemas/         machine-readable policy schemas
AGENTS.md        instructions for AI agents contributing to ZigGuard
zigguard.example.yml
```

## GitHub Action

ZigGuard includes an official composite GitHub Action for pull-request governance checks. Once the first prerelease tag is published, repositories can pin that tag and run `zigguard check` directly in CI.

See [docs/github-action.md](docs/github-action.md).

## Releases

Cross-platform release automation is prepared for Linux, macOS, and Windows on amd64/arm64. No public binary release is published yet. See [docs/releasing.md](docs/releasing.md).

## Roadmap

See [ROADMAP.md](ROADMAP.md).

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## Security

See [SECURITY.md](SECURITY.md).

## License

Licensed under the [Apache License 2.0](LICENSE).
