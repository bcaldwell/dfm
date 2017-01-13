package dfm

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"

	"github.com/benjamincaldwell/devctl/printer"
	"github.com/benjamincaldwell/devctl/shell"
)

var dryRun *bool
var verbose *bool
var force *bool
var overwrite *bool

var Version string
var BuildDate string

// Execute kicks of dfm command
func Execute() {
	var config *Configuration

	verbose = flag.Bool("verbose", false, "verbose output")
	dryRun = flag.Bool("dryrun", false, "print changes")
	force = flag.Bool("force", false, "force")
	overwrite = flag.Bool("overwrite", false, "overwrite existing files or folders when linking")
	configFile := flag.String("config", "", "Sets location of dfm config file or url with config file. Defaults to ~/.dfm.yml.")
	configURL := ""

	flag.Parse()

	printer.Verbose = *verbose
	shell.DryRun = *dryRun

	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		printer.Fail("unable to determine user's homeDir")
		os.Exit(1)
	}

	printer.VerboseInfoBar("Running in verbose mode")
	if *dryRun {
		printer.VerboseInfoBar("Running in dryrun mode")
	}
	if *force {
		printer.VerboseInfoBar("Running in force mode")
	}
	if *overwrite {
		printer.VerboseInfoBar("Running in overwrite mode")
	}

	rcFile := path.Join(homeDir, ".dfmrc")
	if *configFile == "" {
		if _, err := os.Stat(rcFile); err == nil {
			dat, err := ioutil.ReadFile(rcFile)
			fatalErrorCheck(err, "Couldn't read .dfmrc file")
			*configFile = string(dat)
		} else {
			*configFile = detectConfigFileLocation(*configFile)
			printer.VerboseWarning("config file not specified. Defaulting to %s", *configFile)
		}
	} else {
		if _, err := os.Stat(*configFile); os.IsNotExist(err) {
			r, _ := regexp.Compile(`(?i)^(http|https):\/\/[a-z0-9]+([\-\.]{1}[a-z0-9]+)*\.[a-z]{2,5}(([0-9]{1,5})?\/.*)?$`)
			if r.MatchString(*configFile) {
				configURL = *configFile
				printer.VerboseInfoBar("Fetching config from %s", configURL)
				file, err := ioutil.TempFile(os.TempDir(), "dfm-")
				errorCheck(err, "creating temp file")
				defer os.Remove(file.Name())
				*configFile = file.Name()

				err = DownloadToFile(configURL, *configFile)
				errorCheck(err, "downloading configuration file")
			}
		}

		// offer to create dfmrc file at the end if it doesnt exist
		if _, err := os.Stat(rcFile); os.IsNotExist(err) {
			defer func() {
				if configURL != "" {
					if config.SrcDir != ".dotfiles" {
						return
					}
					configPrediction := path.Join(config.SrcDir, "dfm.yml")
					if _, err := os.Stat(configPrediction); err == nil {
						*configFile = configPrediction
					}
				}
				fmt.Println()
				if AskForConfirmation("set default configuration to " + *configFile + "?") {
					err := ioutil.WriteFile(rcFile, []byte(*configFile), 0644)
					errorCheck(err, "writing ~/.dfmrc file")
				}
			}()
		}
	}

	printer.VerboseInfoBar("Using configuration file located at %s", *configFile)
	config, err := parseConfig(*configFile)
	fatalErrorCheck(err, "Unable to parse configuration file: %s")

	if config.SrcDir == "" {
		config.SrcDir = path.Join(homeDir, ".dotfiles")
		printer.VerboseWarning("srcDir not specified. Defaulting to %s", config.SrcDir)
	}

	if config.DestDir == "" {
		config.DestDir = homeDir
		printer.VerboseWarning("destDir not specified. Defaulting to %s", config.DestDir)
	}

	if _, err := os.Stat(config.SrcDir); os.IsNotExist(err) {
		printer.VerboseInfoBar("cloning %s to %s", config.Repo, config.SrcDir)
		output, err := shell.Command("git", "clone", config.Repo, config.SrcDir).Output()
		if err != nil {
			printer.InfoBar(string(output))
			printer.Fail("Failed to clone repo %s with %s", config.Repo, err)
		}
	}
	// preclone commands?? (ssh keys)

	args := flag.Args()

	var commandName string
	var arguments []string
	if len(args) == 0 {
		printer.Fail("Atleast 1 argument is required")
		os.Exit(1)
	} else if len(args) > 1 {
		arguments = args[1:]
	}
	commandName = args[0]

	switch commandName {
	case "install":
		processInstall(arguments, config)
	case "update":
		processUpdate(arguments, config)
	case "upgrade":
		if processUpdate(arguments, config) {
			processInstall(arguments, config)
		}
	case "path":
		fmt.Println(config.SrcDir)
	case "git":
		processGit(arguments, config)
	case "version":
		fmt.Println(Version)
	default:
		printer.Error("Unknown subcommand %s", commandName)
	}
}

func detectConfigFileLocation(configFile string) string {
	if configFile != "" {
		return configFile
	}
	defaultConfigFiles := []string{"dfm.yml", "$HOME/.dotfiles/dfm.yml", "$HOME/dfm.yml", "$HOME/.dfm.yml"}
	for _, file := range defaultConfigFiles {
		file = os.ExpandEnv(file)
		if _, err := os.Stat(file); err == nil {
			return file
		}
	}
	printer.Error("could not find a configuration file")
	os.Exit(1)
	return ""
}
