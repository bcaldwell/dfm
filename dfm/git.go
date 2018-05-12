package dfm

import (
	"github.com/bcaldwell/dfm/utilities"
	"github.com/bcaldwell/go-sh"
)

func gitAction(args []string, config *Configuration) {
	err := sh.Command("git", args...).SetDir(config.SrcDir).PrintOutput()
	utilities.ErrorCheck(err, "")
}
