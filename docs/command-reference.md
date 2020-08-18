# Command reference

This is an overview for the verless CLI and its commands.

## Installation

To install the verless CLI tool, check out the
[installation instructions](https://github.com/verless/verless/tree/docs#-installation).

## Commands

### verless

The top-level verless command does not provide any functionality.

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

| Option        | Short | Type   | Example                    | Description                                                      |
|---------------|-------|--------|----------------------------|------------------------------------------------------------------|
| `--output`    | `-o`  | String | `--output="/var/www/html"` | An alternative output directory where the website is written to. |
| `--overwrite` | -     | Bool   | `--overwrite`              | Allow verless to overwrite the output directory.                 |

### verless create

The `verless create` command does not provide any functionality.

### verless create project

`verless create project NAME` initializes a new verless standard project with all directories and files required for
running a build. If the `NAME` directory already exists, the command will fail. Use `--overwrite` to overwrite the
directory with the new project.

| Option        | Short | Type   | Example                    | Description                                                      |
|---------------|-------|--------|----------------------------|------------------------------------------------------------------|
| `--overwrite` | -     | Bool   | `--overwrite`              | Overwrite the specified directory if it already exists.          |

...
