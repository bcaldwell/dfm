package dfm

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestIntegration(t *testing.T) {
	Convey("Integration test", t, func() {
		configFile, err := filepath.Abs("../resources/testing/dfm.yml")
		So(err, ShouldEqual, nil)

		dir, err := ioutil.TempDir("", "dfm-test")
		So(err, ShouldEqual, nil)
		defer os.RemoveAll(dir)

		os.Args = []string{"dfm", "--config", configFile, "--destdir", dir, "install"}
		Execute()

		srcDir := filepath.Dir(configFile)

		link, err := os.Readlink(path.Join(dir, "link-test-1"))
		So(err, ShouldEqual, nil)
		expected := path.Join(srcDir, "link-test")
		So(link, ShouldEqual, expected)

		link, err = os.Readlink(path.Join(dir, "link-test"))
		So(err, ShouldEqual, nil)
		expected = path.Join(srcDir, "link-test")
		So(link, ShouldEqual, expected)

		compiledContents, err := ioutil.ReadFile(path.Join(dir, "temp/file_compiled"))
		So(err, ShouldEqual, nil)
		expected = fmt.Sprintf("really important configuration file for %s", runtime.GOOS)
		So(string(compiledContents), ShouldEqual, expected)
	})

}
