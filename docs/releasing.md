# Releasing ZigGuard

ZigGuard releases are built from Git tags by `.github/workflows/release.yml`.

## Versioning

The CLI product version and the policy format version are separate concepts.

- CLI example: `0.1.0-alpha.1`
- policy format: `0.1`

The release workflow injects the tag version into the Go CLI at build time.

## First planned prerelease

```text
v0.1.0-alpha.1
```

Do not create the tag until the release candidate commit is on `main` and CI is green.

## Release process

```bash
git checkout main
git pull
git tag v0.1.0-alpha.1
git push origin v0.1.0-alpha.1
```

The workflow then:

1. runs tests and `go vet`;
2. cross-compiles Linux, macOS, and Windows for amd64 and arm64;
3. packages the binary with README and LICENSE;
4. creates a GitHub Release;
5. marks tags containing a hyphen as prereleases.

## Artifacts

Expected archive names include:

```text
zigguard_0.1.0-alpha.1_linux_amd64.tar.gz
zigguard_0.1.0-alpha.1_linux_arm64.tar.gz
zigguard_0.1.0-alpha.1_darwin_amd64.tar.gz
zigguard_0.1.0-alpha.1_darwin_arm64.tar.gz
zigguard_0.1.0-alpha.1_windows_amd64.zip
zigguard_0.1.0-alpha.1_windows_arm64.zip
```

Publishing a tag is a public release action. The workflow is prepared in advance so a release can be reviewed before the tag is created.
