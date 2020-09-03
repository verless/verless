# Template reference

verless renders Markdown content as static HTML sites using templates. This reference explains how to create and use
these templates.

## Contents

* [Template path](#template-path)
* [Required template](#required-templates)
* [Optional templates](#optional-templates)
* [Template syntax](#template-syntax)
* [Field reference](#field-reference)

## Template path

All templates have to live in your project's `templates` directory. You can use the `verless create project` command
to initialize a new verless project, which will create the directory automatically for you.

## Required templates

There are two types of templates, and each is represented with its own file:

* `page.html`: This template is used to render your Markdown content.
* `index-page.html`: This template is used to render generated index pages. Verless creates an index page for each
directory in your content path, and all pages inside that directory are available to the index page.

Both templates are required to build a static website.

## Optional templates

If you've got a special page like an 'About' page, you may provide a distinct template in the Markdown file:

```markdown
# File: content/about.md
---
Title: About
Template: about.html
---
```

You can then provide a template called `about.html` in your `templates` directory. verless will use the custom template
instead of `page.html`.

## Template syntax

As verless is written in Go, it uses the [Go template syntax](https://golang.org/pkg/text/template/).

### Accessing fields

Fields can be accessed by a dot followed by the capitalized field name:

```html
<h1>{{.Title}}</h1>
``` 

### Looping

You can loop over an array like an index page's `Pages` and access its items as follows:

```html
<li>
    {{range $page := .Pages}}
        <li>{{$page.Title}}</li>
    {{end}}
</li>
```

Make sure to check out the [example templates](https://github.com/verless/verless/tree/master/example/templates).

## Field reference

### Meta

Available in:
* `page.html`
* `index-page.html`

| Field                   | Source      | Description                                                                                    |
|-------------------------|-------------|------------------------------------------------------------------------------------------------|
| `{{.Meta.Title}}`       | verless.yml | See [example/verless.yml](https://github.com/verless/verless/blob/master/example/verless.yml). |
| `{{.Meta.Subtitle}}`    | verless.yml | See [example/verless.yml](https://github.com/verless/verless/blob/master/example/verless.yml). |
| `{{.Meta.Description}}` | verless.yml | See [example/verless.yml](https://github.com/verless/verless/blob/master/example/verless.yml). |
| `{{.Meta.Author}}`      | verless.yml | See [example/verless.yml](https://github.com/verless/verless/blob/master/example/verless.yml). |
| `{{.Meta.Base}}`        | verless.yml | See [example/verless.yml](https://github.com/verless/verless/blob/master/example/verless.yml). |

### Nav

Available in:
* `page.html`
* `index-page.html`

| Field            | Source      | Description                                                                                    |
|------------------|-------------|------------------------------------------------------------------------------------------------|
| `{{.Nav.Items}}` | verless.yml | See [example/verless.yml](https://github.com/verless/verless/blob/master/example/verless.yml). |

### NavItem

Available in:
* `{{.Nav.Items}}`

| Field         | Source      | Description                                                                                    |
|---------------|-------------|------------------------------------------------------------------------------------------------|
| `{{.Label}}`  | verless.yml | See [example/verless.yml](https://github.com/verless/verless/blob/master/example/verless.yml). |
| `{{.Target}}` | verless.yml | See [example/verless.yml](https://github.com/verless/verless/blob/master/example/verless.yml). |

### Page

Available in:
* `page.html`
* `index-page.html`

| Field                   | Source   | Description                                                                                      |
|-------------------------|----------|--------------------------------------------------------------------------------------------------|
| `{{.Page.Route}}`       | Filepath | Page path in the form `/my-blog/coffee`. Useful for creating links to other pages.               |
| `{{.Page.ID}}`          | Filename | Useful for creating links to other pages.                                                        |
| `{{.Page.Title}}`       | Markdown |                                                                                                  |
| `{{.Page.Author}}`      | Markdown | For the global website author, see `{{.Meta.Author`.                                             |
| `{{.Page.Date}}`        | Markdown |                                                                                                  |
| `{{.Page.Tags}}`        | Markdown | Array of strings. You can loop through tags with `{{range $t := .Page.Tags}} ... {{end}}`.       |
| `{{.Page.Img}}`         | Markdown | It is recommended to use an URL like `/assets/img/picture.jpg`.                                  |
| `{{.Page.Credit}}`      | Markdown | This may be the image credit or something related.                                               |
| `{{.Page.Description}}` | Markdown |                                                                                                  |
| `{{.Page.Content}}`     | Markdown |                                                                                                  |
| `{{.Page.Related}}`     | Markdown | Array of `Page`. You can loop through tags with `{{range $r := .Page.Related}} ... {{end}}`.     |
| `{{.Page.Type}}`        | Markdown | An optional page type. Has to be declared in `verless.yml` (see `types` key) first.              |

### Pages

Available in:
* `index-page.html`

| Field        | Source   | Description                                                                                  |
|--------------|----------|----------------------------------------------------------------------------------------------|
| `{{.Pages}}` | Markdown | Array of `Page`. You can loop through tags with `{{range $r := .Page.Related}} ... {{end}}`. |

### Footer

Available in:
* `page.html`
* `index-page.html`

| Field               | Source      | Description                                                                                    |
|---------------------|-------------|------------------------------------------------------------------------------------------------|
| `{{.Footer.Items}}` | verless.yml | See [example/verless.yml](https://github.com/verless/verless/blob/master/example/verless.yml). |

### FooterItem

Available in:
* `{{.Footer.Items}}`

| Field         | Source      | Description                                                                                    |
|---------------|-------------|------------------------------------------------------------------------------------------------|
| `{{.Label}}`  | verless.yml | See [example/verless.yml](https://github.com/verless/verless/blob/master/example/verless.yml). |
| `{{.Target}}` | verless.yml | See [example/verless.yml](https://github.com/verless/verless/blob/master/example/verless.yml). |

<p align="center">
<br>
<a href="https://github.com/verless/verless"><img src="https://verless.dominikbraun.io/assets/img/icon-light.png"></a>
</p>
