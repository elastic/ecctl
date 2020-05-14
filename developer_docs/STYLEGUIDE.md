# Style Guide

- [General patterns](#general-patterns)
- [Errors](#errors)
  - [Error handling](#error-handling)
  - [API Errors](#api-errors)
  - [Multiple errors](#multiple-errors)
- [Packages](#packages)
  - [Structure](#structure)
  - [Util packages](#util-packages)
- [Testing](#testing)
  - [Unit tests](#unit-tests)
  - [Testing commands](#testing-commands)
- [General Style](#general-style)
  - [Usage of pointers](#usage-of-pointers)
  - [Arrays and slices](#arrays-and-slices)
  - [Config or params as the only function parameter](#config-or-params-as-the-only-function-parameter)
  - [Functions and methods](#functions-and-methods)
- [Documentation](#documentation)
- [Extras](#extras)

## General patterns

Our codebase is structured very similarly to what [DDD (Domain Driven Design)](https://martinfowler.com/bliki/EvansClassification.html) dictates, and even though there’s no strict naming or patterns being enforced we try to follow it as a good practice.

## Errors

All errors should be informative and helpful. Something like "a cluster ID is required for this command" is more useful than "no ID". This is a helpful blog post on the subject: [How to Write Good Error Messages.](https://uxplanet.org/how-to-write-good-error-messages-858e4551cd4)

Error strings should not be capitalised or end with punctuation, since they are usually printed following other context.

In regards to packages, we prefer to use the `"errors"` standard library package over `"github.com/pkg/errors"`

### Error handling

All errors must be handled and returned, `_` variables must not be used to discard them.

#### Exceptions

Parsing command line arguments/flags is done like this

```go
threshold, _ := cmd.Flags().GetUint32(thresholdArg)
```

The above line would return an error only if the flag is not defined, or the datatype does not match the flag declaration data type.
Adding many `if err` checks will make the code a little bit noisy so we can ignore the errors in these cases.

### API Errors

API errors should always be encapsulated with `apierr.Unwrap()`, this function tries to break down and inspect the encapsulated and multi-layer wraps that the API errors contain.

### Multiple errors

When multiple errors can be returned, it is preferable to use the `mutlierror.Prefixed` type to return all the possible errors with a prefixed string to include some context.

yes! :smile:

```go
func (params StopParams) Validate() error {
    var merr = multierror.NewPrefixed("stop operation")
    if params.ID == "" {
        merr = merr.Append(errors.New("ID cannot be empty"))
    }

    if params.API == nil {
        merr = merr.Append(errors.New("api reference is required"))
    }

   return merr.ErrorOrNil()
}
```

preferably not :confused:

```go
func (params StopParams) Validate() error {
    if params.ID == "" {
        return errors.New("ID cannot be empty")
    }

    if params.API == nil {
        return errors.New("api reference is required")
    }

   return nil
}
```

## Packages

### Structure

A package per context, a context applies to any of the high level containers like `platform`, `deployment`, etc. When a context becomes too large to be contained in a package we can start breaking it down into sub-packages.

An example would be  [pkg/deployment/elasticsearch/plan](../pkg/deployment/elasticsearch/plan).

### Util packages

When a function can be made generic it should go in one of the utils packages (e.g. [pkg/util](../pkg/util), [pkg/util](../pkg/util)) to remove complexity and give the ability to be reused.

If the function is not specific to `ecctl`, it should be part of [cloud-sdk-go](https://github.com/elastic/cloud-sdk-go) or a standalone repository if the functionality is big enough.

## Testing

### Unit tests

All files containing functions or methods must have a corresponding unit test file, and we aim to have 100% coverage.

#### API Mocks

When unit testing functions which will call the external API, please use the provided `api.NewMock` in conjunction with `mock.Response`.

yes! :smile:

``` go
import (
    "net/http"

    "github.com/elastic/cloud-sdk-go/pkg/api"
    "github.com/elastic/cloud-sdk-go/pkg/api/mock"
)

//Test case
{
    name: "succeeds",
    args: args{params{
        API: api.NewMock(mock.Response{
            Response: http.Response{
                Body:       mock.NewStringBody(`{}`),
                StatusCode: 200,
            },
        }),
    }},
},
// More tests ...
```

### Testing commands

When writing unit tests for commands, we only look to assert that the command is constructing the correct API call. API responses are mocked and tested only in the `pkg/` directory.

See [TestRunShowKibanaClusterCmd()](./cmd/kibana/show_test.go) as a good example to base your tests on.

## General Style

Before committing to your feature branch, run `make lint` and `make format` to ensure the code has the proper styling format. We run linters on Jenkins, so this also prevents your builds from failing.

### Usage of pointers

Unless a pointer is necessary, always use the value type before switching to the pointer type. Of course, if your structure contains a mutex, or needs to be synced between different goroutines, it needs to be a pointer, otherwise there’s no reason why it should be.

### Arrays and slices

To remove complexity, when it's not necessary to set length or capacity, it's preferable to declare a nil slice as it has no underlying array. We should avoid using the `make` function, as this function allocates a zeroed array that returns a slice of said array.

yes! :smile:

```go
var slice []string
```

preferably not :confused:

```go
slice := make([]string, 0)
```

### `Config` or `Params` as the only function parameter

A `Config` or `Params` structure is used to encapsulate all the parameters in a structure that has a `.Validate()` error signature, so it can be validated inside the receiver of that structure.

Unless the `Params` struct needs to satisfy the `pool.Validator` interface for concurrency on that receiver it should always remain a value type and not a pointer type.

### Functions and methods

Names should be descriptive and we should avoid redundancy.

yes! :smile:

```go
kibana.Create()
```

preferably not :confused:

```go
kibana.CreateKibanaDeployment()
```

When using method chaining make sure to put each method in it's own line to improve readability.

yes! :smile:

```go
res, err := a.API.V1API.ClustersApm.GetApmClusterPlanActivity(
clusters_apm.NewGetApmClusterPlanActivityParams().
     WithClusterID(params.id).
     WithShowPlanDefaults(params.defaults),
)
```

preferably not :confused:

```go
res, err := a.API.V1API.ClustersApm.GetApmClusterPlanActivity(
    clusters_apm.NewGetApmClusterPlanActivityParams().WithClusterID(params.id).WithShowPlanDefaults(params.defaults),
)
```

When possible we try to avoid `else` and nested `if`s. This makes our code more readable and removes complexity.
Specifically, we prefer having multiple `return` statements over having nested code.

yes! :smile:

``` go
if params.Hide {
    params.Do()
    return data, nil
}

if isHidden {
    return nil, fmt.Errorf("example error", params.Name)
}
return data, nil
```

preferably not :confused:

``` go
if params.Hide {
        params.Do()
    } else {
        if isHidden {
            return nil, fmt.Errorf("example error", params.Name)
        }
    }
```

## Documentation

We use `make docs` to automatically generate documentation for our commands which live in the `cmd` folder.

It is important when writing the descriptions to our commands or flags, that we use simple language and are as clear as possible to provide good UX. If you need to explain more about the command or give examples, please do so using the `Example` field, a good example is the [deployment elasticsearch list](cmd/deployment/elasticsearch/list.go) command.

The package wide description and documentation is provided in a godoc `doc.go` file. Aside form packages with a very small context, all packages should have this file.

## Extras

For further information on good practices with Go in general, check out this [document](https://github.com/golang/go/wiki/CodeReviewComments).
