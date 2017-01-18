package tasks

import (
	"github.com/benjamincaldwell/devctl/printer"
	"github.com/benjamincaldwell/devctl/shell"
)

func processCmd(cmd string) error {
	if DryRun {
		printer.InfoBar(cmd)
		return nil
	}
	command := shell.Command("sh", "-c", cmd).SetDir(SrcDir)
	if Verbose {
		return command.PrintOutput()
	}
	return command.Run()
}
