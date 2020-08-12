<p align="center">
<br>
<br>
<br>
<img src="https://verless.dominikbraun.io/assets/img/logo-github.png">
<br>
<br>
<br>
</p>

<h3 align="center">A simple and lightweight Static Site Generator.</h3>

<p align="center">
<a href="https://circleci.com/gh/dominikbraun/cleanup"><img src="https://circleci.com/gh/dominikbraun/cleanup.svg?style=shield"></a>
<a href="https://goreportcard.com/report/github.com/dominikbraun/cleanup"><img src="https://goreportcard.com/badge/github.com/dominikbraun/cleanup"></a>
<a href="https://www.codefactor.io/repository/github/dominikbraun/cleanup"><img src="https://www.codefactor.io/repository/github/dominikbraun/cleanup/badge" /></a>
<a href="https://github.com/dominikbraun/cleanup/releases"><img src="https://img.shields.io/github/v/release/dominikbraun/cleanup?sort=semver"></a>
<a href="https://github.com/dominikbraun/cleanup/blob/master/LICENSE"><img src="https://img.shields.io/badge/license-Apache--2.0-brightgreen"></a>
<br>
<br>
<br>
</p>

---

**verless** (pronounced like _serverless_) is a Static Site Generator designed for Markdown-based content with focus on
simplicity and performance. It reads your Markdown files, applies your HTML templates and renders them as a website.

## <img src="https://verless.dominikbraun.io/assets/img/dot.png"> Features

* Flexibility: Provide and use your default template for all pages, or use a different template for a specific page
* Simplicity: Build your entire website within a single CLI command
* Performance: Generating your website only takes a few seconds, even for thousands of pages
* Rapid Development: Get started quickly with verless' small and reduced feature set
* Configurability: Provide additional information or override defaults in `verless.yml`
* Portability: verless is packaged as a single binary without any dependencies for multiple platforms

## <img src="https://verless.dominikbraun.io/assets/img/dot.png"> Examples

* Example project structure: [verless/example](https://github.com/verless/example)
* Real-world example website: [dominikbraun.io](https://dominikbraun.io)

## <img src="https://verless.dominikbraun.io/assets/img/dot.png"> Installation

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

## <img src="https://verless.dominikbraun.io/assets/img/dot.png"> Getting started

If you are ready to get your hands dirty and test if everything works as expected, you may want to start with a quick
example. Create a project directory on your local machine, e.g. `coffee-blog`, and a file called `verless.yml` inside of
it:

```yaml
# File: coffee-blog/verless.yml
site:
  meta:
    title: Coffee Blog
    subtitle: About Espresso & Cappuccino
```

This file stores global metadata for your website and offers the possibility to modify your website manually, for
example you could override the generated navigation items here.

After that, create a directory `templates/` inside your project and add a template file called `index-page.html`:

```html
<!-- File: coffee-blog/templates/index-page.html -->
<!DOCTYPE html>
<html>
    <head>
        <title>{{.Meta.Title}}</title>
    </head>
    <body>
        <h1>{{.Meta.Title}}</h1>
        <h3>{{.Meta.Subtitle}}</h3>
    </body>
</html>
```

Now run the following command in your project directory:

```shell script
$ verless build
```

The generated website is inside the created `target` directory. You're now good to
[create your first content](https://github.com/verless/verless/tree/master/docs)!

## <img src="https://verless.dominikbraun.io/assets/img/dot.png"> Documentation

For a detailed reference, check out the [documentation](https://github.com/verless/verless/tree/master/docs).

## <img src="https://verless.dominikbraun.io/assets/img/dot.png"> Contributing

Anyone is welcome to contribute to verless. Please refer to our
[Contribution guidelines](https://github.com/verless/verless/blob/master/.github/CONTRIBUTING.md).

<p align="center">
<img src="https://verless.dominikbraun.io/assets/img/icon-light.png">
</p>
