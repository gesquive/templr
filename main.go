package main

import (
	"fmt"

	"github.com/gesquive/shield/cmd"
)

var version = "v0.7.2-git"
var dirty = ""

func main() {
	displayVersion := fmt.Sprintf("shield %s%s",
		version,
		dirty)
	cmd.Execute(displayVersion)
}
