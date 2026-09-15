<!-- zigguard:managed:start -->
# Shared ZigGuard policy for codex, copilot, cursor

> These are AI-agent instructions generated from ZigGuard policy v0.1. They do not replace branch protection, CI, hooks, scanners, or other technical controls.

## Project

- Name: `legacy-integrity-service`
- Stack: `java`, `oracle`, `spring-boot`

## Governance rules

- **BLOCK:** Do not refactor code outside the user's requested scope.
- **REQUIRE APPROVAL:** Obtain explicit human approval before changing a public API or external contract.
- **BLOCK:** Never expose, commit, log, or generate real secrets, tokens, credentials, or private keys.
- **BLOCK:** Do not build SQL by concatenating untrusted or runtime values; use parameter binding or an equivalent safe mechanism.
- **BLOCK:** Never force-push to a protected branch.
- **REQUIRED:** Run the relevant automated tests for changed behavior before declaring the task complete.
- **REQUIRED:** Run the repository's relevant lint/static checks for changed files before declaring the task complete.
<!-- zigguard:managed:end -->
