package core

var (
	defaultConfig = []byte(`version: 1
site:
  meta:
    title: Your verless Project
theme: default
build:
  # overwrite: true
`)

	defaultTpl = []byte(`<!DOCTYPE html>
<html lang="en">
    <head>
        <title>{{.Meta.Title}}</title>
        <link rel="stylesheet" href="/css/style.css" />
    </head>
    <body>
        <main>
            <img src="https://verless.dominikbraun.io/static/img/logo-default.png"
                 alt="verless" id="logo" />
            <h1>{{.Meta.Title}}</h1>
            <p>
                Welcome to your new project! Take a look at the <a
                        href="https://github.com/verless/verless/tree/master/example" target="_blank">example
                    project</a> or visit the <a
                        href="https://github.com/verless/verless/tree/master/docs" target="_blank">documentation</a>.
            </p>
        </main>
    </body>
</html>`)

	defaultCss = []byte(`* {
    font-family: Arial, Tahoma, sans-serif;
    color: #32343D;
}

body {
    padding: 4rem 2rem;
}

main {
    text-align: center;
}

h1 {
    padding-bottom: 2rem;
    font-weight: normal;
}

#logo {
    margin: 2rem;
}`)

	defaultThemeConfig = []byte(`version: 1
build:
  # Here you can specify commands you need to run in order to build
  # your theme, like generating CSS.
  before:
    # - Your command here - just comment out the line
    # - Another command there
`)

	defaultGitignore = []byte(`generated/`)
)
