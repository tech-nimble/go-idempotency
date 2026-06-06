# Changelog

All notable changes to this package are documented here. The format is based on
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and the project follows
[Semantic Versioning](https://semver.org/).

## [1.0.0]

First release under Nimble Tech.

### Added
- Gin idempotency middleware with Redis and no-op storages.
- Functional options (`WithStorage`, `WithHeaderKey`).
- Tests, golangci-lint config, GitHub Actions CI and dependabot.

### Fixed
- `Initialize` now wires the storage through functional options (previously did
  not compile).
- Removed the duplicate `headerIdempotencyKey` declaration.
- `NewIdempotency` returns an independent instance instead of mutating a shared
  global.
