# How to add a new command to ecctl

This document aims to provide a detailed guide on adding a new command to ecctl.

## The `cmd` package

The [`cmd`](../cmd) package contains all the commands that the user can see in the CLI. We are using the [cobra](https://github.com/spf13/cobra)
library to build these. This package is in charge of building the "scaffolding" of a command. It parses the arguments, calls
into the `app` package for business logic and returns results.

[`cmd/root.go`](../cmd/root.go) contains the root command (`ecctl`) and global flags.
[`cmd/commands.go`](../cmd/commands.go) attaches the top level commands (`deployment`, `platform`, etc) to the root command.
The subdirectories define the subcommand structure (e.g. `deployment` -> `list`). The lowest level
(e.g. [`cmd/deployment/list.go`](../cmd/deployment/list.go)) contains the definition of the actual command
(e.g. `deployment list`).

### Example: adding a new command to list deployments

Create a new subdirectory `cmd/deployment`, and add a `list.go` file to it. Some of these files/directories
might already exist (i.e. if you are adding a new operation for an existing entity). You will define your new command here.

Define the `list` subcommand in `cmd/deployment/list.go`:

```go
var listCmd = &cobra.Command{
    Use:     "list",
    // ...
}
```

There are a few things to note about the `list` command:

- The `list` command will contain a call to the `deploymentapi.List()` function from the [`/pkg/api/deploymentapi`](https://github.com/elastic/cloud-sdk-go/tree/master/pkg/api/deploymentapi)
package found in the cloud-sdk-go. You will define this function later. You will also need to define the parameters
struct used to pass arguments to this function - e.g. `deploymentapi.ListParams`.
- Ideally, you would implement a custom text formatter to present the deployment list in a user-friendly way. We currently
support two types of output: json and text. Absent a custom formatter, the output defaults to *json*. If you choose to add
a formatter to support *text* output, you'll need to create a template such as
[`pkg/formatter/templates/text/deployment/list.gotmpl`](../pkg/formatter/templates/text/deployment/list.gotmpl). To ensure this template
gets used, invoke a `ecctl.Get().Formatter.Format()` from the  [`listCmd`](https://github.com/elastic/ecctl/blob/master/cmd/deployment/list.go#L39).

Next, the `init()` function attaches the `list` subcommand to the `deployment` command, and adds any custom
command line parameters.

```go
func init() {
    Command.AddCommand(listCmd)
    // ...
}

```

Finally, you need to attach the `deployment` subcommand to the `ecctl` command, by adding a reference to it
in [`cmd/commands.go`](../cmd/commands.go).

## The `pkg` package

The [`pkg`](../pkg) package is used by the `cmd` package to create commands and is library code that's ok to use by external applications.

### [`ecctl`](../pkg/ecctl)

Contains the business logic and configuration for the ecctl app.

### [`formatter`](../pkg/formatter)

Functions and templates to format command output from json into user friendly text.

### [`util`](../pkg/util)

Common resources, such as utility functions, constants, parameters, that are shared among different packages.

### Example: adding a new command to list deployments (application logic)

Next, you need to add the business logic to your `ecctl deployment list` command. That's the `List()`
function mentioned earlier. These APIs can be found in [cloud-sdk-go](https://github.com/elastic/cloud-sdk-go/tree/master/pkg/api).
The API for our list command should go in the [`cloud-sdk-go/pkg/api/deploymentapi/list.go`](https://github.com/elastic/cloud-sdk-go/tree/master/pkg/api/deploymentapi/list.go) file
which you'll need create if it does not already exist:

```go
func List(params ListParams) (*models.DeploymentsListResponse, error) {
    // ...
}
```

Things to note about the `List` function:

- This is where the main business logic happens: it should include a call to the Elastic Cloud API to retrieve all the
deployments and return them to the calling function.
- Make sure to properly catch and handle all possible validation errors. Use the [`multierror`](https://github.com/elastic/cloud-sdk-go/blob/master/pkg/api/deploymentapi/get.go#L46-L57)
package even if you're returning a single error. This provides consistency and good UX since the errors will be be properly prefixed.
- For a general set of guidelines on the project's code style, see our [Style Guide](https://github.com/elastic/ecctl/blob/master/developer_docs/STYLEGUIDE.md).

You'll also need to define the `ListParams` struct, as we normally use structs to pass several arguments to API functions.

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
method for `ListParams` need to be tested. These tests should go to `cloud-sdk-go/pkg/api/deploymentapi/list_test.go`
There are many examples of this in the code base, so feel free to browse and use them as inspiration.

That concludes all the steps necessary for creating a new command. You can easily manually test your new command by
importing your local cloud-sdk-go changes running `make fake-sdk`, and making use of our [helper scripts](../CONTRIBUTING.md#helpers):

`dev-cli --config config deployment list`

If your command behaves as expected, all that's left is to make sure you followed all the
[code contribution guidelines](../CONTRIBUTING.md#code-contribution-guidelines) before submitting your PR.

Congratulations on writing your first ecctl command!
