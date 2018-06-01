package main

import (
	"fmt"

	"github.com/gesquive/templr/cmd"
)

var version = "v1.0.0-git"
var dirty = ""

func main() {
	displayVersion := fmt.Sprintf("templr %s%s",
		version,
		dirty)
	cmd.Execute(displayVersion)
}
