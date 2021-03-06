package dfm

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/bcaldwell/dfm/pkg/pragma"
	"github.com/bcaldwell/go-printer"
)

const (
	ignoreFile = ".dfmpragmaignore"
)

func pragmaAction(args []string, config *Configuration) error {
	return applyPragmaToFolder(config.SrcDir)
}

func applyPragmaToFolder(folder string) error {
	if filepath.Base(folder) == ".git" {
		return nil
	}

	if _, err := os.Stat(filepath.Join(folder, ignoreFile)); !os.IsNotExist(err) {
		return nil
	}

	files, err := ioutil.ReadDir(folder)
	if err != nil {
		return err
	}

	for _, f := range files {
		filePath := path.Join(folder, f.Name())

		if f.IsDir() {
			err = applyPragmaToFolder(filePath)
			if err != nil {
				return err
			}
		} else {
			printer.InfoBar("Processing pragma on %s", filePath)
			err = processPragmaFile(filePath, f.Mode())
			if err != nil {
				printer.ErrorBar("error processing file %v", err)
			}
		}
	}

	return nil
}

func processPragmaFile(file string, mode os.FileMode) error {
	folderName := filepath.Base(file)
	if _, err := os.Stat(filepath.Join(folderName, ignoreFile)); !os.IsNotExist(err) {
		return err
	}

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	p := pragma.NewFile(string(b))

	s, err := p.Process()
	if err != nil {
		return nil
	}

	err = ioutil.WriteFile(file, []byte(s), mode)
	if err != nil {
		return err
	}

	return nil
}
