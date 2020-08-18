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

That's it! For a full configuration example, see the
[example project](https://github.com/verless/verless/blob/master/example/verless.yml).

## Available plugins

As verless has just been released, we're constantly adding new plugins.

### atom

* **Plugin key:** `atom`
* **What it does:** Generates an Atom RSS feed for your pages. You can exclude a page with `Hide: true`. The generated
RSS feed will be available in your project root.

### tags

* **Plugin key:** `tags`
* **What it does:** Creates a top-level `tags` directory containing a directory for each tag, and those directories
contain an index page rendered with your
[`index-page.html`](https://github.com/verless/verless/blob/docs/docs/template-reference.md#required-templates)
template. This is where all pages are available as
[`Pages`](https://github.com/verless/verless/blob/master/docs/template-reference.md#pages). From there, you can link to
the each page's actual location. As a result, the overview for all articles with the `coffee` tag are available under
`/tags/coffee`.
