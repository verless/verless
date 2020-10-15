# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to
[Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.5.0] - 2020-10-13

### Added
- Introduce the `related` plugin.
- Introduce page types for themes.

### Deprecated
- Deprecate the `types` section in `verless.yml`.
- Deprecate the `css` and `js` directories in themes.

### Fixed
- Fix `serve` behavior when an error occurs.

## [0.4.8] - 2020-10-10

### Changed
- Create directories for files created with `verless create file`.
- Improve performance of copying files.
- Improve log messages.

### Deprecated
- Deprecate `css` and `js` directories in themes.

## [0.4.7] - 2020-10-07

### Changed
- Panic when there are invalid list pages.
- Introduce the `assets` directory for themes.

## [0.4.6] - 2020-10-04

### Fixed
- Fix generated route of list page for `/tags`.

## [0.4.5] - 2020-10-03

### Changed
- Introduce configurations for themes.
- Introduce `--project` flag for `verless create theme` command.
- Create default `.gitignore` file for new verless projects.
- Make `verless create file` command available.

### Fixed
- Fix generated routes for `tags` plugin.
- Fix initialization of list pages for tags.
- Fix `version` key in configuration reference.

## [0.4.4] - 2020-10-02

### Changed
- Make the paths for `create project` and `serve` optional.

### Fixed
- Fix `create project` inside the current working directory.

## [0.4.3] - 2020-09-28

### Fixed
- Fix public verless Docker image.

## [0.4.2] - 2020-09-28

### Changed
- Make pre-build hooks available in project configuration.
- Ignore `static/generated` directory when watching the site.
- Improve the verless documentation.

## [0.4.1] - 2020-09-25

### Changed
- Ignore hidden pages in list pages.
- Publish verless Docker images on Docker Hub and GitHub Packages.

## [0.4.0] - 2020-09-23

### Added
- Introduce the verless theming system.
- Add the `verless create theme` command.

### Changed
- Create a default theme when creating a new project.
- Update the example project to use the default theme.
- Rename the `assets` directory to `static`.

### Removed
- Remove the `templates` directory.

## [0.3.9] - 2020-09-22

### Changed
- Make entire page URL available as `Page.Href`.

### Fixed
- Fix wrong URLs for pages in the root directory.

## [0.3.8] - 2020-09-22

### Changed
- Generate `list-page.html` even for empty routes.

## [0.3.7] - 2020-09-20

### Changed
- `list-page.html`: Make pages from sub-directories available.

## [0.3.6] - 2020-09-17

### Fixed
- Fix relative build paths like `./example`.
- Fix error when serving the site without `-w`.

## [0.3.5] - 2020-09-17

### Changed
- Use relative links in the documentation.
- Don't create physical files when testing.
- Order pages in `ListPage.Pages` by date.

## [0.3.4] - 2020-09-16

### Changed
- Don't include hidden pages in RSS feeds.

### Fixed
- Correctly assign routes to list pages.

## [0.3.3] - 2020-09-16

### Changed
- Make `Hidden` field of pages available in templates.

### Deprecated
- Remove the `Template` field in Markdown files.

### Fixed
- Fix incorrect routes for `.` as build path.

## [0.3.2] - 2020-09-15

### Changed
- Print short description for commands in the CLI.

### Fixed
- Render files in correct target directory when using `verless serve`.

## [0.3.1] - 2020-09-14

### Changed
- Change 'index page' terminology to 'list page'.
- Make list pages overridable by providing an own `index.md` file.

## [0.3.0] - 2020-09-12

### Added
- Add the `verless serve` command.
- Add file watching with automatic rebuilds for `serve`.

## Older releases

* [0.2.2](https://github.com/verless/verless/releases/tag/v0.2.2)
* [0.2.1](https://github.com/verless/verless/releases/tag/v0.2.1)
* [0.2.0](https://github.com/verless/verless/releases/tag/v0.2.0)
* [0.1.7](https://github.com/verless/verless/releases/tag/v0.1.7)
* [0.1.6](https://github.com/verless/verless/releases/tag/v0.1.6)
* [0.1.5](https://github.com/verless/verless/releases/tag/v0.1.5)
* [0.1.4](https://github.com/verless/verless/releases/tag/v0.1.4)
* [0.1.3](https://github.com/verless/verless/releases/tag/v0.1.3)
* [0.1.2](https://github.com/verless/verless/releases/tag/v0.1.2)
* [0.1.1](https://github.com/verless/verless/releases/tag/v0.1.1)
* [0.1.0](https://github.com/verless/verless/releases/tag/v0.1.0)

<p align="center">
<br>
<a href="https://github.com/verless/verless">
<img src="https://verless.dominikbraun.io/static/img/logo-footer-v1.0.0.png">
</a>
</p>
