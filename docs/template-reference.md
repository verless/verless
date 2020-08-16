# Template reference

verless renders Markdown content as static HTML sites using templates. This reference explains how to create and use
these templates.

## Template path

All templates have to live in your project's `templates` directory. You can use the `verless create project` command
to initialize a new verless project, which will create the directory automatically for you.

## Required templates

There are two types of templates, and each is represented with its own file:

* `page.html`: This template is used to render your Markdown content.
* `index-page.html`: This template is used to render generated index pages. verless creates an index page for each
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

Make sure to check out the [example templates](https://github.com/verless/verless/tree/docs/example/templates).

## Field reference

### Meta

Available in:
* `page.html`
* `index-page.html`

| Field                 | Source      | Note |
|-----------------------|-------------|------|
| {{.Meta.Title}}       | verless.yml |      |
| {{.Meta.Subtitle}}    | verless.yml |      |
| {{.Meta.Description}} | verless.yml |      |
| {{.Meta.Author}}      | verless.yml |      |
| {{.Meta.Base}}        | verless.yml |      |
| {{.Meta.Subtitle}}    | verless.yml |      |

...
