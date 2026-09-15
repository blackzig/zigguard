# GitHub Copilot Adapter

## MVP surface

GitHub Copilot consumes the generated shared policy from `AGENTS.md` for agent instruction scenarios that support it.

Copilot also supports repository-wide `.github/copilot-instructions.md` and path-specific `.github/instructions/*.instructions.md`. ZigGuard will use those mechanisms later when the neutral policy includes capabilities that benefit from those scopes.

## Current capability

Instruction context only.

Support varies across Copilot surfaces, so future capability matrices must identify the exact Copilot feature a generated artifact targets.
