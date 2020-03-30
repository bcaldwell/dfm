package dfm

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"

	"github.com/bcaldwell/dfm/utilities"
	"github.com/bcaldwell/go-printer"
	"github.com/bcaldwell/go-sh"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var (
	// global flag
	dryRun        bool
	verbose       bool
	force         bool
	overwrite     bool
	noDfmrcCreate bool
	destDir       string

	showDate bool

	// Version represents the current version of the dfm cli --> Set by build flags
	Version string
	// BuildDate is the date the current version was built --> Set by build flags
	BuildDate string

	// ErrNoHomeEnv is the error that is returned when HOME environment variable is not set
	ErrNoHomeEnv = errors.New("HOME environment variable not set")
	// ErrNoConfigFile is the error that is returned when the configuration file could not be found
	ErrNoConfigFile = errors.New("Could not find a configuration file")

	// Fs is the afero filesystem used when making file system manipulations
	Fs = afero.NewOsFs()

	defaultConfigFiles = []string{"dfm.yml", "$HOME/.dotfiles/dfm.yml", "$HOME/dotfiles/dfm.yml", "$HOME/dfm.yml", "$HOME/.dfm.yml"}
)

// set bootstrap env on clone???

// Execute kicks of dfm command
func Execute() {
	var configFile string

	var config *Configuration

	var rootCmd = &cobra.Command{
		Use:   "dfm",
		Short: "Another way to manage dotfiles",
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
			err := cmd.Help()
			if err != nil {
				printer.ErrorBar("Failed to show help - %s", err)
			}
		},
	}

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dryrun", false, "prints changes")
	rootCmd.PersistentFlags().BoolVarP(&force, "force", "f", false, "force output")
	rootCmd.PersistentFlags().BoolVarP(&overwrite, "overwrite", "o", false, "overwrite existing files or folders when linking")
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "sets location of dfm config file or url with config file. Defaults to ~/.dfm.yml.")
	rootCmd.PersistentFlags().StringVar(&destDir, "destdir", "", "sets the destination directory. Overwrites destdir in configuration file")

	installCommand := &cobra.Command{
		Use:     "install",
		Aliases: []string{"i"},
		Short:   "process each tasks and excuses them",
		Run: func(cmd *cobra.Command, args []string) {
			config = getConfig(configFile)
			installAction(args, config)
		},
	}

	updateCommand := &cobra.Command{
		Use:     "update",
		Aliases: []string{"u"},
		Short:   "updates dfm source repo using git",
		Run: func(cmd *cobra.Command, args []string) {
			config = getConfig(configFile)
			err := updateAction(args, config)
			utilities.ErrorCheck(err, "fetch updates")
		},
	}

	upgradeCommand := &cobra.Command{
		Use:     "upgrade",
		Aliases: []string{"up"},
		Short:   "update dfm source repo and then processes and excuses tasks",
		Run: func(cmd *cobra.Command, args []string) {
			config = getConfig(configFile)
			err := updateAction(args, config)
			if utilities.ErrorCheck(err, "fetch updates") {
				return
			}
			installAction(args, config)
		},
	}

	gitCommand := &cobra.Command{
		Use:     "git",
		Aliases: []string{"g"},
		Short:   "alias to running git in dfm source directory",
		Run: func(cmd *cobra.Command, args []string) {
			config = getConfig(configFile)
			gitAction(args, config)
		},
	}

	pathCommand := &cobra.Command{
		Use:     "path",
		Aliases: []string{"up"},
		Short:   "process each tasks and excuses them",
		Run: func(cmd *cobra.Command, args []string) {
			config = getConfig(configFile)
			fmt.Print(config.SrcDir)
		},
	}

	configCommand := &cobra.Command{
		Use:   "config",
		Short: "print current configuration file",
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var jsonConfig []byte
			config = getConfig(configFile)

			if len(args) > 0 {
				filteredConfig, ok := config.Tasks[args[0]]
				if !ok {
					printer.ErrorBarIcon("Failed to parse json %s", err)
					return
				}
				printer.InfoBarIcon(" %s configuration", args[0])
				jsonConfig, err = json.MarshalIndent(filteredConfig, "", "    ")
			} else {
				jsonConfig, err = json.MarshalIndent(config, "", "    ")
			}

			if err != nil {
				printer.ErrorBarIcon("Failed to parse json %s", err)
				return
			}
			fmt.Println(string(jsonConfig))
		},
	}

	pragmaCommand := &cobra.Command{
		Use:   "pragma",
		Short: "Applies pragma to all files in src directory",
		Run: func(cmd *cobra.Command, args []string) {
			config = getConfig(configFile)
			err := pragmaAction(args, config)
			utilities.ErrorCheck(err, "pragma")
		},
	}

	configFileCommand := &cobra.Command{
		Use:   "configfile",
		Short: "Prints current configuration file",
		Run: func(cmd *cobra.Command, args []string) {
			config = getConfig(configFile)
			fmt.Println(config.configFile)
		},
	}

	versionCommand := &cobra.Command{
		Use: "version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(Version)
			if showDate {
				fmt.Println("Built: " + BuildDate)
			}
		},
	}
	versionCommand.Flags().BoolVarP(&showDate, "date", "d", false, "Show build date")

	rootCmd.AddCommand(
		installCommand,
		updateCommand,
		upgradeCommand,
		gitCommand,
		pragmaCommand,
		pathCommand,
		configCommand,
		versionCommand,
		configFileCommand,
	)

	if err := rootCmd.Execute(); err != nil && !noDfmrcCreate {
		printer.Fail("Unexpected failure:", err)
		os.Exit(1)
	}

	// if config != nil {
	// 	createDfmrc(config.homeDir, config.configFile, config.SrcDir)
	// }
}

