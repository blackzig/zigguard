# ZigGuard Policy Specification — Draft v0.1

This document defines the initial shape of `zigguard.yml`.

The format is intentionally small. The goal of v0.1 is to validate universal policy semantics before expanding the rule catalog.

## Top-level fields

### `version`

Required string. For this draft:

```yaml
version: "0.1"
```

### `project`

Required project metadata.

```yaml
project:
  name: billing-api
  stack:
    - java
    - spring-boot
```

### `targets`

Required list of adapter targets.

Draft values:

- `claude`
- `codex`
- `cursor`
- `copilot`

### `rules`

Required policy rules.

## Actions

Rules that use an action share this vocabulary:

| Action | Meaning |
| --- | --- |
| `allow` | Explicitly permitted |
| `warn` | Allowed, but the user should be informed |
| `require_approval` | Must not proceed without explicit human approval |
| `block` | Must not proceed |

Adapters may have different technical capabilities. They must never silently compile `block` or `require_approval` into weaker semantics.

## Draft rule catalog

### Scope

```yaml
rules:
  scope:
    refactor_outside_request: block
    public_api_changes: require_approval
```

### Security

```yaml
rules:
  security:
    secrets: block
    dynamic_sql: block
```

### Quality

```yaml
rules:
  quality:
    tests_required: true
    lint_required: true
```

### Git

```yaml
rules:
  git:
    force_push_protected_branches: block
```

## Compatibility

Within policy version `0.1`:

- removing a field or changing its meaning is considered breaking;
- adding an optional rule may be backward compatible;
- unknown fields should fail validation by default during this draft phase.

## Compilation safety

If a target cannot safely represent a requested rule, compilation should return a clear unsupported-capability error rather than generating misleading output.

## Machine-readable schema

See [../schemas/zigguard.schema.json](../schemas/zigguard.schema.json).
