package dfm

import "github.com/benjamincaldwell/devctl/shell"

func processGit(args []string, config *Configuration) {
	err := shell.Command("git", args...).SetDir(config.SrcDir).PrintOutput()
	errorCheck(err, "")
}
