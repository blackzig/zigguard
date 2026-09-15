# ZigGuard

> Universal governance layer for AI coding agents.

ZigGuard is an open-core project for defining engineering policies once and compiling them into safe, reviewable instruction and enforcement surfaces for AI-assisted software development.

> **Status:** pre-alpha. MVP 0.1 now includes a working Go CLI for policy initialization, validation, and deterministic instruction compilation.

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

Future:
     +--> local/CI checks
     +--> hooks and enforceable controls
```

## MVP 0.1

The first milestone is intentionally narrow:

- versioned `zigguard.yml` policy;
- strict YAML parsing and semantic validation;
- `zigguard init`;
- `zigguard validate`;
- `zigguard compile`;
- deterministic generated output;
- safe refusal to overwrite unmanaged agent files;
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
```

The compiler refuses to overwrite an existing unmanaged `AGENTS.md` or `CLAUDE.md` unless `--force` is explicitly supplied.

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

See [zigguard.example.yml](zigguard.example.yml) and [docs/policy-spec.md](docs/policy-spec.md).

## Important: instructions are not enforcement

MVP 0.1 generates **instruction context**. A rule rendered as `BLOCK` is a strong instruction to the agent, but it is not equivalent to branch protection, a CI gate, a hook, or a security scanner.

ZigGuard will add mechanical enforcement paths where they are technically possible. It will not market natural-language guidance as a hard security boundary.

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
internal/        policy parser, compiler, output safety, CLI
adapters/        target behavior and capability documentation
docs/            product, architecture, CLI and policy specification
examples/        example projects and future fixtures
rules/           built-in rule catalog
schemas/         machine-readable policy schemas
AGENTS.md        instructions for AI agents contributing to ZigGuard
zigguard.example.yml
```

## Roadmap

See [ROADMAP.md](ROADMAP.md).

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## Security

See [SECURITY.md](SECURITY.md).

## License

Licensed under the [Apache License 2.0](LICENSE).
