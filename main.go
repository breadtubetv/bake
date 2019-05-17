package main

import (
	"github.com/breadtubetv/bake/cmd"
)

var (
	// Populated by goreleaser during build
	version = "master"
	commit  = "?"
	date    = ""
)

func main() {
	cmd.Execute(version, commit, date)
}