func setenv(config *Configuration) {
	os.Setenv("DFM_SRC_DIR", config.SrcDir)
	os.Setenv("DFM_DEST_DIR", config.DestDir)
	os.Setenv("DFM_REPO", config.Repo)
}

func getConfig(configFile string) (config *Configuration) {
	var err error
	config = &Configuration{}
	printer.Verbose = verbose
	sh.DryRun = dryRun

	printFlagOptions()

	config.configFile = configFile

	config.homeDir, err = detectHomeDir()
	utilities.FatalErrorCheck(err, "determining user's homeDir")

	config.configFile, err = detectConfigFile(config.configFile, config.homeDir)
	utilities.FatalErrorCheck(err, "determining configuration file")

	err = config.Parse(config.configFile)
	utilities.FatalErrorCheck(err, fmt.Sprintf("Unable to parse configuration file: %s", configFile))

	err = config.SetDefaults()
	utilities.FatalErrorCheck(err, fmt.Sprintf("Unable to set defaults on configuration: %s", err))

	if destDir != "" {
		config.DestDir = destDir
	}

	if _, err := Fs.Stat(config.SrcDir); os.IsNotExist(err) {
		err = cloneRepo(config.Repo, config.SrcDir)
		utilities.FatalErrorCheck(err, "Unable to clone desired repo")
	}

	setenv(config)

	return config
}

func detectHomeDir() (homeDir string, err error) {
	homeDir = os.Getenv("HOME")
	if homeDir == "" {
		err = ErrNoHomeEnv
	}

	return homeDir, err
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

func detectConfigFile(configFileFlag, homeDir string) (configFile string, err error) {
	var configURL string

	rcFile := determineRcFile(homeDir)

	if configFileFlag == "" {
		if _, err := Fs.Stat(rcFile); err == nil {
			dat, err := ioutil.ReadFile(rcFile)
			utilities.FatalErrorCheck(err, "Couldn't read .dfmrc file")

			configFile = string(dat)
		} else {
			configFile, err = detectDefaultConfigFileLocation()
			printer.VerboseWarning("config file not specified. Defaulting to %s", configFile)
		}
	} else {
		if _, err := Fs.Stat(configFileFlag); os.IsNotExist(err) {
			r, _ := regexp.Compile(`(?i)^(http|https):\/\/[a-z0-9]+([\-\.]{1}[a-z0-9]+)*\.[a-z]{2,5}(([0-9]{1,5})?\/.*)?$`)
			if r.MatchString(configFileFlag) {
				configURL = configFileFlag
				printer.VerboseInfoBar("Fetching config from %s", configURL)
				file, err := afero.TempFile(Fs, "", "dfm-")
				utilities.ErrorCheck(err, "creating temp file")
				defer Fs.Remove(file.Name())
				configFile = file.Name()

				err = utilities.DownloadToFile(configURL, configFile)
				utilities.ErrorCheck(err, "downloading configuration file")
			}
		} else {
			configFile = configFileFlag
		}
	}

	printer.VerboseInfoBar("Using configuration file located at %s", configFile)

	return configFile, err
}

func createDfmrc(homeDir, configFile, scrDir string) {
	// offer to create dfmrc file at the end if it doesnt exist
	rcFile := determineRcFile(homeDir)
	if _, err := Fs.Stat(rcFile); os.IsNotExist(err) {
		tempDir := afero.GetTempDir(Fs, "")
		if path.Dir(scrDir) != tempDir {
			configPrediction := path.Join(scrDir, "dfm.yml")
			if _, err := Fs.Stat(configPrediction); err == nil {
				configFile = configPrediction
			}
		}

		for _, file := range defaultConfigFiles {
			file = os.ExpandEnv(file)
			if path.Clean(path.Dir(scrDir)) == path.Clean(file) {
				return
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

func detectDefaultConfigFileLocation() (string, error) {
	for _, file := range defaultConfigFiles {
		file = os.ExpandEnv(file)

		if _, err := Fs.Stat(file); err == nil {
			return file, nil
		}
	}

	return "", ErrNoConfigFile
}
