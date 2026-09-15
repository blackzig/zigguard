# Java + Spring Boot + Oracle example

This fixture demonstrates a conservative ZigGuard policy for a Java service using Spring Boot and Oracle.

The policy blocks out-of-scope refactors, secret exposure, unsafe dynamic SQL, and protected-branch force pushes. Public API changes require explicit approval, while tests and lint/static checks are required before task completion.

## Compile

From the repository root:

```bash
go run ./cmd/zigguard validate --file examples/java-spring-oracle/zigguard.yml
go run ./cmd/zigguard compile \
  --file examples/java-spring-oracle/zigguard.yml \
  --root /tmp/zigguard-example
```

The `generated/` directory contains the expected deterministic output used by the compiler contract test.

This example is intentionally generic and contains no proprietary application rules.
