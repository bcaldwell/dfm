package tasks

import (
	"github.com/benjamincaldwell/go-printer"
	"github.com/benjamincaldwell/go-sh"
)

func processCmd(cmd string) error {
	if DryRun {
		printer.InfoBar(cmd)
		return nil
	}
	command := sh.Command("sh", "-c", cmd).SetDir(SrcDir)
	if Verbose {
		return command.PrintOutput()
	}
	return command.Run()
}
