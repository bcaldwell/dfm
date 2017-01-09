package dfm

import (
	"github.com/benjamincaldwell/devctl/printer"
	"github.com/benjamincaldwell/devctl/shell"
)

func processCmd(cmd string) error {
	if *dryRun {
		printer.InfoBar(cmd)
		return nil
	}
	command := shell.Command("sh", "-c", cmd)
	if *verbose {
		return command.PrintOutput()
	}
	return command.Run()
}
