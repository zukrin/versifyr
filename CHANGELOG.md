# Changelog

All notable changes to this project will be documented in this file.

## [0.1.0] - 2026-03-31

### Added
- Comprehensive unit testing suite for `logging`, `configuration`, and `command` packages.
- Integration test `TestCLIWorkflow` for full CLI lifecycle validation.
- `go-test-coverage` integration for monitoring coverage thresholds.
- `golangci-lint` workflow and local `task lint` for code quality.
- Split GitHub workflows into `ci.yml` (tests/coverage), `lint.yml` (quality), and `release.yml` (build/deploy).
- Enhanced `Taskfile.yml` with tasks for testing, coverage, linting, and automated versioning.
- Coverage badge in README.md.

### Changed
- Major refactoring: extracted template application logic to `ConfigFile.ApplyTemplates` for improved testability.
- Updated `advance-version` task to fetch current version from latest git tag instead of source.
- Pinned all GitHub Actions to commit SHAs for enhanced security.
- Improved logging safety across the project using explicit format strings.

### Fixed
- `init` command now correctly returns `os.ErrExist` when the target folder already exists.
- Fixed version calculation bug in `Taskfile.yml`.
- Resolved various linting issues (unchecked errors, syntax errors in tests).
- Fixed `go-test-coverage` action SHA and `golangci-lint-action` SHA.

## [0.0.16] - 2026-03-22
- Pre-refactoring version with initial features.
