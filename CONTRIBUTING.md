# Contributing to ZigGuard

ZigGuard is currently in an early specification-first phase.

## Before contributing

Please read:

- [AGENTS.md](AGENTS.md)
- [docs/product-vision.md](docs/product-vision.md)
- [docs/architecture.md](docs/architecture.md)
- [docs/policy-spec.md](docs/policy-spec.md)

## Pull request expectations

A contribution should:

- solve one clearly scoped problem;
- avoid unrelated refactoring;
- explain behavioral changes;
- update documentation when policy behavior changes;
- include tests when implementation code is added;
- avoid adding mandatory cloud dependencies to the core.

## Policy changes

Changes to `zigguard.yml` are API changes. They require synchronized updates to the specification, schema, example configuration, adapters, and tests.

## Security issues

Do not report vulnerabilities in a public issue. Follow [SECURITY.md](SECURITY.md).
