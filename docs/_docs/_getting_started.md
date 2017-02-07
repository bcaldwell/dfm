# dfm

[![Go Report Card](https://goreportcard.com/badge/github.com/benjamincaldwell/dfm)](https://goreportcard.com/report/github.com/benjamincaldwell/dfm) [![codebeat badge](https://codebeat.co/badges/9d235a39-896d-4dd3-9cf8-bb5f85b8cf66)](https://codebeat.co/projects/github-com-benjamincaldwell-dfm)
[![Build Status](https://travis-ci.org/benjamincaldwell/dfm.svg?branch=master)](https://travis-ci.org/benjamincaldwell/dfm)

dfm is a tool for managing dotfiles. dfm works best when using git to manage dotfile but will also work without.

## Table of Contents

  * [dfm](#dfm)
    * [Table of Contents](#table-of-contents)
    * [Installation](#installation)
      * [Install script](#install-script)
      * [Manual Install](#manual-install)
    * [Yaml configuration file](#yaml-configuration-file)
    * [Usage](#usage)
      * [Commands](#commands)
          * [install](#install)
          * [update](#update)
          * [upgrade](#upgrade)
          * [git](#git)
          * [path](#path)
      * [Flags](#flags)
          * [config](#config)
          * [verbose](#verbose)
          * [dryrun](#dryrun)
          * [overwrite](#overwrite)
          * [force](#force)
      * [Enable cd command](#enable-cd-command)

## Installation

### Install script
```
curl -o- https://raw.githubusercontent.com/benjamincaldwell/dfm/master/scripts/install.sh | bash
```


### Manual Install
Download the applicable binary from [releases](https://github.com/benjamincaldwell/dfm/releases)


## Yaml configuration file

### Example
```
---

# remote location of repository [optional]
repo: git@github.com:benjamincaldwell/dotfiles.git

# Links occur fron srcDir/path to destDir/path
# location of source files. Defaults to $HOME/.dotfiles
srcDir: /Home/user/.dotfiles

# default location to link files relatively to. Defaults to $HOME
destDir: /home/user

# Tasks to run
tasks:
  name:
    # Parameter to run the task
    when:
      # runs tasks when os matches
      os: darwin
      # runs the shell command and runs the task if exit code is 0
      condition: "echo"
      # runs when the text after install matchs Ex: dfm install darwin
      parameter "darwin"
      # Runs if command doesnt exists in the shell
      notInstalled: "brew"

    # commands to run
    cmd:
      - "brew install"
    
    # Links to create
    links:
      - gitconfig:.gitconfig

    # task dependencies
    deps:
      - depend

```

## Usage

### Commands

##### install
Process each tasks and excuses it

```
dfm install
```

##### update
To use git to update the repository by running `git fetch && git pull` run:

```
dfm update
```

##### upgrade
Same as `dfm update` but runs `dfm install` afterwards:

```
dfm upgrade
```

##### git
Runs the passed in git command in the source respository:

```
dfm git args
```

Ex: `dfm git status`

##### path
Returns the path of the source respository. Useful for `cd`

```
cd $(dfm cd)
```

### Flags
##### config
Sets the location of the configuration file

```
dfm -config pathtoconfig/config.yml
``` 

##### verbose
Enables verbose logging

```
dfm -verbose install
```

##### dryrun
Prints shell commands that would be excused without excusing them

```
dfm -dryrun install
```

##### overwrite
Overwrite previously existing files, folders or symlinks when linking files

```
dfm -overwrite install
```

##### force
Forces commands to run with more power. Also overwrites previously existing files, folders or symlinks when linking files

```
dfm -force install
```


### Enable cd command
To enable the `dfm cd` command, add the following to your profile (~/.bash_profile, ~/.zshrc, ~/.profile, or ~/.bashrc).

```
function dfm () {
  if [ "$1" == "cd" ]
  then
    cd $(command dfm path)
  	return 0
  fi
  command dfm "$@"
}
```
