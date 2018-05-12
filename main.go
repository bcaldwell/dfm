package main

import "github.com/bcaldwell/dfm/dfm"

var Version = "dev"
var BuildDate = "n/a"

func main() {
	dfm.Version = Version
	dfm.BuildDate = BuildDate
	dfm.Execute()
}
