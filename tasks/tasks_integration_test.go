package tasks

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/ghodss/yaml"

	"github.com/benjamincaldwell/dfm/templates"
	. "github.com/smartystreets/goconvey/convey"
)

func TestIntegration(t *testing.T) {
	Convey("Integration test", t, func() {
		tasks := new(struct {
			Tasks map[string]Task
		})
		configFile, err := filepath.Abs("../resources/testing/dfm.yml")
		So(err, ShouldEqual, nil)
		tpl, err := templates.New(
			templates.AppendFiles(configFile),
		)
		var data bytes.Buffer
		tpl.Execute(&data)
		err = yaml.Unmarshal(data.Bytes(), &tasks)
		So(err, ShouldEqual, nil)

		SrcDir = path.Dir(configFile)

		dir, err := ioutil.TempDir("", "dfm-test")
		So(err, ShouldEqual, nil)
		// defer os.RemoveAll(dir)
		fmt.Println(dir)

		DestDir = dir

		err = ExecuteTasks(tasks.Tasks, "")
		So(err, ShouldEqual, nil)

		link, err := os.Readlink(path.Join(dir, "link-test-1"))
		So(err, ShouldEqual, nil)
		expected := path.Join(SrcDir, "link-test")
		So(link, ShouldEqual, expected)

		link, err = os.Readlink(path.Join(dir, "link-test"))
		So(err, ShouldEqual, nil)
		expected = path.Join(SrcDir, "link-test")
		So(link, ShouldEqual, expected)

		compiledContents, err := ioutil.ReadFile(path.Join(dir, "temp/file_compiled"))
		So(err, ShouldEqual, nil)
		expected = fmt.Sprintf("really important configuration file for %s", runtime.GOOS)
		So(string(compiledContents), ShouldEqual, expected)
	})

}
