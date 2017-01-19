package dfm

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"

	"github.com/benjamincaldwell/dfm/utilities"
	"github.com/benjamincaldwell/go-printer"
	"github.com/benjamincaldwell/go-sh"
	"github.com/urfave/cli"
)

var dryRun bool
var verbose bool
var force bool
var overwrite bool

var Version string
var BuildDate string

// set bootstrap env on clone???

// Execute kicks of dfm command
func Execute() {

	var configFile string

	app := cli.NewApp()
	app.Name = "dfm"
	app.Usage = "an easy way to manage dotfiles"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "verbose",
			Usage:       "verbose output",
			Destination: &verbose,
		},
		cli.BoolFlag{
			Name:        "dryrun",
			Usage:       "prints changes",
			Destination: &dryRun,
		},
		cli.BoolFlag{
			Name:  "force, f",
			Usage: "force",
		},
		cli.BoolFlag{
			Name:  "overwrite",
			Usage: "overwrite existing files or folders when linking",
		},
		cli.StringFlag{
			Name:        "config, c",
			Usage:       "sets location of dfm config file or url with config file. Defaults to ~/.dfm.yml.",
			Destination: &configFile,
		},
	}

	printer.Verbose = verbose
	sh.DryRun = dryRun

	printFlagOptions()
	homeDir := detectHomeDir()

	configFile = detectConfigFile(configFile, homeDir)

	config, err := parseConfig(configFile)
	utilities.FatalErrorCheck(err, "Unable to parse configuration file: %s")

	config.SetDefaults(homeDir)

	if _, err := os.Stat(config.SrcDir); os.IsNotExist(err) {
		err = cloneRepo(config.Repo, config.SrcDir)
		utilities.FatalErrorCheck(err, "Unable to clone desired repo")
	}

	app.Version = Version
	app.Author = "Benjamin Caldwell"

	app.Commands = []cli.Command{
		{
			Name:      "install",
			ShortName: "i",
			Usage:     "process each tasks and excuses them",
			Action: func(c *cli.Context) {
				installAction(c.Args(), config)
			},
		},
		{
			Name:      "update",
			ShortName: "u",
			Usage:     "process each tasks and excuses them",
			Action: func(c *cli.Context) {

			},
		},
		{
			Name:      "upgrade",
			ShortName: "up",
			Usage:     "process each tasks and excuses them",
			Action: func(c *cli.Context) {

			},
		},
		{
			Name:      "compile",
			ShortName: "c",
			Usage:     "process each tasks and excuses them",
			Action: func(c *cli.Context) {

			},
		},
		{
			Name:      "git",
			ShortName: "up",
			Usage:     "process each tasks and excuses them",
			Action: func(c *cli.Context) {

			},
		},
		{
			Name:      "path",
			ShortName: "up",
			Usage:     "process each tasks and excuses them",
			Action: func(c *cli.Context) {

			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		printer.Fail("Unexpected failure:", err)
		os.Exit(1)
	}

	createDfmrc(homeDir, configFile, config.SrcDir)
}

func detectHomeDir() (homeDir string) {
	homeDir = os.Getenv("HOME")
	if homeDir == "" {
		printer.Fail("unable to determine user's homeDir")
		os.Exit(1)
	}
	return
}

func printFlagOptions() {
	printer.VerboseInfoBar("Running in verbose mode")
	if dryRun {
		printer.VerboseInfoBar("Running in dryrun mode")
	}
	if force {
		printer.VerboseInfoBar("Running in force mode")
	}
	if overwrite {
		printer.VerboseInfoBar("Running in overwrite mode")
	}
}

func detectConfigFile(configFileFlag, homeDir string) (configFile string) {
	var configURL string
	rcFile := determineRcFile(homeDir)

	if configFileFlag == "" {
		if _, err := os.Stat(rcFile); err == nil {
			dat, err := ioutil.ReadFile(rcFile)
			utilities.FatalErrorCheck(err, "Couldn't read .dfmrc file")
			configFile = string(dat)
		} else {
			configFile = detectDefaultConfigFileLocation()
			printer.VerboseWarning("config file not specified. Defaulting to %s", configFile)
		}
	} else {
		if _, err := os.Stat(configFileFlag); os.IsNotExist(err) {
			r, _ := regexp.Compile(`(?i)^(http|https):\/\/[a-z0-9]+([\-\.]{1}[a-z0-9]+)*\.[a-z]{2,5}(([0-9]{1,5})?\/.*)?$`)
			if r.MatchString(configFileFlag) {
				configURL = configFileFlag
				printer.VerboseInfoBar("Fetching config from %s", configURL)
				file, err := ioutil.TempFile(os.TempDir(), "dfm-")
				utilities.ErrorCheck(err, "creating temp file")
				defer os.Remove(file.Name())
				configFile = file.Name()

				err = utilities.DownloadToFile(configURL, configFile)
				utilities.ErrorCheck(err, "downloading configuration file")
			}
		}
	}

	printer.VerboseInfoBar("Using configuration file located at %s", configFile)
	return
}

func createDfmrc(homeDir, configFile, scrDir string) {
	// offer to create dfmrc file at the end if it doesnt exist
	rcFile := determineRcFile(homeDir)
	if _, err := os.Stat(rcFile); os.IsNotExist(err) {
		if path.Dir(scrDir) != os.TempDir() {
			configPrediction := path.Join(scrDir, "dfm.yml")
			if _, err := os.Stat(configPrediction); err == nil {
				configFile = configPrediction
			}
		}
		if path.Base(scrDir) == ".dotfiles" {
			return
		}
		fmt.Println()
		if utilities.AskForConfirmation("set default configuration to " + configFile + "?") {
			err := ioutil.WriteFile(rcFile, []byte(configFile), 0644)
			utilities.ErrorCheck(err, "writing ~/.dfmrc file")
		}
	}
}

func determineRcFile(homeDir string) string {
	return path.Join(homeDir, ".dfmrc")
}

func cloneRepo(repo, srcDir string) error {
	printer.VerboseInfoBar("cloning %s to %s", repo, srcDir)
	output, err := sh.Command("git", "clone", repo, srcDir).Output()
	if err != nil {
		printer.InfoBar(string(output))
		printer.Fail("Failed to clone repo %s with %s", repo, err)
	}
	return err
}

func detectDefaultConfigFileLocation() string {
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
