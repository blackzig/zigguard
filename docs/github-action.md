# ZigGuard GitHub Action

ZigGuard ships a composite GitHub Action that runs `zigguard check` against a repository in CI.

The action requires only read access to repository contents. It does not execute the consumer project's test scripts, package-manager hooks, or arbitrary project commands.

## Pull request example

After a release tag is available, a repository can add:

```yaml
name: ZigGuard

on:
  pull_request:
  push:
    branches:
      - main

permissions:
  contents: read

jobs:
  governance:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v7

      - name: Check AI coding governance
        uses: blackzig/zigguard@v0.1.0-alpha.1
```

If generated governance artifacts are missing, drifted, or replaced by unmanaged files, the job fails.

## Inputs

| Input | Default | Purpose |
| --- | --- | --- |
| `policy-file` | `zigguard.yml` | policy path relative to the repository workspace |
| `root` | `.` | repository root containing generated governance files |
| `format` | `text` | `text` or `json` check output |
| `go-version` | `1.23.x` | Go version used to build ZigGuard from the pinned action source |

Example with custom paths:

```yaml
- uses: blackzig/zigguard@v0.1.0-alpha.1
  with:
    policy-file: config/zigguard.yml
    root: .
    format: json
```

## Why the action builds from source

During the early alpha, the action builds the ZigGuard CLI from the exact Git ref that the workflow pins. This gives the action and CLI identical version semantics and avoids downloading an unverified moving binary.

A later stable release may use signed/prebuilt artifacts for faster startup while preserving verification.

## Security model

The action:

- uses the source pinned by the caller's `uses: blackzig/zigguard@...` reference;
- passes inputs through environment variables instead of injecting them into shell source;
- needs no secrets;
- requires no write permission;
- runs only ZigGuard's own read-only `check` command against the consumer workspace.

The action validates generated-governance integrity. It does not turn natural-language agent instructions into a hard sandbox.
