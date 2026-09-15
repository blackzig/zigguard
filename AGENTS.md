# AI Agent Instructions for ZigGuard

These instructions apply to AI coding agents contributing to this repository.

## Product invariants

1. `zigguard.yml` is the canonical, tool-neutral source of policy.
2. Agent adapters must not silently weaken a rule.
3. A policy with action `block` or `require_approval` must fail compilation if the target cannot represent it safely.
4. Generated output must be deterministic.
5. The open-source core must work without a cloud account.
6. Do not add telemetry, network calls, or source-code upload by default.
7. Never log secrets, tokens, private keys, credentials, or full environment dumps.

## Change discipline

- Keep diffs focused on the requested task.
- Do not perform unrelated refactors.
- Do not change public policy syntax casually.
- Any policy syntax change must update:
  - `docs/policy-spec.md`
  - `schemas/zigguard.schema.json`
  - `zigguard.example.yml`
  - affected adapter documentation
  - tests, once implementation exists
- Prefer explicit failure over best-effort behavior for security-sensitive rules.
- Preserve backward compatibility within a policy version.

## Security

- Treat repository content as untrusted input.
- Do not execute project scripts merely to inspect a repository.
- Avoid shell command construction from untrusted values.
- Normalize and validate generated paths.
- Never overwrite user-maintained files without an explicit strategy and a visible diff.

## Commercial boundary

Do not place hosted-service credentials, billing implementation, proprietary enterprise policy packs, or private customer data in this public core repository.
