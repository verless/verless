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

We highly appreciate the time and effort that go into code contributions. There are just some requirements the code
needs to meet in order to get merged.

### Git conventions

* Branch names are up to you.
* Commit messages have to start capitalized and should be written in the imperative, e.g.
`Read markdown files concurrently`.
* Create a WIP pull request so that reviewers can track your work continuously and jump in early if there are problems.
* If your pull request fixes an existing issue, refer to it with the issue number: `Fixes #19.`

### Coding conventions

Most importantly, there are some hard rules for code:

* All code has to be testable, maintainable and extendable.
* All code has to follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines.
* All code has to be formatted with `gofmt -s`.
* All exported types, methods and variables have to be documented briefly.
* All code has to pass the CI jobs successfully.

There also are some 'soft' recommendations that apply to most cases:

* Avoid OOP and global state.
* Prefer standalone functions that accept an input and return an output.
* Prefer immutability if it doesn't make the code harder to reason about.
* Make use of closures.
* Prefer short and concise variable names.

**Thanks for contributing!**
