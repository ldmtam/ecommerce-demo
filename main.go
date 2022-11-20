package main

import (
	"github.com/ldmtam/ecommerce-demo/cmd"
)

var version = "0.0.1"

func main() {
	cmd.SetVersion(version)
	cmd.Execute()
}
