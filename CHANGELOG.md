# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog], and this project adheres to
[Semantic Versioning].

<!-- references -->

[keep a changelog]: https://keepachangelog.com/en/1.0.0/
[semantic versioning]: https://semver.org/spec/v2.0.0.html

## [0.13.2] - 2024-08-12

### Fixed

- Fixed panicking in static analysis functionality.

## [0.13.1] - 2024-07-16

### Changed

- Added support for `Disable()` method in `dogmatiq/dogma` v0.13.1.

### Fixed

- Fixed incorrect panic value when calling
  `ProjectionConfigurer.DeliveryPolicy()` with a `nil` policy.

## [0.13.0] - 2024-03-26

- **[BC]** Updated to Dogma v0.13.0

## [0.12.2] - 2024-01-19

### Fixed

- Ensure static analysis correctly identifies handlers with non-pointer
  receivers when registered via a pointer.
- Ensure static analysis ignores routes that are not actually passed to a
  configurer via the `Routes()`.

## [0.12.1] - 2023-05-04

### Changed

- Modify `static` package to analyze the latest changes introduced by Dogma
  v0.12.0.

## [0.12.0] - 2023-04-09

This release updates the `configkit` implementation to adhere to Dogma v0.12.0
interfaces. It does not yet include support for static analysis of the new
`Routes()` API.

### Added

- Added `RichProjection.DeliveryPolicy()`

### Changed

- **[BC]** Application and handler identity keys must be RFC 4122 UUIDs

## [0.11.1] - 2021-09-17

### Added

- Add `static` package, which produces application configurations using static analysis

## [0.11.0] - 2021-02-24

### Changed

- **[BC]** Updated to Dogma v0.11.0

## [0.10.0] - 2021-01-20

### Added

- Add `api.NewServer()`
- Add `api.Server`

### Changed

- **[BC]** The `api` package now implements the `configspec` from `dogmatiq/interopspec`
- **[BC]** `api.Client` is now a struct, not an interface

### Removed

- **[BC]** Remove `api.Client.ListApplicationIdentities()`
- **[BC]** Remove `api.RegisterServer()`
- **[BC]** Remove the `api/discovery` package, use `dogmatiq/discoverkit` instead
- **[BC]** Remove the `api/fixtures` package

## [0.9.1] - 2020-11-20

### Added

- Add `ValidateIdentityName()` and `ValidateIdentityKey()`

## [0.9.0] - 2020-11-06

### Changed

- **[BC]** Updated to Dogma v0.9.0

## [0.8.0] - 2020-11-03

### Changed

- **[BC]** Updated to Dogma v0.8.0

## [0.7.4] - 2020-11-01

### Added

- Add `visualization/dot` package (migrated from `dogmatiq/graphkit`)

## [0.7.3] - 2020-04-30

### Added

- Add `HandlerSet.MessageNames()`
- Add `RichHandlerSet.MessageTypes()`

## [0.7.2] - 2020-03-25

### Added

- Add `Len()` to `NameCollection` and `TypeCollection` interfaces

## [0.7.1] - 2020-03-16

### Added

- Add `message.IsEqualSetN()` and `IsEqualSetT()`

## [0.7.0] - 2020-03-15

### Added

- Add `discovery.Connector.IsFatal`
- Add `discovery.Inspector.IsFatal`
- Add `NameRoles.RangeByRole()` and `FilterByRole()`
- Add `TypeRoles.RangeByRole()` and `FilterByRole()`

### Removed

- **[BC]** Remove `discovery.Connector.Logger`
- **[BC]** Remove `discovery.Inspector.Logger`

## [0.6.0] - 2020-03-14

### Changed

- **[BC]** Change `discovery.Inspector.Ignore()` to accept a `*discovery.Application`

## [0.5.0] - 2020-03-14

### Added

- Add `MessageEntityNames.Foreign()`
- Add `MessageEntityTypes.Foreign()`

### Removed

- **[BC]** Remove `ForeignMessageNames()` and `ForeignMessageTypes()`

## [0.4.2] - 2020-03-13

### Added

- Add `static.Discoverer` for "discovering" static lists of targets

## [0.4.1] - 2020-03-10

### Added

- Add `message.TypeFromReflect()`

## [0.4.0] - 2020-03-09

### Added

