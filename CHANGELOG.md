# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Add dry-run flag to output mnanifests without applying them.

## [0.12.0] - 2025-07-22

### Changed

- Update Arch PKGBUILD for 0.11.0.
- Read shell KUBECONFIG env var if kubeconfig flag not provided.

## [0.11.0] - 2024-08-07

### Changed

- Added pod security context to ensure the pod runs as a non-root user.

## [0.10.0] - 2024-03-21

### Removed

- Remove deprecated PodSecurityPolicy.

### Changed

- Update Arch PKGBUILD for 0.9.0.

## [0.9.0] - 2024-02-15

### Added

- Add flags to configure group ID, running as non-root and allowing privilege escalation.

### Changed

- Update release workflows.

## [0.8.1] - 2023-05-05

### Changed

- Update release workflows.

## [0.8.0] - 2023-05-04

### Changed

- Support all client-go auth plugins.

## [0.7.1] - 2022-02-05

### Fixed

- Fix client-go to work with OIDC enabled clusters.
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

[Unreleased]: https://github.com/a7d-corp/sonar/compare/v0.12.0...HEAD
[0.12.0]: https://github.com/a7d-corp/sonar/compare/v0.11.0...v0.12.0
[0.11.0]: https://github.com/a7d-corp/sonar/compare/v0.10.0...v0.11.0
[0.10.0]: https://github.com/a7d-corp/sonar/compare/v0.9.0...v0.10.0
[0.9.0]: https://github.com/a7d-corp/sonar/compare/v0.8.1...v0.9.0
[0.8.1]: https://github.com/a7d-corp/sonar/compare/v0.8.0...v0.8.1
[0.8.0]: https://github.com/a7d-corp/sonar/compare/v0.7.1...v0.8.0
[0.7.1]: https://github.com/glitchcrab/sonar/compare/v0.7.0...v0.7.1
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
