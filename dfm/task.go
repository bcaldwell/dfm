package dfm

import (
	"runtime"

	"strings"

	"github.com/benjamincaldwell/devctl/printer"
	"github.com/benjamincaldwell/devctl/shell"
)

type Task struct {
	When struct {
		OS           string
		Condition    string
		Parameter    string
		NotInstalled string `yaml:"notInstalled"`
	}
	Deps  []string
	Cmd   []string
	Links []string
	Env   []string
	// 0 not enabled, 2 enabled, 1 can be enabled if dependent on
	Importance byte
}

// look at multiple
func (t Task) calculateImportance(parameter string) byte {
	if strings.ToLower(runtime.GOOS) == strings.ToLower(t.When.OS) {
		return 2
	} else if parameter != "" && strings.ToLower(parameter) == strings.ToLower(t.When.Parameter) {
		return 2
	} else if t.When.Condition != "" {
		shell.DryRun = false
		if err := shell.Command("sh", "-c", t.When.Condition).Run(); err == nil {
			return 2
		}
		shell.DryRun = *dryRun
	} else if t.When.NotInstalled != "" {
		shell.DryRun = false
		if err := shell.Command("sh", "-c", "command -v "+t.When.NotInstalled).Run(); err != nil {
			return 2
		}
		shell.DryRun = *dryRun
	} else if t.When.Condition == "" && t.When.OS == "" && t.When.Parameter == "" && t.When.NotInstalled == "" {
		return 1
	}
	return 0
}

func getTaskList(parameter string, config *Configuration) (tasks []string) {
	if parameter != "" {
		if task, ok := config.Tasks[parameter]; ok {
			tasks = append(tasks, task.appendTaskDependencyList(tasks, parameter, config)...)
			tasks = append(tasks, parameter)
		}
	} else {
		for i, task := range config.Tasks {
			if task.calculateImportance(parameter) == 2 {
				tasks = append(tasks, task.appendTaskDependencyList(tasks, parameter, config)...)
				tasks = append(tasks, i)
			}
		}
	}
	if len(tasks) == 0 {
		printer.Error("no tasks to run")
	}

	// taskNames := reflect.ValueOf(config.Tasks).MapKeys()
	return tasks
}

// check for circular
func (t Task) appendTaskDependencyList(dependencies []string, parameter string, config *Configuration) []string {
	if len(t.Deps) == 0 {
		return []string{}
	}
	for _, depString := range t.Deps {
		if dep, ok := config.Tasks[depString]; ok {
			// fmt.Printf("%s %+v %b\n", depString, dependencies, stringInSlice(depString, dependencies))
			if dep.calculateImportance(parameter) > 0 && !stringInSlice(depString, dependencies) {
				dependencies = append(dependencies, dep.appendTaskDependencyList(dependencies, parameter, config)...)
				dependencies = append(dependencies, depString)
			}
		} else {
			printer.Warning("could find task %s", depString)
		}
	}
	return dependencies
}

func (t Task) execute(config *Configuration) error {

	if len(t.Links) > 0 {
		printer.Info("Running commands")
	}

	for _, command := range t.Cmd {
		err := processCmd(command, config)
		if err != nil {
			printer.Error("%s", err)
		}
	}

	if len(t.Links) > 0 {
		printer.Info("Linking")
	}

	for _, link := range t.Links {
		err := processLink(link, config)
		if err != nil {
			printer.Error("%s", err)
		}
	}

	return nil
}
