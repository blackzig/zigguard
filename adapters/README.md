# Adapters

Adapters translate the neutral ZigGuard policy model into the configuration understood by a specific AI coding tool.

Initial adapter targets:

- [Claude Code](claude/README.md)
- [OpenAI Codex](codex/README.md)
- [Cursor](cursor/README.md)
- [GitHub Copilot](copilot/README.md)

Every adapter must eventually publish a capability matrix for each ZigGuard rule:

- enforceable;
- instruction-only;
- unsupported.

Unsupported `block` and `require_approval` semantics must not be silently downgraded.
