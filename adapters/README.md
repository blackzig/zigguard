# Adapters

Adapters translate the neutral ZigGuard policy model into instruction or enforcement surfaces supported by AI coding tools.

## MVP 0.1 surfaces

| Target | Surface | Current capability |
| --- | --- | --- |
| Claude Code | `CLAUDE.md` | instruction context |
| OpenAI Codex | `AGENTS.md` | instruction context |
| Cursor | `AGENTS.md` | instruction context |
| GitHub Copilot | `AGENTS.md` | instruction context |

The MVP deliberately reuses `AGENTS.md` where multiple tools support it. ZigGuard should not create four divergent copies of the same policy merely because four tools are selected.

Target-specific surfaces such as Claude hooks, Cursor project rules, and Copilot path-specific instructions will be added when they provide semantics that the shared surface cannot express.

Every adapter must distinguish:

- **instruction context** — guidance placed into an agent's context;
- **enforceable** — a technical mechanism can deny or gate the action;
- **unsupported** — the target cannot safely represent the requested behavior.

ZigGuard must never describe instruction-only behavior as hard enforcement.
