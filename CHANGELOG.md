# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Fixed

- Mount host filesystem into container when `--node-exec` flag is used.

## [0.7.0] - 2021-08-29

### Added

- Add --node-exec flag to grant access to a host's IPC/net/PID namespaces.

## [0.6.0] - 2021-08-28

### Added

- Allow setting the cluster context via the --context flag.

### Changed

- Align the kube config flag with kubectl.

## [0.5.0] - 2021-08-15

### Added

- Github workflow to build and test on PRs.
- Flag to provide a nodename to schedule the pod on.

## [0.4.0] - 2021-08-05

### Added

- Documentation for all flags and provided usage examples.

## [0.3.1] - 2021-08-05

### Removed

- Remove print leftover from debugging.

## [0.3.0] - 2021-08-05

### Added

- Expose the ability to set the pod's userID.
- Expose the ability to set the pod's command and args.
- Added short versions of some flags.

## [0.2.2] - 2021-08-05

### Added

- Added PKGBUILD for Arch Linux.

## [0.2.1] - 2021-08-05

### Added

- Added Makefile.

## [0.2.0] - 2021-08-04

### Added

- Added version command.

### Changed

- Aligned logging style.

## [0.1.0] - 2021-08-04

### Added

- Initial release; rough around the edges but usable.

[Unreleased]: https://github.com/glitchcrab/sonar/compare/v0.7.0...HEAD
[0.7.0]: https://github.com/glitchcrab/sonar/compare/v0.6.0...v0.7.0
[0.6.0]: https://github.com/glitchcrab/sonar/compare/v0.5.0...v0.6.0
[0.5.0]: https://github.com/glitchcrab/sonar/compare/v0.4.0...v0.5.0
[0.4.0]: https://github.com/glitchcrab/sonar/compare/v0.3.1...v0.4.0
[0.3.1]: https://github.com/glitchcrab/sonar/compare/v0.3.0...v0.3.1
[0.3.0]: https://github.com/glitchcrab/sonar/compare/v0.2.2...v0.3.0
[0.2.2]: https://github.com/glitchcrab/sonar/compare/v0.2.1...v0.2.2
[0.2.1]: https://github.com/glitchcrab/sonar/compare/v0.2.0...v0.2.1
[0.2.0]: https://github.com/glitchcrab/sonar/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/glitchcrab/sonar/releases/tag/v0.1.0
