package dfm

import (
	"github.com/benjamincaldwell/go-printer"
	"github.com/benjamincaldwell/go-sh"
)

func updateAction(args []string, config *Configuration) bool {
	output, err := sh.Command("git", "fetch").SetDir(config.SrcDir).Output()
	if err != nil {
		printer.Fail("%s failed with %s", "Failed to fetch updates", err)
		printer.InfoBar(string(output))
		return false
	}
	output, err = sh.Command("git", "pull").SetDir(config.SrcDir).Output()
	if err != nil {
		printer.Fail("%s failed with %s", "Failed to fetch updates", err)
		printer.InfoBar(string(output))
		return false
	}
	printer.Success("Pulled latest version from git")
	return true
}
