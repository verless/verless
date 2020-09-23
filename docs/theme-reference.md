# Theme reference

In order to render your Markdown content, verless needs HTML templates. Those templates are provided by _themes_, and
you therefore need to customize the default theme, create your own theme or use a third-party theme.

## Contents

* [Theme structure](#theme-structure)
* [Required templates](#required-templates)
* [Custom templates](#custom-templates)
* [Customize the default theme](#customize-the-default-theme)

## Theme structure

All themes are stored inside the `themes` directory in your project, and each theme has its own directory. Inside this
directory, there has to be a `templates` directory. The `css` and `js` directories are optional.

Let's take a theme called `dark-theme` as an example. The directory structure has to look as follows:

```shell script
my-blog/
    themes/
        dark-theme/
            css/
                style.css
            js/
            templates/
                list-page.html
                page.html
```

The `css` and `js` directories will be copied into the root of your website, so your stylesheet will be directly
available as `/css/style.css`, for example.

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

## Customize the default theme

When you create a new project using `verless create project`, verless generates a default theme inside the `themes`
directory.

```shell script
$ verless create project my-blog
$ cd my-blog/themes/default
$ ls -l
  Sep 23 08:18 css/
  Sep 23 10:02 templates/
```

You can customize this theme by changing the templates and stylesheets. To view your changes, re-build the project or
just use `verless serve -w my-blog`.

## Create your own theme

The simplest way to create your own theme is to use the `verless create theme` command. It expects two arguments: The
project for which you want to create a theme and the theme name itself.

For example, creating the `dark-theme` theme for the `my-blog` project is as easy as:

```shell script
$ verless create theme my-blog dark-theme
```

Or, if you already are _inside_ the `my-blog` directory:

```shell script
$ verless create theme . dark-theme
```