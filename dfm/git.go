package dfm

import (
	"github.com/benjamincaldwell/devctl/shell"
	"github.com/benjamincaldwell/dfm/utilities"
)

func gitAction(args []string, config *Configuration) {
	err := shell.Command("git", args...).SetDir(config.SrcDir).PrintOutput()
	utilities.ErrorCheck(err, "")
}
