package dfm

import (
	"flag"
	"fmt"

	"github.com/benjamincaldwell/devctl/printer"
)

func Execute() {
	verbose := flag.Bool("verbose", false, "verbose output")
	configFile := flag.String("config", "~/dfm.yml", "Sets location of dfm config file or url with config file. Defaults to ~/.dfm.yml.")

	*configFile = "dfm.example.yml"

	flag.Parse()

	printer.Verbose = *verbose

	// printer.Info(*configFile)
	// printer.Info("%b", *verbose)

	config, err := parseConfig(*configFile)
	fmt.Println(err)
	// fmt.Printf("%+v\n", config)

	args := flag.Args()
	var parameter string
	if len(args) > 0 {
		parameter = args[0]
	}

	a := getTaskList(parameter, config)
	fmt.Println(a)
	fmt.Println(uniqueSliceTransform((a)))

	// for i, task := range config.Tasks {
	// 	fmt.Printf("\n\n%s\n", i)
	// 	fmt.Println(task.calculateImportance(parameter))
	// }

	// for _, task := range config.Tasks {
	// 	fmt.Println(task)
	// }

	// fmt.Printf("%+v\n", config.Tasks)

}
