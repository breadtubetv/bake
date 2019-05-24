# BreadTube Bake CLI

The `bake` cli is a companion CLI to the [breadtubetv/breadtubetv](https://github.com/breadtubetv/breadtubetv) project.

## [Contributing](https://github.com/breadtubetv/bake/blob/master/CONTRIBUTING.md)

Scripts are being written in Go, this keeps the scripting and operational language the same, provides cross system compatibility, and gives everyone an opportunity to learn a new programming language.

## Installation

The compiled bake binary has no dependencies, so you can download your platforms release from the [release page](https://github.com/breadtubetv/bake/releases), pop it in your `$PATH` and off you go.

If you want to build from source you will need the Go programming language installed, and your `$GOPATH` configured.

```bash
# Enable Go 1.11+ module support
export GO111MODULE=on
go get github.com/breadtubetv/bake
```

This will install the binary into your `$GOPATH/bin` folder. If you don't already it's recommended you add this directory to your `$PATH`.

## Usage

### Configuring the `bake` CLI

The `.bake.yaml` configuration file can be stored in the following locations:

- `$HOME/.bake.yaml`
- `./bake.yaml` (In other words, the current directory from which you're running the CLI)

Current configuration options and default values:

- `projectRoot: "../"` : Directory of channel data files.  
  E.g. `$GOPATH/src/github.com/breadtubetv/breadtubetv/data/channels`

### Commands

Once you have `bake` installed:

#### Import a Channel

```bash
bake channel import creator_slug youtube channel_url
```

#### Import a Video

##### Using the Video ID

```bash
bake import video --creator creator_slug --provider youtube --id VIDEO_ID
```

##### Using the Video URL

```bash
bake import video --creator creator_slug --provider youtube --url https://VIDEO_URL
```

Note: The following formats are supported

- https://www.youtube.com/watch?v=xspEtjnSfQA
- https://www.youtube.com/embed/xspEtjnSfQA

## Contributing

You'll need the Go programming language installed. We recommend version v1.12+. This is going to be dependent on your system, we recommend following <https://golang.org/doc/install>

Clone the repo and pull dependencies:

```bash
# Enable Go 1.11 module support
export GO111MODULE=on
# Clone directory
git clone https://github.com/breadtubetv/bake
cd bake
# Download dependencies
go mod download
```

> **NOTE:** The `spf13/cobra` generator CLI won't work if you don't clone the project into `$GOPATH`. This is not a requirement to develop for the project.

### Linting & Pre-Commit

We recommend using [`golangci/golangci-lint`](https://github.com/golangci/golangci-lint) for linting. A development config (`.golangci.dev.yaml`) is provided. To use this config:

```shell
golangci-lint run --config=.golangci.dev.yml
```

A `.pre-commit-config.yaml` is provided for pre-commit hooks using [`pre-commit/pre-commit`](https://github.com/pre-commit/pre-commit).

### Submitting a PR

We welcome PRs! The only thing we ask is that you ensure you keep the `go.mod` and `go.sum` files clean by running the following:

```shell
GO111MODULE=on go mod tidy
```

### Testing

Bake has some very basic tests for now, they can be run with the standard go test command line:

```bash
go get -t ./...
go test ./...
```

### Releasing

Releasing is automated via `git tag` and CircleCI. Users with write permissions will be able to create tags. To create a new release:

```shell
# vX.Y.Z needs to be a valid SemVer version number
git tag vX.Y.Z
git push --tags
```

CircleCI will do the rest!
