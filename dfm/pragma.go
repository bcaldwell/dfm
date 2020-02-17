package dfm

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/bcaldwell/dfm/pkg/pragma"
)

func pragmaAction(args []string, config *Configuration) error {
	return filepath.Walk(config.SrcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) == ".png" {
			return nil
		}
		if strings.HasPrefix(path, "/home/bcaldwell/.dotfiles/.git") || strings.HasPrefix(path, "/home/bcaldwell/.dotfiles/macos") {
			return nil
		}

		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		p := pragma.NewFile(string(b))
		s, err := p.Process()
		if err != nil {
			return nil
		}

		err = ioutil.WriteFile(path, []byte(s), info.Mode())
		if err != nil {
			return err
		}

		return nil
	})
}
