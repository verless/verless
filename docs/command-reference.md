# Command reference

This is an overview for the verless CLI and its commands.

## Installation

To install the verless CLI tool, check out the
[installation instructions](https://github.com/verless/verless/tree/docs#-installation).

## Commands

### verless

The top-level verless command doesn't provide any functionality.

### verless build

`verless build PATH` runs a build for the specified `PATH`. If your current directory is the project directory, you
can build your project using `verless build .`.

The `build` command will generate all pages and collect all errors that occurred during the build. Those errors will be
returned as a list of things that have to be fixed - the build itself will _not_ finish.

For security reasons, verless doesn't overwrite an existing output directory if it already exists. If you've run a
build before, creating a `target` directory, this cannot be just overwritten by the next build. You explicitly have to
allow verless to overwrite it using `--overwrite`. If you're getting tired of this and know what you're doing, you may
allow this in the project configuration:

```yaml
build:
  overwrite: true
```

**Caution:** This will also overwrite any other output directory specified with `--output`.

...
