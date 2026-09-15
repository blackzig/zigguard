# ZigGuard CLI

The MVP CLI is implemented in Go so ZigGuard can be distributed as a single executable for common developer platforms.

## Commands

### zigguard init

Creates a starter `zigguard.yml`.

The command refuses to replace an existing policy unless `--force` is explicitly supplied.

### zigguard validate

Strictly parses the policy and checks:

- policy version;
- required project metadata;
- supported targets;
- duplicate targets/stack values;
- supported rule actions;
- presence of at least one configured rule;
- unknown YAML fields.

### zigguard compile

Compiles a valid policy into managed agent instruction surfaces.

The MVP intentionally prefers a shared instruction surface when tools already support the same convention:

- Codex → `AGENTS.md`
- Cursor → `AGENTS.md`
- GitHub Copilot → `AGENTS.md`
- Claude Code → `CLAUDE.md`

When Claude Code and any AGENTS-compatible target are selected together, ZigGuard generates the canonical policy in `AGENTS.md` and a small `CLAUDE.md` that imports `@AGENTS.md`. This avoids maintaining two independent copies of the policy.

The compiler also writes `.zigguard/manifest.json`, which records the target-to-surface bindings.

## Safety

Generated files carry a ZigGuard marker. The compiler may update files it previously generated, but it refuses to overwrite an unmanaged `AGENTS.md` or `CLAUDE.md` unless the user explicitly passes `--force`.

This is deliberately conservative. A later milestone will support managed sections/merge strategies for repositories that already maintain instruction files.

## Instruction vs enforcement

The initial surfaces are **instruction context**, not hard security controls.

A `block` rule is rendered as an explicit mandatory instruction, but ZigGuard does not claim that a natural-language agent instruction is technically equivalent to branch protection, a hook, a CI gate, or a security scanner.

Future local enforcement and hook adapters will complement instruction surfaces.

## Examples

```bash
zigguard init
zigguard validate
zigguard compile --dry-run
zigguard compile
```

Use a different policy or repository root with:

```bash
zigguard validate --file config/zigguard.yml
zigguard compile --file config/zigguard.yml --root /path/to/repo
```
