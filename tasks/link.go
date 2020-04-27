package tasks

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/bcaldwell/dfm/pkg/pragma"
	"github.com/bcaldwell/go-printer"
)

func parseLink(params string) (string, string, os.FileMode, error) {
	var src, dest string

	var mode os.FileMode

	paramParts := strings.Split(params, ":")
	if len(paramParts) > 3 {
		printer.Info("droping last part of link param %s from %s", strings.Join(paramParts[3:], ":"), params)
	}

	if len(paramParts) == 3 {
		var parsed uint64

		parsed, err := strconv.ParseUint(paramParts[2], 0, 64)
		if err != nil {
			return src, dest, mode, err
		}

		mode = os.FileMode(parsed)
	} else {
		mode = 0755
	}

	src = os.ExpandEnv(paramParts[0])

	if len(paramParts) == 2 {
		dest = os.ExpandEnv(paramParts[1])
	} else {
		dest = path.Join(DestDir, path.Base(src))
	}

	return src, dest, mode, nil
}

func processLink(params string, vars DfmVars) error {
	src, dest, mode, err := parseLink(params)
	if err != nil {
		return fmt.Errorf("failed to parse link %w", err)
	}

	srcAbs := absPath(src, SrcDir)

	err = pragma.ProcessFile(srcAbs, vars)
	if err != nil {
		return fmt.Errorf("failed to process pragma on link %s->%s %w", src, dest, err)
	}

	return createLink(src, dest, mode)
}

func createLink(src string, dest string, mode os.FileMode) error {
	srcAbs := absPath(src, SrcDir)

	destAbs := absPath(dest, DestDir)

	if _, err := os.Stat(srcAbs); os.IsNotExist(err) {
		return fmt.Errorf("source path %s does not exist in filesystem", srcAbs)
	}

	if _, err := os.Stat(destAbs); err == nil {
		if Force || Overwrite {
			if DryRun || Verbose {
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

	if DryRun || Verbose {
		printer.InfoBar("Linking %s to %s", srcAbs, destAbs)

		if DryRun {
			return nil
		}
	}

	return os.Symlink(srcAbs, destAbs)
}
