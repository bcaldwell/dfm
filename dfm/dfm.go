package dfm

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/benjamincaldwell/devctl/printer"
	"github.com/benjamincaldwell/devctl/shell"
)

var dryRun *bool
var verbose *bool
var force *bool
var overwrite *bool

func Execute() {
	verbose = flag.Bool("verbose", false, "verbose output")
	dryRun = flag.Bool("dryrun", false, "print changes")
	force = flag.Bool("force", false, "force")
	overwrite = flag.Bool("overwrite", false, "overwrite existing files or folders when linking")
	configFile := flag.String("config", "~/dfm.yml", "Sets location of dfm config file or url with config file. Defaults to ~/.dfm.yml.")

	*configFile = "dfm.example.yml"

	flag.Parse()

	printer.Verbose = *verbose
	shell.DryRun = *dryRun

	shell.Command("sh", "-c")

	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		printer.Fail("unable to determine user's homeDir")
		os.Exit(1)
	}

	// printer.Info(*configFile)
	printer.Info("%b", *dryRun)

	config, err := parseConfig(*configFile)

	if config.SrcDir == "" {
		config.SrcDir = path.Join(homeDir, ".dotfiles")
		printer.VerboseWarning("srcDir not specified. Defaulting to %s", config.SrcDir)
	}

	if config.DestDir == "" {
		config.DestDir = homeDir
		printer.VerboseWarning("destDir not specified. Defaulting to %s", config.DestDir)
	}

	fmt.Println(err)

	args := flag.Args()
	var parameter string
	if len(args) > 0 {
		parameter = args[0]
	}

	taskList := getTaskList(parameter, config)
	taskList = uniqueSliceTransform((taskList))
	fmt.Println(taskList)

	for _, task := range taskList {
		printer.Info("Executing %s\n", task)
		config.Tasks[task].execute(config)
	}

	// for _, task := range config.Tasks {
	// 	fmt.Println(task)
	// }

	// fmt.Printf("%+v\n", config.Tasks)

}
