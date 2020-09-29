# Template reference

verless renders Markdown content as static HTML sites using templates. This reference explains how to use these
templates. If you want to know how to _create_ them, see the [theme reference](theme-reference.md).

## Contents

* [Template syntax](#template-syntax)
* [Field reference](#field-reference)

## Template syntax

As verless is written in Go, it uses the [Go template syntax](https://golang.org/pkg/text/template/).

### Accessing fields

Fields can be accessed by a dot followed by the capitalized field name:

```html
<h1>{{.Title}}</h1>
``` 

### Looping

You can loop over an array like a list page's `Pages` and access its items as follows:

```html
<li>
    {{range $page := .Pages}}
        <li>{{$page.Title}}</li>
    {{end}}
</li>
```

Make sure to check out the [example templates](../example/templates).

## Field reference

### Meta

Available in:
* `page.html`
* `list-page.html`
* Templates used by an `index.md` page

| Field                   | Source      | Description                                        |
|-------------------------|-------------|----------------------------------------------------|
| `{{.Meta.Title}}`       | verless.yml | See [example/verless.yml](../example/verless.yml). |
| `{{.Meta.Subtitle}}`    | verless.yml | See [example/verless.yml](../example/verless.yml). |
| `{{.Meta.Description}}` | verless.yml | See [example/verless.yml](../example/verless.yml). |
| `{{.Meta.Author}}`      | verless.yml | See [example/verless.yml](../example/verless.yml). |
| `{{.Meta.Base}}`        | verless.yml | See [example/verless.yml](../example/verless.yml). |

### Nav

Available in:
* `page.html`
* `list-page.html`
* Templates used by an `index.md` page

| Field            | Source      | Description                                        |
|------------------|-------------|----------------------------------------------------|
| `{{.Nav.Items}}` | verless.yml | See [example/verless.yml](../example/verless.yml). |

### NavItem

Available in:
* `{{.Nav.Items}}`

| Field         | Source      | Description                                        |
|---------------|-------------|----------------------------------------------------|
| `{{.Label}}`  | verless.yml | See [example/verless.yml](../example/verless.yml). |
| `{{.Target}}` | verless.yml | See [example/verless.yml](../example/verless.yml). |

### Page

Available in:
* `page.html`
* `list-page.html`
* Templates used by an `index.md` page

| Field                   | Source   | Description                                                                                                              |
|-------------------------|----------|--------------------------------------------------------------------------------------------------------------------------|
| `{{.Page.Href}}`        | Filepath | Ready to use path to the page for links.                                                                                 |
| `{{.Page.Route}}`       | Filepath | Page path in the form `/my-blog/coffee`. Useful for creating links to other pages. If possible, prefer `{{.Page.Href}}`. |
| `{{.Page.ID}}`          | Filename | Useful for creating links to other pages. If possible, prefer `{{.Page.Href}}`.                                          |
| `{{.Page.Title}}`       | Markdown |                                                                                                                          |
| `{{.Page.Author}}`      | Markdown | For the global website author, see `{{.Meta.Author`.                                                                     |
| `{{.Page.Date}}`        | Markdown |                                                                                                                          |
| `{{.Page.Tags}}`        | Markdown | Array of strings. You can loop through tags with `{{range $t := .Page.Tags}} ... {{end}}`.                               |
| `{{.Page.Img}}`         | Markdown | It is recommended to use an URL like `/assets/img/picture.jpg`.                                                          |
| `{{.Page.Credit}}`      | Markdown | This may be the image credit or something related.                                                                       |
| `{{.Page.Description}}` | Markdown |                                                                                                                          |
| `{{.Page.Content}}`     | Markdown |                                                                                                                          |
| `{{.Page.Related}}`     | Markdown | Array of `Page`. You can loop through tags with `{{range $r := .Page.Related}} ... {{end}}`.                             |
| `{{.Page.Type}}`        | Markdown | An optional page type. Has to be declared in `verless.yml` (see `types` key) first.                                      |
| `{{.Page.Hidden}}`      | Markdown |                                                                                                                          |

### Links to pages

Normally you should use `{{.Page.Href}}` as it already provides a ready to use file path.  
Concatenating `{{.Page.Route}}` with `{{.Page.ID}}` manually can lead to undesired effects and therefore this should be avoided.  
Example:  
`<p><a href="{{$page.Href}}">read post</a></p>`

### Pages

Available in:
* `list-page.html`
* Templates used by an `index.md` page

| Field        | Source   | Description                                                                                  |
|--------------|----------|----------------------------------------------------------------------------------------------|
| `{{.Pages}}` | Markdown | Array of `Page`. You can loop through tags with `{{range $r := .Page.Related}} ... {{end}}`. |

### Footer

Available in:
* `page.html`
* `list-page.html`
* Templates used by an `index.md` page

| Field               | Source      | Description                                        |
|---------------------|-------------|----------------------------------------------------|
| `{{.Footer.Items}}` | verless.yml | See [example/verless.yml](../example/verless.yml). |

### Footer

Available in:
* `{{.Footer.Items}}`

| Field         | Source      | Description                                        |
|---------------|-------------|----------------------------------------------------|
| `{{.Label}}`  | verless.yml | See [example/verless.yml](../example/verless.yml). |
| `{{.Target}}` | verless.yml | See [example/verless.yml](../example/verless.yml). |

<p align="center">
<br>
<a href="https://github.com/verless/verless">
<img src="https://verless.dominikbraun.io/static/img/logo-footer-v1.0.0.png">
</a>
</p>
