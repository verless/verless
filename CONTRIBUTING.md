# Contribution guidelines

## Contents

* [Reporting issues](#reporting-issues)
* [Proposing features](#proposing-features)
* [Improving documentation](#improving-documentation)
* [Contributing code](#contributing-code)

## Reporting issues

If you encountered an unexpected behavior or a bug, feel free to
[file an issue](https://github.com/verless/verless/issues/new). When you do so, please make sure to ...
* include version information from the output of `verless version`.
* provide steps to reproduce the behavior.
* hide very long stack traces or logs in a [spoiler](https://gist.github.com/jbsulli/03df3cdce94ee97937ebda0ffef28287).

## Proposing features

Desire a new verless feature? Just propose your idea by
[creating an issue](https://github.com/verless/verless/issues/new) for it, there's nothing to lose!

Good feature proposals ...
* explain the problem that the feature solves.
* explain why it would be a desirable feature for the majority of users.
* bear in mind that verless should stay simple and lightweight.

## Improving documentation

If you've noticed something in the documentation that could be improved, please give us feedback. The easiest way to
request changes is to [file an issue](https://github.com/verless/verless/issues/new) or to fork the repository, make the
changes yourself and open a pull request. When you do so, please take note of our [Git conventions](#git-conventions).

## Contributing code

### Philosophy

Since verless isn't a company governed by customer expectations that have to be met for generating revenue, we try to
avoid quick-and-dirty solutions, workarounds and hotfixes at all costs. Instead, we strive for perfection,
maintainability and a thoughtful design. Working on the verless codebase should be fun!

### Git conventions

* Commit messages have to start capitalized and should be written in the imperative, e.g.
`Read markdown files concurrently`.
* If your pull request fixes an existing issue, refer to it with the issue number: `Fixes #19.`
* You may create a Draft Pull Request so reviewers can track your work and jump in early if there are problems.

### Coding conventions

Most importantly, there are some hard rules for code:

* All code has to follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines.
* All code has to be formatted with `gofmt -s`.
* All exported types, methods and variables have to be documented briefly.

Some general recommendations:

* Avoid OOP and global state if possible.
* Prefer standalone functions that accept an input and return an output.
* Prefer short and concise variable names.
* Complex logic and edge cases in your code should be explained.

We are very supportive and helpful, so don't hesitate to ask any questions during your work. Thanks for contributing!

<p align="center">
<br>
<a href="https://github.com/verless/verless">
<img src="https://verless.dominikbraun.io/static/img/logo-footer-v1.0.0.png">
</a>
</p>
