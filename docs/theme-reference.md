# Theme reference

In order to render your Markdown content, verless needs HTML templates. Those templates are provided by _themes_, and
you therefore need to customize the default theme, create your own theme or use a third-party theme.

## Contents

* [Customize the default theme](#customize-the-default-theme)
* [Create your own theme](#create-your-own-theme)
* [Theme structure](#theme-structure)
* [Theme configuration](#theme-configuration)
* [Required templates](#required-templates)
* [Custom templates](#custom-templates)
* [Pre-build hooks](#pre-build-hooks)

## Customize the default theme

When you create a new project using `verless create project`, verless generates a default theme inside the `themes`
directory.

You can customize this theme by changing the templates and stylesheets. To view your changes, re-build the project or
just use `verless serve -w my-blog`.

## Create your own theme

The simplest way to create your own theme is to use the `verless create theme` command. It expects two arguments: The
project for which you want to create a theme and the theme name itself.

For example, creating the `dark-theme` theme for the `my-blog` project is as easy as:

```shell script
$ verless create theme -p my-blog dark-theme
```

You can leave out the `-p` flag if you already are _inside_ the `my-blog` directory:

```shell script
$ verless create theme dark-theme
```

## Theme structure

All themes are stored inside the `themes` directory in your project, and each theme has its own directory. Inside this
directory, there has to be a `templates` directory.

Let's take a theme called `dark-theme` as an example. The directory structure has to look as follows:

```shell script
my-blog/
└── themes/
    └── dark-theme/
        ├── theme.yml
        ├── assets/
        ├── generated/ (optional)
        │   └── css/
        │       └── style.css
        └── templates/
            ├── list-page.html
            └── page.html
```

Stylesheets, JavaScript files or even images can be stored in `assets`. This directory will be copied into the root of
your website, so the stylesheet from the example above is directly available as `/assets/css/style.css`. Any
theme-specific configuration goes into `theme.yml`.

**To activate your theme, set it in `verless.yml`:**

```yaml
theme: dark-theme
```

## Required templates

A theme requires two templates inside its `templates` directory:

* `page.html`: This template is used to render your Markdown content.
* `list-page.html`: This template is used to render generated list pages. Verless creates a list page for each
directory in your content path, and all pages inside that directory are available to the list page.

## Custom templates

If your theme offers any special pages, you may provide an additional template inside `templates`. To use this
template for a certain page, first define a new page type in the project configuration:

```yaml
# File: verless.yml

types:
   my-special-page:
      template: my-special-template.html
```

You can now use that type in your page:

```markdown
# File: content/about.md
---
Title: About
Type: my-special-page
---
```

## Pre-build hooks

Modern front-end development often requires preprocessing CSS or JS files, for example when using Sass for CSS. For
this purpose, verless offers _pre-build hooks_.

If it doesn't exist yet, create a `theme.yml` inside the directory of your theme. All you have to do is to add the
command to the `build` section:

```yaml
# File: theme.yml

build:
   before:
      - sass css/style.scss generated/style.css
```

**We highly recommend to store generated files like `style.css` inside a directory called `generated`.**

If you don't do this, `verless serve -w` will run into an infinite loop - because the build triggers the generation of
files, and the generated files will trigger another build because verless notices that something changed.

<p align="center">
<br>
<a href="https://github.com/verless/verless"><img src="https://verless.dominikbraun.io/assets/img/icon-light.png"></a>
</p>
