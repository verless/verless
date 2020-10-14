# Plugin reference

The philosophy of verless is that it provides a very basic feature set by default - it basically just renders your
Markdown content using HTML templates -, and you explicitly have to enable more advanced features in your project
configuration.

Those advanced features are provided by _plugins_.

## Contents

* [Enabling plugins](#enabling-plugins)
* [Available plugins](#available-plugins)

## Enabling plugins

To enable a certain plugin, you have to add them to the top-level `plugins` key in your project configuration. Say you
have some blog posts, and they're tagged:

```markdown
---
Title: Choosing a coffee machine
Tags:
    - Coffee
    - Coffee Machine
---
```

To enable overview pages for the `Coffee` and `Coffee Machine` tags, you have to add the `tags` plugin to your
configuration:

```yaml
plugins:
- tags
```

That's it! For a full configuration example, see the [example project](../example/verless.yml).

## Available plugins

As verless has just been released, we're constantly adding new plugins.

### atom

* **Plugin key:** `atom`
* **What it does:** Generates an Atom RSS feed for your pages. You can exclude a page with `Hide: true`. The generated
RSS feed will be available in your project root.

### related

* **Plugin key:** `related`
* **What it does:** Allows you to put page URIs in the `Related` section of your Markdown Front Matter and access those
related pages as [`{{.Page.Related}}`](template-reference.md#page) in templates.

```markdown
---
Title: Making Barista-Quality Espresso
Related:
    - /blog/steaming-milk-for-cappuccino
---

...
```

The page URI is a verless path inside the `content` directory - in the example above, the related page physically is
`content/blog/steaming-milk-for-cappuccino.md`.

The `{{.Page.Related}}` array contains full `Page` instances with _all_ page data available.

### tags

* **Plugin key:** `tags`
* **What it does:** Creates a top-level `tags` directory containing a directory for each tag, and those directories
contain a list page rendered with your [`list-page.html`](template-reference.md#required-templates) template. This is
where all pages are available as [`Pages`](template-reference.md#pages). From there, you can link to the each page's
actual location. As a result, the overview for all articles with the `coffee` tag are available under `/tags/coffee`.

<p align="center">
<br>
<a href="https://github.com/verless/verless">
<img src="https://verless.dominikbraun.io/static/img/logo-footer-v1.0.0.png">
</a>
</p>
