package dfm

import (
	"strings"

	"github.com/benjamincaldwell/devctl/printer"
)

func processInstall(args []string, config *Configuration) {
	parameter := ""
	if len(args) > 0 {
		parameter = args[0]
	}
	taskList := getTaskList(parameter, config)
	taskList = uniqueSliceTransform((taskList))
	printer.VerboseInfoBar("Running tasks: %s", strings.Join(taskList, ","))

	for _, task := range taskList {
		printer.Info("Executing %s\n", task)
		config.Tasks[task].execute(config)
	}
}
