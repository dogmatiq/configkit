# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog], and this project adheres to
[Semantic Versioning].

<!-- references -->
[Keep a Changelog]: https://keepachangelog.com/en/1.0.0/
[Semantic Versioning]: https://semver.org/spec/v2.0.0.html

## [Unreleased]

### Removed

- **[BC]** Remove `Application.ForeignMessageNames()`
- **[BC]** Remove `RichApplication.ForeignMessageTypes()`

### Added

- Add `ForeignMessageNames()`
- Add `ForeignMessageTypes()`
- Add `IsApplicationEqual()` and `IsHandlerEqual()`
- Add `HandlerSet.IsEqual()`
- Add `RichHandlerSet.IsEqual()`
- Add `Identity.[Un]MarshalText()` and `[Un]MarshalBinary()`
- Add `NameRoles.IsEqual()` and `TypeRoles.IsEqual()`
- Add `NameSet.IsEqual()` and `TypeSet.IsEqual()`
- Add `EntityMessageNames.IsEqual()` and `EntityMessageTypes.IsEqual()`

## [0.1.1]

### Added

- Add `HandlerSet.AcceptVisitor()` and `RichHandlerSet.AcceptRichVisitor()`

## [0.1.0] - 2019-12-02

- Initial release

<!-- references -->
[Unreleased]: https://github.com/dogmatiq/configkit
[0.1.0]: https://github.com/dogmatiq/configkit/releases/v0.1.0
[0.1.1]: https://github.com/dogmatiq/configkit/releases/v0.1.1

<!-- version template
## [0.0.1] - YYYY-MM-DD

### Added
### Changed
### Deprecated
### Removed
### Fixed
### Security
-->
