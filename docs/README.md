# verless Documentation

## Contents

* [Configuration reference]()
* [Template reference]()
* [Command reference]()

## Configuration reference

All project-related settings are stored in the verless configuration file. Even though the file is optional from a
technical perspective, it makes sense to maintain a project configuration for the majority of cases.

### Configuration file

verless expects the project configuration file to be stored in the project root. The file name has to be `verless`, and
the file extension is `.yml`, `.toml` or `.json`, depending on the configuration format you want to use.

This documentation uses YAML for examples.

### Full configuration example

This is a full YAML configuration with all available configuration keys:

```yaml
site:
  meta:
    title: <Global website title>
    subtitle: <Global website subtitle>
    description: <Global website description>
    author: <Global website author>
    base: <Website base URL>
  nav:
    items:
      - label: <Navigation item label>
        target: <Navigation item URL>
      - ...
    overwrite: <Overwrite generated navigation>
  footer:
    items:
      - label: <Footer item label>
        target: <Footer item URL>
      - ...
    overwrite: <Overwrite generated footer>
plugins:
  - atom
  - ...
build:
  overwrite: <Overwrite output directory>
```

### Configuration key reference

...

<p align="center">
<br>
<a href="https://github.com/verless/verless"><img src="https://verless.dominikbraun.io/assets/img/icon-light.png"></a>
</p>
