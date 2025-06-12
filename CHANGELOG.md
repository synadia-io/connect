# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Comprehensive test coverage for CLI, builders, and client packages
- Test coverage increased from 0.9% to 13.2%
- Mock client implementation for testing without NATS connection
- Helper functions for CLI testing to avoid os.Exit calls
- Fluent Url() method to NatsConfigBuilder
- Detailed API documentation in `docs/api/`
- CLI reference documentation in `docs/cli-reference.md`
- Component library reference in `docs/reference/README.md`
- Expanded README with better structure and examples
- Comprehensive CONTRIBUTING guide with development workflow

### Changed
- Updated Go dependencies to latest versions
- Improved error handling in splitCommand to prevent panic on empty input
- Enhanced README with badges, table of contents, and clear sections
- Expanded CONTRIBUTING.md with detailed guidelines

### Fixed
- Fixed panic in splitCommand when given empty string
- Fixed test compilation issues with proper type assertions

### Security
- Updated outdated dependencies with security patches

## [Previous Releases]

For previous releases, see the [releases page](https://github.com/synadia-io/connect/releases).