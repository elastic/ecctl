# How to add a new command to ecctl

This document aims to provide a detailed guide on adding a new command to ecctl.

## The `cmd` package

The [`cmd`](../cmd) package contains all the commands that the user can see in the CLI. We are using the [cobra](https://github.com/spf13/cobra)
library to build these. This package is in charge of building the "scaffolding" of a command. It parses the arguments, calls
into the `app` package for business logic and returns results.

[`cmd/root.go`](../cmd/root.go) contains the root command (`ecctl`) and global flags.
[`cmd/commands.go`](../cmd/commands.go) attaches the top level commands (`deployment`, `platform`, etc) to the root command.
The subdirectories define the subcommand structure (e.g. `platform` -> `proxy`). The lowest level `command.go`
(e.g. [`cmd/platform/proxy/command.go`](../cmd/platform/proxy/command.go)) contains the definition of the actual command
(e.g. `platform proxy list`).

### Example: adding a new command to list proxies

Create a new subdirectory `cmd/platform/proxy`, and add a `command.go` file to it. (Some of these files/directories
might already exist, e.g. if you are adding a new operation for an existing entity.) You will define your new command here.

First, you need to define the `proxy` command:

```go
var Command = &cobra.Command{
    Use: "proxy",
    // ...
}
```

Next, define the `list` subcommand:

```go
var listProxiesCmd = &cobra.Command{
    Use:     "list",
    // ...
}
```

There are a few things to note about the `list` command:

- The `list` command will contain a call to the `proxy.List()` method from the `pkg/platform/proxy` package which you have
not defined yet. You will define this method later. You will also need to define the parameters struct used to pass
arguments to this method - e.g. `proxy.ListParams`.
- Ideally, you would implement a custom text formatter to present the proxy list in a user-friendly way. We currently
support two types of output: json and text. Absent a custom formatter, the output defaults to *json*. If you choose to add
a formatter to support *text* output, you'll need to create a template such as
[`pkg/formatter/templates/text/proxy/list.gotmpl`](../pkg/formatter/templates/text/proxy/list.gotmpl). To ensure this template
gets used, invoke a `ecctl.Get().Formatter.Format()` from the  [`listProxiesCmd`](https://github.com/elastic/ecctl/blob/a90daa0c4411905c8d5c3fa06f5b6250395c4730/cmd/platform/proxy/command.go#L60).

Next, the `init()` function attaches the `list` subcommand to the `proxy` command, and adds any custom
command line parameters.

```go
func init() {
    Command.AddCommand(listProxiesCmd)
    // ...
}

```

Finally, you need to attach the `proxy` subcommand to the `platform` command, by adding a reference to it
in [`cmd/platform/platform.go`](../cmd/platform/platform.go).

## The `pkg` package

The [`pkg`](../pkg) package is used by the `cmd` package to create commands and is library code that's ok to use by external applications.

### [`deployment`](../pkg/deployment)

Business logic behind `deployment` commands. Typically, methods in this package are called from the `cmd` package,
they accept parameters and call into the ES Cloud API, process and return results to the commands.

### [`ecctl`](../pkg/ecctl)

Contains the business logic and configuration for the ecctl app.

### [`formatter`](../pkg/formatter)

Functions and templates to format command output from json into user friendly text.

### [`platform`](../pkg/platform)

Business logic behind `platform` commands. Typically, methods in this package are called from the `cmd` package,
they accept parameters and call into the ES Cloud API, process and return results to the commands.

### [`util`](../pkg/util)

Common resources, such as utility functions, constants, parameters, that are shared among different packages.

### Example: adding a new command to list proxies (application logic)

Next, you need to add the business logic to your `platform proxy list` command. That's the `List()`
method mentioned earlier. It should go into the `pkg/platform/proxy/proxy.go` file (that you'll need create if it does
not already exist):

```go
func List(params ListParams) (*models.ProxyOverview, error) {
    // ...
}
```

Things to note about the `List` function:

- This is where the main business logic happens: it should include a call to the cloud API to retrieve all the
proxies and return them to the calling function.
- Make sure to properly catch and handle all possible errors. Consider using the  `multierror`
library if that makes sense. In this simple case, where we only have one API call, and not much else in terms of processing,
it's probably fine to use regular `errors`.
- For a general set of guidelines on how to write good Go code, see our [Style Guide](https://github.com/elastic/ecctl/blob/master/developer_docs/STYLEGUIDE.md).

You'll also need to define the `ListParams` struct, as we normally use structs to pass several arguments to functions. If you only have one or two
functions and corresponding parameter structs, it's ok to define the parameters structs in the same file. If it starts growing beyond that,
a good practice is to define parameter structs in their own file - in this case it would be `pkg/platform/proxy/proxy_params.go`.

```go
type ListParams struct {
    // ...
}
```

Make sure you also define the `Validate()` method for `ListParams`:

```go
func (params ListParams) Validate() error {
    // ...
}
```

Additionally, please create unit tests for your changes. Both the `List()` function and the `Validate()`
method for `ListParams` need to be tested. These tests should go to `pkg/platform/proxy/proxy_test.go`
(and `pkg/platform/proxy/proxy_params_test.go`, if you separated the params in their own file).
There are many examples of this in the code base, so feel free to browse and use them as inspiration.

That concludes all the steps necessary for creating a new command. You can easily manually test your new command by
making use of our [helper scripts](../CONTRIBUTING.md#helpers):

`dev-cli --config ~/path/to/your/ecctl/config.yaml platform proxy list`

If your command behaves as expected, all that's left is to make sure you followed all the
[code contribution guidelines](../CONTRIBUTING.md#code-contribution-guidelines) before submitting your PR.

Congratulations on writing your first ecctl command!
