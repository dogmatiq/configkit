# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog], and this project adheres to
[Semantic Versioning].

<!-- references -->
[Keep a Changelog]: https://keepachangelog.com/en/1.0.0/
[Semantic Versioning]: https://semver.org/spec/v2.0.0.html

## [Unreleased]

### Added

- Add `message.IsIntersectingN()` and `IsSubsetN()`
- Add `message.IntersectionN()` and `UnionN()` and `DiffN()`

## [0.2.1] - 2020-01-16

## Changed

- `NameOf()` and `TypeOf()` now produce a more meaningful panic message when passed a `nil` message

## [0.2.0] - 2019-12-16

### Removed

- **[BC]** Remove `Application.ForeignMessageNames()`
- **[BC]** Remove `RichApplication.ForeignMessageTypes()`
- **[BC]** Remove `Errorf` and `Panicf()`
- **[BC]** Remove `Roles` field from `EntityMessageNames` and `EntityMessageTypes`

### Added

- Add `ForeignMessageNames()` and `ForeignMessageTypes()`
- Add `IsApplicationEqual()` and `IsHandlerEqual()`
- Add `ToString()`
- Add `IsEqual()` method to `HandlerSet` and `RichHandlerSet`
- Add `Identity.[Un]MarshalText()` and `[Un]MarshalBinary()`
- Add `IsEqual()` method to `NameRoles`, `NameSet`, `TypeRoles` and `TypeSet`
- Add `IsEqual()`, `All()` and `RoleOf()` methods to `EntityMessageNames` and `EntityMessageTypes`

## [0.1.1]

### Added

- Add `HandlerSet.AcceptVisitor()` and `RichHandlerSet.AcceptRichVisitor()`

## [0.1.0] - 2019-12-02

- Initial release

<!-- references -->
[Unreleased]: https://github.com/dogmatiq/configkit
[0.1.0]: https://github.com/dogmatiq/configkit/releases/v0.1.0
[0.1.1]: https://github.com/dogmatiq/configkit/releases/v0.1.1
[0.2.0]: https://github.com/dogmatiq/configkit/releases/v0.2.0
[0.2.1]: https://github.com/dogmatiq/configkit/releases/v0.2.1

<!-- version template
## [0.0.1] - YYYY-MM-DD

### Added
### Changed
### Deprecated
### Removed
### Fixed
### Security
-->
