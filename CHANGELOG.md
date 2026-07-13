# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.3.0] - 2026-07-14

### Added
- Markdown renderer for terminal display
- Code block background padding to the terminal width

### Fixed
- Prevent code block padding from wrapping wide characters
- Show the installed module version for tagged `go install` builds

## [0.2.0] - 2026-07-12

### Added
- Config command with interactive settings management
- Init command to create `.dircard` file
- Configurable candidate order for file finding

## [0.1.0] - 2026-03-16

### Added
- Install and uninstall commands for shell hooks
- Dircard shell hook command
- Show command to display directory notes
- Initial project setup

### Changed
- Pin Go version from go.mod in release workflow
- Update README and Japanese translation

### Infrastructure
- Add goreleaser config and release workflow

[Unreleased]: https://github.com/dircard/dircard/compare/v0.3.0...HEAD
[0.3.0]: https://github.com/dircard/dircard/releases/tag/v0.3.0
[0.2.0]: https://github.com/dircard/dircard/releases/tag/v0.2.0
[0.1.0]: https://github.com/dircard/dircard/releases/tag/v0.1.0
