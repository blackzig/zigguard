# Target capability matrix

Last reviewed: 2026-09-15.

ZigGuard distinguishes **instruction context** from **technical enforcement**. A target supporting a rule file does not automatically mean the platform can guarantee that an agent cannot violate the instruction.

| Target | MVP surface | ZigGuard capability | Notes |
| --- | --- | --- | --- |
| Claude Code | `CLAUDE.md` | instruction-context | Claude Code documents CLAUDE.md as persistent project context. Its documentation explicitly distinguishes this from enforcement and points to `PreToolUse` hooks when an action must be blocked. |
| OpenAI Codex | `AGENTS.md` | instruction-context | Codex uses AGENTS.md for persistent repository guidance. MVP 0.1 does not claim hard enforcement through this surface. |
| Cursor | `AGENTS.md` | instruction-context | Cursor supports AGENTS.md as project instructions and also supports `.cursor/rules/*.mdc`. Cursor documentation warns that AI guidance should not be the only security control. |
| GitHub Copilot | `AGENTS.md` | instruction-context | AGENTS.md support varies by Copilot feature. Copilot also supports repository-wide and path-specific instruction files. |

## Why the MVP shares AGENTS.md

Codex, Cursor, and several GitHub Copilot experiences understand `AGENTS.md`. Maintaining three generated copies of identical policy text would increase drift and may cause tools that recognize several instruction formats to receive redundant guidance.

For that reason, MVP 0.1 emits one shared `AGENTS.md`.

Claude Code currently documents `CLAUDE.md` rather than reading `AGENTS.md` directly. When both target families are enabled, ZigGuard generates a small `CLAUDE.md` that imports `@AGENTS.md`.

## Current enforcement status

MVP 0.1 has no rule marked `enforceable`. All generated target bindings are `instruction-context`.

Phase 2 will add mechanical checks and target-specific mechanisms only where they can actually enforce or gate behavior.

## Official references

- Claude Code memory/instructions: https://code.claude.com/docs/en/memory
- Claude Code hooks: https://code.claude.com/docs/en/hooks
- Cursor rules: https://prod.cursor.com/docs/rules
- GitHub Copilot custom instruction support: https://docs.github.com/en/copilot/reference/custom-instructions-support
- OpenAI Codex / AGENTS.md guidance: https://openai.com/business/guides-and-resources/how-openai-uses-codex/
