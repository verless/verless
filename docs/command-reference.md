# Command reference

This is an overview for the verless CLI and its commands.

## Contents

* [Installation](#installation)
* [`verless`](#verless)
* [`verless build`](#verless-build)
* [`verless create`](#verless-create)
    * [`verless create project`](#verless-create-project)
* [`verless serve`](#verless-serve)
* [`verless version`](#verless-version)

## Installation

To install the verless CLI tool, check out the
[installation instructions](../README.md#img-srchttpsverlessdominikbraunioassetsimgdotpng-installation).

## verless

The top-level verless command does not provide any functionality.

## verless build

`verless build PATH` runs a build for the specified `PATH`. If your current directory is the project directory, you
can build your project using `verless build .`.

The `build` command will generate all pages and collect all errors that occurred during the build. Those errors will be
returned as a list of things that have to be fixed - the build itself will _not_ finish.

For security reasons, verless doesn't overwrite an existing output directory if it already exists. If you've run a
build before, the created `target` directory cannot be just overwritten by the next build. You explicitly have to
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

## verless create

The `verless create` command does not provide any functionality.

## verless create project

`verless create project NAME` initializes a new verless default project with all directories and files required for
running a build. If the `NAME` directory already exists, the command will fail. Use `--overwrite` to overwrite the
directory with the new project.

**Caution:** The entire directory will be deleted when doing so.

| Option        | Short | Type   | Example       | Description                                             |
|---------------|-------|--------|---------------|---------------------------------------------------------|
| `--overwrite` | -     | Bool   | `--overwrite` | Overwrite the specified directory if it already exists. |

## verless create theme

`verless create theme PROJECT NAME` initializes a new verless theme called `NAME` within an existing project called
`PROJECT`. For example, if there's a project `my-blog` and you want to create a theme called `dark-theme`, use the
following command:

```shell script
$ verless create theme my-blog dark-theme
```

Just as with other commands, `PROJECT` is the path to your project. If you already are _inside_ `my-blog` directory, the
command is:

```shell script
$ verless create theme . dark-theme
```

## verless serve

`verless serve PROJECT` starts a tiny webserver that serves your static site. By default, verless listens to port 8080
on all network interfaces, so your project is available under `localhost:8080` for example.

The `--watch` flag is useful for local development because verless re-builds your website when a file has changed, so
you're able to view your changes immediately.

Because `verless serve` re-builds your static site when the `--watch` flag is used, it additionally accepts all options
that [`verless build`](#verless-build) does.

| Option    | Short | Type   | Example          | Description                                                        |
|-----------|-------|--------|------------------|--------------------------------------------------------------------|
| `--port`  | `-p`  | UInt16 | `--port 8000`    | The TCP port for serving the static site.                          |
| `--watch` | `-w`  | Bool   | `--watch`        | Watch all project file and re-build the site if something changed. |
| `--ip`    | `-i`  | String | `--ip 127.0.0.1` | The network address for serving the static site.                   |

## verless version

`verless version` prints the installed verless version.

| Option    | Short | Type   | Example   | Description                          |
|-----------|-------|--------|-----------|--------------------------------------|
| `--quiet` | `-q`  | Bool   | `--quiet` | Only print the plain version number. |

<p align="center">
<br>
<a href="https://github.com/verless/verless"><img src="https://verless.dominikbraun.io/assets/img/icon-light.png"></a>
</p>
