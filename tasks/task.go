package tasks

import (
	"runtime"

	"strings"

	"github.com/bcaldwell/dfm/utilities"
	"github.com/bcaldwell/go-printer"
	"github.com/bcaldwell/go-sh"
)

type Task struct {
	When struct {
		OS           string `yaml:"os"`
		Condition    string `yaml:"condition"`
		Parameter    string `yaml:"parameter"`
		NotInstalled string `yaml:"notInstalled"`
	}
	Deps     []string
	Cmd      []string
	Links    []string
	Env      []string
	Template Template
	// 0 not enabled, 2 enabled, 1 can be enabled if dependent on
	// importance byte
}

var (
	SrcDir    string
	DestDir   string
	Verbose   bool
	DryRun    bool
	Force     bool
	Overwrite bool
	absPath   = utilities.AbsPath
)

func ExecuteTasks(tasks map[string]Task, task string) error {
	taskList := getTaskList(task, tasks)
	taskList = utilities.UniqueSliceTransform((taskList))
	printer.VerboseInfoBar("Running tasks: %s", strings.Join(taskList, ","))

	printer.Verbose = Verbose
	sh.DryRun = DryRun

	for _, task := range taskList {
		printer.Info("Executing %s\n", task)
		err := tasks[task].Execute(task)

		if err != nil {
			return err
		}
	}

	return nil
}

// look at multiple
func (t Task) calculateImportance(parameter string) byte {
	if strings.EqualFold(runtime.GOOS, t.When.OS) {
		return 2
	} else if parameter != "" && strings.EqualFold(parameter, t.When.Parameter) {
		return 2
	} else if t.When.Condition != "" {
		sh.DryRun = false
		if err := sh.Command("sh", "-c", t.When.Condition).Run(); err == nil {
			return 2
		}
		sh.DryRun = DryRun
	} else if t.When.NotInstalled != "" {
		sh.DryRun = false
		if err := sh.Command("sh", "-c", "command -v "+t.When.NotInstalled).Run(); err != nil {
			return 2
		}
		sh.DryRun = DryRun
	} else if t.When.Condition == "" && t.When.OS == "" && t.When.Parameter == "" && t.When.NotInstalled == "" {
		return 1
	}

	return 0
}

func getTaskList(parameter string, tasks map[string]Task) (tasksList []string) {
	if parameter != "" {
		if task, ok := tasks[parameter]; ok {
			tasksList = append(tasksList, task.appendTaskDependencyList(tasksList, parameter, tasks)...)
			tasksList = append(tasksList, parameter)
		}
	} else {
		for i, task := range tasks {
			if task.calculateImportance(parameter) == 2 {
				tasksList = append(tasksList, task.appendTaskDependencyList(tasksList, parameter, tasks)...)
				tasksList = append(tasksList, i)
			}
		}
	}

	if len(tasksList) == 0 {
		printer.Error("no tasks to run")
	}

	// taskNames := reflect.ValueOf(config.Tasks).MapKeys()
	return tasksList
}

// check for circular
func (t Task) appendTaskDependencyList(dependencies []string, parameter string, tasks map[string]Task) []string {
	if len(t.Deps) == 0 {
		return []string{}
	}

	for _, depString := range t.Deps {
		if dep, ok := tasks[depString]; ok {
			// fmt.Printf("%s %+v %t\n", depString, dependencies, utilities.StringInSlice(depString, dependencies))
			if dep.calculateImportance(parameter) > 0 && !utilities.StringInSlice(depString, dependencies) {
				dependencies = append(dependencies, dep.appendTaskDependencyList(dependencies, parameter, tasks)...)
				dependencies = append(dependencies, depString)
			}
		} else {
			printer.Warning("could find task %s", depString)
		}
	}

	return dependencies
}

func (t Task) Execute(name string) error {
	if len(t.Links) > 0 {
		printer.VerboseInfo("Running commands")
	}

	for _, command := range t.Cmd {
		err := processCmd(command)
		if err != nil {
			printer.Error("%s: %s", name, err)
		}
	}

	if t.Template.isDefined() {
		printer.VerboseInfo("Processing template")

		err := processTemplate(t.Template)
		if err != nil {
			printer.Error("%s: %s", name, err)
		}
	}

	if len(t.Links) > 0 {
		printer.VerboseInfo("Linking")
	}

	for _, link := range t.Links {
		err := processLink(link)
		if err != nil {
			printer.Error("%s: %s", name, err)
		}
	}

	return nil
}
