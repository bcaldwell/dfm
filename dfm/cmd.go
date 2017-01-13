package dfm

import (
	"github.com/benjamincaldwell/devctl/printer"
	"github.com/benjamincaldwell/devctl/shell"
)

func processCmd(cmd string, config *Configuration) error {
	if *dryRun {
		printer.InfoBar(cmd)
		return nil
	}
	command := shell.Command("sh", "-c", cmd).SetDir(config.SrcDir)
	if *verbose {
		return command.PrintOutput()
	}
	return command.Run()
}
