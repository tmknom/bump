# bump

`bump` bumps up version automatically.
It supports only for semantic versioning, such as X.Y.Z.

## Getting Started

First, the following command generates `VERSION` file.

```shell
bump init
```

Then you can bump up to version.
For example, you can bump up to minor version:

```shell
bump minor
```

Of course, major version or patch version bumps up are also possible:

```shell
bump major
bump patch
```

## Usage

```shell
Usage:
  bump <command> [<version>] [flags]

Commands:
  init              Init version
  major             Bump up major version
  minor             Bump up minor version
  patch             Bump up patch version
  show              Show the current version

Flags:
  --help            Show help for command
  --version         Show bump version

Examples:
  $ bump init
  $ bump patch
  $ bump minor 1.0.0
```

## Changelog

See [CHANGELOG.md](/CHANGELOG.md).

## Author

[@tmknom](https://github.com/tmknom/)

## License

Apache 2 Licensed. See LICENSE for full details.
