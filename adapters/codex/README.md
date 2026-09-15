# OpenAI Codex Adapter

## MVP surface

Codex consumes the generated shared policy from `AGENTS.md`.

This keeps the ZigGuard output compatible with the repository-level instruction mechanism used by Codex while avoiding unnecessary target-specific duplication.

## Current capability

Instruction context only.

Future Codex-specific capabilities can be added only when they provide semantics beyond the shared `AGENTS.md` policy.
