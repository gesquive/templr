package main

import (
	"fmt"

	"github.com/gesquive/shield/cmd"
)

var version = "v0.6.1-git"
var dirty = ""

func main() {
	displayVersion := fmt.Sprintf("shield %s%s",
		version,
		dirty)
	cmd.Execute(displayVersion)
}
