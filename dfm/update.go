package dfm

import (
	"github.com/benjamincaldwell/devctl/printer"
	"github.com/benjamincaldwell/devctl/shell"
)

func updateAction(args []string, config *Configuration) bool {
	output, err := shell.Command("git", "fetch").SetDir(config.SrcDir).Output()
	if err != nil {
		printer.Fail("%s failed with %s", "Failed to fetch updates", err)
		printer.InfoBar(string(output))
		return false
	}
	output, err = shell.Command("git", "pull").SetDir(config.SrcDir).Output()
	if err != nil {
		printer.Fail("%s failed with %s", "Failed to fetch updates", err)
		printer.InfoBar(string(output))
		return false
	}
	printer.Success("Pulled latest version from git")
	return true
}
