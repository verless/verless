/*
Package core provides the verless business logic.

Each verless sub-command maintains its own file in the core package,
like build.go for `verless build`. This file provides a function in
the form Run<Command>, e.g. RunBuild, serving as an entry point for
other packages.

An entry point function either implements the business logic itself
like RunVersion does, or prepares components and passes them to an
inner core package for more complex scenarios.

Typically, a verless command has multiple options. These options are
types in the core package, declared in the command's file. They have
to be initialized by the caller - for example the cli package - and
then passed to the entry point function. The name of an option type
must be in the form <Command>Options.

As a result, most entry point functions accept a dedicated options
instance. Some of them also require an initialized config instance
or fixed command line arguments - see RunBuild as an example.

Core functions typically shouldn't have outgoing dependencies except
for dependencies like the fs or model packages. Instead, they should
define their dependencies as interfaces which can be implemented by
outside packages. A good example for this is build.Parser implemented
by parser.Markdown.

It is the entry point function's job to initialize those external
dependencies, like calling parser.NewMarkdown, and passing it to the
particular core function. Again, RunBuild is good example for this.
*/
package core
