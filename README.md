# boa - all snake, no venom
[![Documentation](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/mod/github.com/oncilla/boa)
[![Go](https://img.shields.io/github/workflow/status/oncilla/boa/Go)](https://github.com/Oncilla/boa/actions?query=workflow%3AGo)
[![Go Report Card](https://goreportcard.com/badge/github.com/oncilla/boa)](https://goreportcard.com/report/github.com/oncilla/boa)
[![GitHub issues](https://img.shields.io/github/issues/oncilla/boa/help%20wanted.svg?label=help%20wanted&color=purple)](https://github.com/oncilla/boa/issues?q=is%3Aopen+is%3Aissue+label%3A%22help+wanted%22)
[![GitHub issues](https://img.shields.io/github/issues/oncilla/boa/good%20first%20issue.svg?label=good%20first%20issue&color=purple)](https://github.com/oncilla/boa/issues?q=is%3Aopen+is%3Aissue+label%3A%22good+first+issue%22)
[![Release](https://img.shields.io/github/release-pre/oncilla/boa.svg)](https://github.com/oncilla/boa/releases)
[![License](https://img.shields.io/github/license/oncilla/boa.svg?maxAge=2592000)](https://github.com/oncilla/boa/blob/master/LICENSE)


A venom-free generator for cobra applications.

## Installation

```txt
go get -u github.com/oncilla/boa/cmd/boa
```

## Docker

If you prefer using docker instead:

```txt
WORKDIR=$(pwd)/path/to/application
docker run \
    -v $WORKDIR/:/workdir \
    --user "$(id -u):$(id -g)" \
    docker.pkg.github.com/oncilla/boa/boa:latest
```

## Getting started

`boa` proposes a different layout than what `cobra` proposes as a [typical
structure](https://github.com/spf13/cobra#getting-started).

The commands all live in main package, alongside the main function.

```
cmd/my-app/
├── my-app.go      # main function
├── completion.go  # completion command
└── version.go     # version command
```

The commands should be fairly slim and only take care of checking flags and
initializing the application. Business logic should not go here, but move to
a separate package to ensure it is decoupled from the CLI, reusable and has
no access to global values.

Where you put it is up to you, just try to have the main package as slim as
possible.

### Initializing an application

With boa, your project can have one or multiple cobra based applications.
The simplest approach is to have the `main` package in the project root.

```txt
boa init my-app
```

This creates a runnable cobra application skeleton. Try it out with:

```txt
go run *.go --help
```

The `completion` and `version` are created by default, to showcase how a boa
style application can look like.

If you want to support multiple cobra applications in your project, or want to
keep the root directory clean, you can create the application in a
sub-directory:

```txt
boa init --path cmd/my-app my-app
```

### Adding a command

To add a command, simply run the `add` command:

```txt
boa add greet --flags name:string,age:int
```

If you have provided a path before, make sure to provide the same path again,
or navigate to the corresponding main package first.

This command creates new file called `greet.go` with a license header and
a cobra command generation function.

```go
func newGreet(pather CommandPather) *cobra.Command {
        var flags struct {
                name string
                age  int
        }

        var cmd = &cobra.Command{
                Use:     "greet <arg>",
                Short:   "greet does amazing work!",
                Example: fmt.Sprintf("  %[1]s greet --sample", pather.CommandPath()),
                RunE: func(cmd *cobra.Command, args []string) error {
                        // Add basic sanity checks, where the usage help message should be
                        // printed on error, before this line. After this line, the usage
                        // message is no longer printed on error.
                        cmd.SilenceUsage = true

                        // TODO: Amazing work goes here!
                        return nil
                },
        }

        cmd.Flags().StringVar(&flags.name, "name", "", "name description")
        cmd.Flags().IntVar(&flags.age, "age", 0, "age description")
        return cmd
}
```

Boa suggests to use the `SilenceErrors` and `SilenceUsage`.
For more information, see: https://github.com/spf13/cobra/issues/340#issuecomment-374617413

You now need to register the command with its parent. For the sake of this
example, it is simply the root command. Update `my-app.go` with:

```go
    cmd.AddCommand(
        newCompletion(cmd),
        newGreet(cmd),
        newVersion(cmd),
    )

```

That's it, the new command is now registered and can already be used:

```txt
$ go run *.go greet --help
greet does amazing work!

Usage:
  my-app greet <arg> [flags]

Examples:
  my-app greet --sample

Flags:
      --age int       age description
  -h, --help          help for greet
      --name string   name description
```

## Why boa?

The [cobra](https://github.com/spf13/cobra) library is an amazing and powerful
toolkit for creating command line applications. However, the example projects
display some drawbacks that boa tries to improve upon.

### Directory structure

The proposed rigid directory structure in the examples and the generator go
against the commonly used patterns how applications are structured nowadays.

boa proposes to have all commands in the same directory. The key here is, that
the `main` package should only be used for initialization and very simple tasks.
Bigger business logic should live in a separate package, ensuring that it is
reusable, and decoupled from the CLI interface.

### Command registration

Commands are proposed to be registered inside an `init` function. Essentially
forcing global state, and making it hard for commands to be tested.

With the approach proposed by boa, testing a command is as simple as:

```go
func TestMyCommand(t *testing.T) {
    cmd := newMyCommand(boa.Pather(""))
    cmd.SetArgs([]string{"--my", "args"})
    err := cmd.Execute()
    if err != nil {
        t.Fail()
    }
}
```

### Global flags

Flags are proposed to be package global variables and registered in the `init`
function. This can lead to code that is full of surprises, as package globals
taint logic very easily if you do not take care.

With the approach proposed by boa, each instance of a command has its own set
of flags, and there is no way for other components to access them directly.
