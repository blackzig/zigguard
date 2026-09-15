# ZigGuard

> Universal governance layer for AI coding agents.

ZigGuard is an open-core project for defining engineering policies once and compiling them into agent-specific instructions and enforcement hooks for AI-assisted software development.

> **Status:** pre-alpha / specification phase. The policy format is being designed before the CLI and enforcement engine are implemented.

## Why ZigGuard?

Teams increasingly use more than one AI coding agent. Each tool has its own instruction format, project files, hooks, and conventions. That creates duplicated configuration and inconsistent behavior.

ZigGuard aims to provide one source of truth:

```
zigguard.yml
     |
     +--> Claude adapter
     +--> Codex adapter
     +--> Cursor adapter
     +--> GitHub Copilot adapter
     +--> local/CI policy checks
```

## MVP 0.1

The first milestone is intentionally narrow:

1. Define a versioned `zigguard.yml` policy.
2. Validate it against a published schema.
3. Compile supported policy rules into agent-specific configuration.
4. Keep generated output deterministic and reviewable.
5. Provide a local `check` path for rules that can be enforced mechanically.

Planned initial targets:

- Claude Code
- OpenAI Codex
- Cursor
- GitHub Copilot

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

## Design principles

- **Policy first** — the canonical policy is tool-neutral.
- **Deterministic output** — the same policy and version must generate the same result.
- **Least privilege** — generated rules should not grant capabilities unnecessarily.
- **No silent weakening** — unsupported `block` or approval rules must fail loudly.
- **Human reviewable** — generated instructions should remain understandable in Git.
- **Offline by default** — the core should not require sending source code to a remote service.
- **Open core, optional cloud** — local policy compilation remains useful without a SaaS account.

## Repository layout

```
adapters/        Agent-specific compilation targets
docs/            Product, architecture and policy specification
examples/        Example projects and policies
rules/           Built-in rule catalog
schemas/         Machine-readable policy schemas
AGENTS.md        Instructions for AI agents contributing to ZigGuard
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
