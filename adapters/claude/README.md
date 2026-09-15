# Claude Code Adapter

## MVP surface

Claude Code receives project instructions through `CLAUDE.md`.

When the policy also targets Codex, Cursor, or GitHub Copilot, ZigGuard keeps the full shared policy in `AGENTS.md` and generates a small `CLAUDE.md` containing an `@AGENTS.md` import. Claude Code supports importing additional instruction files with the `@path` syntax.

When Claude is the only target, the complete generated policy is written directly to `CLAUDE.md`.

## Current capability

Instruction context only.

Future versions can add Claude-specific hooks for rules that have a safe technical enforcement path. A hook-backed rule must be reported separately from a natural-language instruction.
