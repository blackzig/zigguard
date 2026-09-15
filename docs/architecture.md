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
| Agent adapters|                      | Local checks   |
+---+---+---+---+                      +----------------+
    |   |   |   |
 Claude Codex Cursor Copilot
```

## Components

### Policy parser

Responsibilities:

- load `zigguard.yml`;
- validate against a versioned schema;
- normalize values into an internal policy model;
- reject unknown/unsafe combinations when required by the policy version.

### Policy compiler

Responsibilities:

- pass normalized rules to target adapters;
- maintain deterministic ordering;
- produce a plan/diff before writing;
- reject unsupported security-sensitive semantics.

### Adapters

Each adapter translates neutral policy semantics into the native mechanism of one tool.

Adapters must publish a capability matrix indicating whether each rule is:

- enforceable;
- representable as instruction only;
- unsupported.

### Check engine

Mechanical rules belong in the check engine rather than being expressed only as prose.

Examples:

- protected branch checks;
- secret scanner integration;
- required tests/lint execution;
- unsafe generated-file changes.

### Rule catalog

Built-in rules should have stable IDs, documented semantics, severity/action behavior, and test fixtures.

## Determinism

Compilation output should depend only on:

- ZigGuard version;
- policy version;
- canonical policy input;
- adapter version/capabilities.

Environment-specific data should not silently change generated files.

## File ownership

Generated content must distinguish between:

- fully managed files;
- managed sections inside shared files;
- advisory output that ZigGuard does not write automatically.

Overwriting user-maintained content without an explicit merge strategy is prohibited.

## Cloud boundary

The local core must remain useful offline. Hosted functionality should consume explicit metadata/artifacts and should not require source-code upload unless the user knowingly enables a feature that needs it.
