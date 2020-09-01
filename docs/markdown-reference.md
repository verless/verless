# Markdown reference

While templates are responsible for displaying content, the actual content lives in Markdown files inside your `content`
directory.

## Contents

* [Paths and filenames](#paths-and-filenames)
* [Metadata](#metadata)
* [Front Matter reference](#front-matter-reference)

## Paths and filenames

A content file has to meet the following requirements:
* It is stored inside the `content` directory of your project.
* It is a Markdown file with the `.md` extension.

Each file in the `content` directory will be converted to a
[Page](https://github.com/verless/verless/blob/master/docs/template-reference.md#page).

The path of the Markdown file inside `content` defines the route for the corresponding page: A Markdown file stored as
`content/blog/making-barista-quality-espresso.md` will be converted to a page whose URL is
`/blog/making-barista-quality-espresso`.

* The path and name of a Markdown file directly defines its URL on the website.
* Paths and names must not contain spaces.

## Metadata

While the URL for a page is inferred from its filename, other metadata is parsed from the Markdown file. Verless uses
the _YAML Front Matter_ for this: There's a short YAML section at the beginning of the Markdown file.

```markdown
---
Title: Making Barista-Quality Espresso
Description: This is a guide for making italian Espresso.
Tags:
    - Espresso
    - Coffee
---

Do you enjoy a high-quality italian Espresso as much as I do?
```

For broader examples, check out the
[example project](https://github.com/verless/verless/tree/master/example/content/blog).

## Front Matter reference

This reference shows all available YAML keys for providing metadata. **All keys have to be capitalized.**

* **`Title`** _(String)_: The page's title.
* **`Author`** _(String)_: The page's author.
* **`Date`** _(String)_: The creation date in the form `YYYY-MM-DD`.
* **`Tags`** _(Array)_: A list of page tags. Enable the [tags plugin](https://github.com/verless/verless/blob/master/docs/plugin-reference.md#tags) for tag support.
    - **`<tag>`** _(String)_: A page tag.
* **`Img`** _(String)_: An image URL like `assets/img/image.jpg`.
* **`Credit`** _(String)_: Copyright credit for `Img` or other contents.
* **`Description`** _(String)_: The page's description.
* **`Related`** _(Array)_: A list of related pages. Has to contain verless paths like `/blog/making-barista-quality-espresso`. This list will be available as `{{.Related}}` in the `page.html` template and contains [Page](https://github.com/verless/verless/blob/master/docs/template-reference.md#page) instances.
    - **`<verless path>`** _(String)_: The path to a related page.
* **`Type`** _(String)_: The page's content type which automatically defines the template to be used. If the content type is `blog`, a template called `blog.html` will be used for rendering. In future verless versions, content types may provide more configurations other than just the template to use.
* **`Hidden`** _(Bool)_: Don't include the page in lists like [`{{.Pages}}`](https://github.com/verless/verless/blob/master/docs/template-reference.md#pages).

<p align="center">
<br>
<a href="https://github.com/verless/verless"><img src="https://verless.dominikbraun.io/assets/img/icon-light.png"></a>
</p>
