package dfm

import (
	"github.com/benjamincaldwell/dfm/utilities"
	"github.com/benjamincaldwell/go-sh"
)

func gitAction(args []string, config *Configuration) {
	err := sh.Command("git", args...).SetDir(config.SrcDir).PrintOutput()
	utilities.ErrorCheck(err, "")
}