- Add `discovery.TargetExecutor` and `ClientExecutor`
- Add `discovery.Connector.Ignore`
- Add `discovery.ApplicationObserver`, `ApplicationObserverSet` and `ApplicationExecutor`
- Add `discovery.Inspector`

### Changed

- **[BC]** Change the internal gRPC API namespace from `dogma.configkit.v1` to `dogma.config.v1`
- **[BC]** Change `api.Client` to an interface
- **[BC]** `discovery.Connector` no longer implements `TargetObserver`, instead use a `TargetExecutor` to call `Connector.Run()`

## [0.3.1] - 2020-03-08

### Added

- Add the `api` package, a gRPC API for communicating APP configurations over the network

## [0.3.0] - 2020-01-29

### Changed

- **[BC]** Rename `message.NameCollection.Each()` to `Range()` for consistency with the Go standard library
- **[BC]** Rename `message.TypeCollection.Each()` to `Range()` for consistency with the Go standard library

### Added

- Add `HandlerSet.RangeAggregates()`, `RangeProcesses()`, `RangeIntegrations()` and `RangeProjections()`
- Add `HandlerSet.Aggregates()`, `Processes()`, `Integrations()` and `Projections()`
- Add `RichHandlerSet.Aggregates()`, `Processes()`, `Integrations()` and `Projections()`
- Add `RichHandlerSet.RangeAggregates()`, `RangeProcesses()`, `RangeIntegrations()` and `RangeProjections()`

## [0.2.2] - 2020-01-29

### Added

- Add `message.IsIntersectingN()`, `IsSubsetN()`, `IntersectionN()`, `UnionN()` and `DiffN()`
- Add `message.IsIntersectingT()`, `IsSubsetT()`, `IntersectionT()`, `UnionT()` and `DiffT()`

## [0.2.1] - 2020-01-16

### Changed

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

[unreleased]: https://github.com/dogmatiq/configkit
[0.1.0]: https://github.com/dogmatiq/configkit/releases/v0.1.0
[0.1.1]: https://github.com/dogmatiq/configkit/releases/v0.1.1
[0.2.0]: https://github.com/dogmatiq/configkit/releases/v0.2.0
[0.2.1]: https://github.com/dogmatiq/configkit/releases/v0.2.1
[0.2.2]: https://github.com/dogmatiq/configkit/releases/v0.2.2
[0.3.0]: https://github.com/dogmatiq/configkit/releases/v0.3.0
[0.3.1]: https://github.com/dogmatiq/configkit/releases/v0.3.1
[0.4.0]: https://github.com/dogmatiq/configkit/releases/v0.4.0
[0.4.1]: https://github.com/dogmatiq/configkit/releases/v0.4.1
[0.4.2]: https://github.com/dogmatiq/configkit/releases/v0.4.2
[0.5.0]: https://github.com/dogmatiq/configkit/releases/v0.5.0
[0.6.0]: https://github.com/dogmatiq/configkit/releases/v0.6.0
[0.7.0]: https://github.com/dogmatiq/configkit/releases/v0.7.0
[0.7.1]: https://github.com/dogmatiq/configkit/releases/v0.7.1
[0.7.2]: https://github.com/dogmatiq/configkit/releases/v0.7.2
[0.7.3]: https://github.com/dogmatiq/configkit/releases/v0.7.3
[0.7.4]: https://github.com/dogmatiq/configkit/releases/v0.7.4
[0.8.0]: https://github.com/dogmatiq/configkit/releases/v0.8.0
[0.9.0]: https://github.com/dogmatiq/configkit/releases/v0.9.0
[0.9.1]: https://github.com/dogmatiq/configkit/releases/v0.9.1
[0.10.0]: https://github.com/dogmatiq/configkit/releases/v0.10.0
[0.11.0]: https://github.com/dogmatiq/configkit/releases/v0.11.0
[0.11.1]: https://github.com/dogmatiq/configkit/releases/v0.11.1
[0.12.0]: https://github.com/dogmatiq/configkit/releases/v0.12.0
[0.12.1]: https://github.com/dogmatiq/configkit/releases/v0.12.1
[0.12.2]: https://github.com/dogmatiq/configkit/releases/v0.12.2
[0.13.0]: https://github.com/dogmatiq/configkit/releases/v0.13.0
[0.13.1]: https://github.com/dogmatiq/configkit/releases/v0.13.1

<!-- version template
## [0.0.1] - YYYY-MM-DD

### Added
### Changed
### Deprecated
### Removed
### Fixed
### Security
-->
