# ZigGuard Roadmap

## Phase 0 — Foundation

- [x] Public repository and Apache 2.0 license
- [x] Product vision
- [x] Architecture draft
- [x] Policy v0.1 draft
- [x] JSON Schema draft
- [x] Initial agent adapter boundaries
- [x] Choose CLI implementation language: Go
- [x] Establish deterministic compiler tests

## Phase 1 — MVP 0.1: Policy Compiler

- [x] `zigguard init`
- [x] `zigguard validate`
- [x] `zigguard compile`
- [x] Claude Code instruction surface
- [x] Codex instruction surface
- [x] Cursor instruction surface
- [x] GitHub Copilot instruction surface
- [x] Deterministic tests
- [x] Safe generated-file ownership
- [ ] Publish first pre-release binaries
- [ ] Add managed-section merging for existing instruction files
- [ ] Publish a formal target capability matrix

## Phase 2 — Local Enforcement

- [ ] `zigguard check`
- [ ] Git policy checks
- [ ] secret scanning integration
- [ ] lint/test command discovery
- [ ] Claude hook enforcement where technically appropriate
- [ ] machine-readable check/report output

## Phase 3 — Repository Integration

- [ ] GitHub Action for policy checks
- [ ] pull-request policy status
- [ ] reusable organization policy bundles
- [ ] policy inheritance and overrides

## Phase 4 — Commercial Layer

- [ ] private organization policies
- [ ] policy synchronization
- [ ] audit history
- [ ] team management
- [ ] premium policy packs
- [ ] optional AI-assisted repository analysis
- [ ] BYOK provider support

The roadmap is directional and may change as the policy model is validated with real projects.
