<p align="center">
<br>
<br>
<br>
<img src="https://verless.dominikbraun.io/static/img/logo-github-v1.0.0.png">
<br>
<br>
<br>
</p>

<h3 align="center">A simple and lightweight Static Site Generator.</h3>

<p align="center">
<a href="https://circleci.com/gh/verless/verless"><img src="https://circleci.com/gh/verless/verless.svg?style=shield"></a>
<a href="https://goreportcard.com/report/github.com/verless/verless"><img src="https://goreportcard.com/badge/github.com/verless/verless"></a>
<a href="https://www.codefactor.io/repository/github/verless/verless"><img src="https://www.codefactor.io/repository/github/verless/verless/badge" /></a>
<a href="https://github.com/verless/verless/releases"><img src="https://img.shields.io/github/v/release/verless/verless?sort=semver"></a>
<a href="LICENSE"><img src="https://img.shields.io/badge/license-Apache--2.0-brightgreen"></a>
<br>
<br>
<br>
</p>

---

**verless** (pronounced like _serverless_) is a Static Site Generator designed for Markdown-based content with focus on
simplicity and performance. It reads your Markdown files, applies your HTML templates and renders them as a website.

## <img src="https://verless.dominikbraun.io/static/img/list-icon-v1.0.0.png"> Features

* **Flexible theming:** Themes define the HTML templates as well as styles and scripts for your frontend. Initializing
a new theme can be done within a single command, and setting it as active requires a single line of configuration. Since
verless provides pre-defined template variables, themes are interchangeable.

* **Rapid development:** Just like themes, new projects can be created using a single command. Local development
requires no third-party webserver as verless comes with a built-in webserver and rebuilds your project when something
changes.

* **High performance:** Generating a static site typically is a matter of milliseconds. To keep build times short,
verless lets you choose and explicitly enable features you need - only generate RSS feeds or overview pages for tags if
you want to.

* **Focus on simplicity:** Global information, page types, plugins and other settings are configured in a central
configuration file inside your project. We try to keep things small and simple, and if your project isn't simple,
verless probably isn't a good fit.

## <img src="https://verless.dominikbraun.io/static/img/list-icon-v1.0.0.png"> Examples

* Example project structure: [example/](example)
* Real-world example website: [dominikbraun.io](https://dominikbraun.io)

## <img src="https://verless.dominikbraun.io/static/img/list-icon-v1.0.0.png"> Installation

### Linux and macOS

Download the [latest release](https://github.com/verless/verless/releases) for your platform. Extract the
downloaded binary into a directory like `/usr/local/bin`. Make sure the directory is in `PATH`.

### Windows

Download the [latest release](https://github.com/verless/verless/releases), create a directory like
`C:\Program Files\verless` and extract the executable into that directory.
[Add the directory to `Path`](https://www.computerhope.com/issues/ch000549.htm).

### Docker

Assuming your project directory is `my-blog`, the following command builds your website:

```shell script
$ docker container run -v $(pwd)/my-blog:/project verless/verless
```

The container will build the project mounted at `/project` and you'll find the website in `my-blog/target`. To run
another command, just append it to the image name:

```shell script
$ docker container run verless/verless version
```

## <img src="https://verless.dominikbraun.io/static/img/list-icon-v1.0.0.png"> Getting started

The easiest way to create a new project is to use the verless CLI:

```shell script
$ verless create project my-blog
```

This initializes a project called `my-blog` inside a new directory, containing a small default site. You can either
build the project or serve the static site directly:

```shell script
$ verless serve -w my-blog
```

After running the command, you can view your new project under [localhost:8080](http://localhost:8080). Building the
project works similary and generates a deployable website:

```shell script
$ verless build my-blog
```

By default, verless generates the website into `my-blog/target`. You're now good to [create your first content](docs)!

## <img src="https://verless.dominikbraun.io/static/img/list-icon-v1.0.0.png"> Documentation

To find out how a verless project is structured, take a look at the [example project](example). For a detailed
reference, check out the [documentation](docs).

## <img src="https://verless.dominikbraun.io/static/img/list-icon-v1.0.0.png"> Contributing

Anyone is welcome to contribute to verless. Please refer to our [contribution guidelines](CONTRIBUTING.md).

<p align="center">
<br>
<a href="https://github.com/verless/verless">
<img src="https://verless.dominikbraun.io/static/img/logo-footer-v1.0.0.png">
</a>
</p>
