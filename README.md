# Bake CLI

The `bake` cli is a companion CLI to the [breadtubetv/breadtubetv](https://github.com/breadtubetv/breadtubetv) project. 

## [Contributing](https://github.com/breadtubetv/bake/blob/master/CONTRIBUTING.md)

Scripts are being written in Go, this keeps the scripting and operational language the same, provides cross system compatibility, and gives everyone an opportunity to learn a new programming language.

### Installation

#### Installing Go

This is going to be dependent on your system, we recommend following https://golang.org/doc/install

#### Developing

You'll need a copy of the project in your `$GOPATH`.

```
go get github.com/breadtubetv/bake
cd $GOPATH/github.com/breadtubetv/bake
```

And you'll need to install bake's dependencies with:

```
go get
```

#### Installing

This will put the `bake` command in your path:

```
go install
```

#### Configuring the `bake` CLI

The `.bake.yaml` configuration file can be stored in the following locations:

- `$HOME/.bake.yaml`
- `./bake.yaml` (In other words, the current directory from which you're running the CLI)

Current configuration options and default values:

- `projectRoot: "../"` : Directory of channel data files.    
  E.g. `$GOPATH/src/github.com/breadtubetv/breadtubetv/data/channels`

#### Importing a Channel

You can run bake directly from the source like so:

```
go run main.go config youtube #follow prompts
go run main.go channel import contrapoints youtube https://www.youtube.com/user/contrapoints
```

Or if you've installed it:

```
bake config youtube #follow prompts
bake channel import contrapoints youtube https://www.youtube.com/user/contrapoints
```

#### Testing

Bake has some very basic tests for now, they can be run with the standard go test command line:

```
go get -t ./...
go test ./...
```
