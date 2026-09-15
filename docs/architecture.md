# Architecture

## High-level model

```
                 +-------------------+
                 |   zigguard.yml    |
                 +---------+---------+
                           |
                    validate/normalize
                           |
                 +---------v---------+
                 |  Policy Model     |
                 +----+---------+----+
                      |         |
              compile |         | check
                      |         |
        +-------------+         +----------------+
        |                                      |
+-------v-------+                      +-------v--------+
| Agent surfaces|                      | Local checks   |
+---+---+---+---+                      +----------------+
    |   |   |   |
 Claude Codex Cursor Copilot
```

## Components

### Policy parser

Responsibilities:

- load `zigguard.yml`;
- strictly validate the supported v0.1 shape;
- normalize values into an internal policy model;
- reject unknown fields and unsupported values.

A machine-readable JSON Schema remains published as the language-neutral contract. The Go MVP performs semantic validation in code.

### Policy compiler

Responsibilities:

- convert normalized rules into deterministic agent instructions;
- prefer a shared instruction standard when several targets support it;
- maintain deterministic ordering;
- produce a visible dry-run plan;
- never label instruction context as technical enforcement.

### Shared instruction strategy

Duplicating the same policy into multiple files can create drift and can cause tools that read more than one instruction format to load redundant guidance.

For MVP 0.1:

- Codex, Cursor, and GitHub Copilot share `AGENTS.md`;
- Claude Code uses `CLAUDE.md`;
- when both families are selected, `CLAUDE.md` imports `@AGENTS.md`;
- `.zigguard/manifest.json` records which surface each selected target uses.

Target-specific formats are added only when they provide useful semantics that the shared surface cannot express.

### File ownership

Artifacts declare an ownership mode rather than relying on filename conventions:

- `managed-section` — ZigGuard owns only a delimited section inside a potentially human-maintained file;
- `managed-file` — ZigGuard owns the entire artifact.

`AGENTS.md` and `CLAUDE.md` use `managed-section`. Their generated content is enclosed by explicit start/end markers. Existing human content before or after the section is preserved.

`.zigguard/manifest.json` uses `managed-file`.

The writer:

- creates a managed section when the instruction file does not exist;
- appends a managed section to an existing human-maintained instruction file;
- replaces only an existing valid ZigGuard section on recompilation;
- migrates the legacy fully generated instruction-file format;
- rejects duplicate, reversed, or ambiguous managed markers without changing the file;
- validates that generated paths remain inside the selected repository root;
- writes through a temporary file before replacement.

The checker applies the same ownership model: human content outside a valid managed section does not count as policy drift.

### Adapters and capability levels

Each target binding must identify its capability:

- `instruction-context`;
- `enforceable`;
- `unsupported`.

A natural-language instruction is never described as a hard control.

### Check engine

Mechanical rules belong in the check engine rather than being expressed only as prose.

Examples:

- protected branch checks;
- secret scanner integration;
- required tests/lint execution;
- unsafe generated-file changes.

This is a post-MVP 0.1 milestone.

### Determinism

Compilation output depends only on:

- ZigGuard version and policy version;
- canonical policy input;
- adapter capabilities.

Targets and stack values are sorted before rendering so semantically equivalent ordering does not create noisy output.

### Cloud boundary

The local core remains useful offline. Hosted functionality should consume explicit metadata/artifacts and should not require source-code upload unless the user knowingly enables a feature that needs it.
