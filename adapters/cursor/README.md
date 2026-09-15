# Cursor Adapter

## MVP surface

Cursor consumes the generated shared policy from `AGENTS.md`, which Cursor supports as a project instruction mechanism.

Cursor also supports project rules in `.cursor/rules/*.mdc`. ZigGuard will use those files later for conditional/path-specific behavior where they add value beyond the shared policy.

## Current capability

Instruction context only.

The adapter must not claim that a Cursor rule or AGENTS instruction is a hard security control.
