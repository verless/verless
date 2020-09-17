# Configuration reference

All project-related settings are stored in the verless configuration file. Even though the file is optional from a
technical perspective, it makes sense to maintain a project configuration for the majority of cases.

## Contents

* [Configuration file](#configuration-file)
* [Full configuration example](#full-configuration-example)
* [Configuration key reference](#configuration-key-reference)

## Configuration file

verless expects the project configuration file to be stored in the project root. The file name has to be `verless`, and
the file extension is `.yml`, `.toml` or `.json`, depending on the configuration format you want to use.

## Full configuration example

There is a full YAML configuration available in the example project: [example/verless.yml](../example/verless.yml)

Note that all configuration keys except `version` are optional.

## Configuration key reference

* **`site`** _(Map)_:
    * **`version`** _(String)_: The configuration version (currently `1`).
    * **`meta`** _(Map)_:
        * **`title`** _(String)_: The global website title that applies to all pages.
        * **`subtitle`** _(String)_: The global website subtitle that applies to all pages.
        * **`description`** _(String)_: The global website description that applies to all pages.
        * **`author`** _(String)_: The website author or publisher.
        * **`base`** _(String)_: The website's base URL in the form `https://example.com`. Needs to be enclosed in quotes.
    * **`nav`** _(Map)_:
        * **`items`** _(Array)_:
            * **`label`** _(String_): The navigation item's label, e.g. `Home`.  
              **`target`** _(String)_: The navigation item's target URL in the form `https://example.com`. Needs to be enclosed in quotes.
    * **`footer`** _(Map)_:
        * **`items`** _(Array)_:
            * **`label`** _(String_): The footer item's label, e.g. `Home`.   
              **`target`** _(String)_: The footer item's target URL in the form `https://example.com`. Needs to be enclosed in quotes.
* **`types`** _(Map)_:
    * **`<type>`** _(Object)_: A page type.
        * **`template`** _(String)_: The template to use for rendering pages of `<type>`.
* **`plugins`** _(Array)_:
    - **`<plugin key>`** _(String)_: The key of the plugin to be used. You can find the plugin key in the [plugin reference](#plugin-reference).
* **`build`** _(Map)_:
    * **`overwrite`** _(Bool)_: Allow verless to overwrite the output directory completely. This removes the need for the `--overwrite` flag for builds.
    
<p align="center">
<br>
<a href="https://github.com/verless/verless"><img src="https://verless.dominikbraun.io/assets/img/icon-light.png"></a>
</p>