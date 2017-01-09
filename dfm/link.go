package dfm

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/benjamincaldwell/devctl/printer"
)

func parseLink(params string, config *Configuration) (src string, dest string, mode os.FileMode, err error) {
	paramParts := strings.Split(params, ":")
	if len(paramParts) > 3 {
		printer.Info("droping last part of link param %s from %s", strings.Join(paramParts[3:], ":"), params)
	}
	if len(paramParts) == 3 {
		var parsed uint64
		parsed, err = strconv.ParseUint(paramParts[2], 0, 64)
		if err != nil {
			return
		}

		mode = os.FileMode(parsed)
	} else {
		mode = 0755
	}

	src = os.ExpandEnv(paramParts[0])

	if len(paramParts) == 2 {
		dest = os.ExpandEnv(paramParts[1])
	} else {
		dest = path.Join(config.DestDir, path.Base(src))
	}
	return
}

func processLink(params string, config *Configuration) error {
	src, dest, mode, err := parseLink(params, config)
	if err == nil {
		return createLink(src, dest, mode, config)
	}
	return err
}

func createLink(src string, dest string, mode os.FileMode, config *Configuration) error {
	srcAbs := src
	if !path.IsAbs(srcAbs) {
		srcAbs = path.Join(config.SrcDir, src)
	}

	destAbs := dest
	if !path.IsAbs(destAbs) {
		destAbs = path.Join(config.DestDir, dest)
	}

	if _, err := os.Stat(srcAbs); os.IsNotExist(err) {
		return fmt.Errorf("source path %s does not exist in filesystem", srcAbs)
	}

	if _, err := os.Stat(destAbs); err == nil {
		if *force || *overwrite {
			if *dryRun || *verbose {
				printer.WarningBar("removing %s", destAbs)

			} else {
				os.RemoveAll(destAbs)
			}
		} else {
			printer.VerboseInfoBar("Linking %s to %s", srcAbs, destAbs)
			printer.VerboseWarningBar("%s already exists. Use overwrite or force option to overwrite", destAbs)
			return nil
		}
	}

	if *dryRun || *verbose {
		printer.InfoBar("Linking %s to %s", srcAbs, destAbs)
		if *dryRun {
			return nil
		}
	}

	return os.Symlink(srcAbs, destAbs)
}
