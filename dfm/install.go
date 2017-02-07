package dfm

import "github.com/benjamincaldwell/dfm/tasks"

func installAction(args []string, config *Configuration) error {
	parameter := ""
	if len(args) > 0 {
		parameter = args[0]
	}
	tasks.SrcDir = config.SrcDir
	tasks.DestDir = config.DestDir
	tasks.DryRun = dryRun
	tasks.Force = force
	tasks.Overwrite = overwrite
	tasks.Verbose = verbose
	return tasks.ExecuteTasks(config.Tasks, parameter)
}
