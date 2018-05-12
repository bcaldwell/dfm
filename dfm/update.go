package dfm

import (
	"github.com/bcaldwell/go-printer"
	"github.com/bcaldwell/go-sh"
)

func updateAction(args []string, config *Configuration) (err error) {
	output, err := sh.Command("git", "fetch").SetDir(config.SrcDir).Output()
	if err != nil {
		printer.VerboseInfoBar(string(output))
		return
	}
	output, err = sh.Command("git", "pull").SetDir(config.SrcDir).Output()
	if err != nil {
		printer.VerboseInfoBar(string(output))
		return
	}
	printer.Success("Pulled latest version from git")
	return
}
