package tasks

import (
	"github.com/bcaldwell/go-printer"
	"github.com/bcaldwell/go-sh"
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
