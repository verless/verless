# Contribution guidelines

## Contents

* [Reporting issues](#reporting-issues)
* [Proposing features](#proposing-features)
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

## Contributing code

I highly appreciate the time and effort that go into code contributions. There are just some requirements the code needs
to meet in order to get merged.

### Git conventions

* Your branch names are up to you.
* Commit messages have to start capitalized and should be written in the imperative, e.g.
`Read markdown files concurrently`.
* Create a WIP pull request so that reviewers can track your work continuously and jump in early if there are problems.

### Coding conventions

Most importantly, there are some hard rules for code:

* All code has to be testable, maintainable and extendable.
* All code has to follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines.
* All code has to be formatted with `gofmt -s`.
* All type declarations and methods have to be documented, even if they're private.

There also are some 'soft' recommendations that apply to most cases:

* Avoid OOP and global state.
* Prefer standalone functions that accept an input and provide an output.
* Prefer immutability if it doesn't make the code harder to reason about.
* Make use of closures.
* Prefer short and concise variable names.

**Thanks for contributing!**
