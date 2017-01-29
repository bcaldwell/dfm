package dfm

import (
	"os"
	"path"
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/afero"
)

func TestExecute(t *testing.T) {
}

func Test_detectHomeDir(t *testing.T) {
	Convey("Should return value of HOME if it exists", t, func() {
		os.Setenv("HOME", "testing_dir")
		home, err := detectHomeDir()
		So(home, ShouldEqual, "testing_dir")
		So(err, ShouldEqual, nil)
	})
	Convey("Should return an error is HOME env is blank", t, func() {
		os.Setenv("HOME", "")
		home, err := detectHomeDir()
		So(home, ShouldEqual, "")
		So(err, ShouldEqual, ErrNoHomeEnv)
	})
}

func Test_printFlagOptions(t *testing.T) {
	Convey("Should print flag options if in verbose", t, func() {
		dryRun = true
		force = true
		overwrite = true
		printFlagOptions()
	})
}

func Test_createDfmrc(t *testing.T) {
	type args struct {
		homeDir    string
		configFile string
		scrDir     string
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createDfmrc(tt.args.homeDir, tt.args.configFile, tt.args.scrDir)
		})
	}
}

func Test_determineRcFile(t *testing.T) {
	Convey("Should return homeDir with .dfmrc appended", t, func() {
		So(determineRcFile("/testing"), ShouldEqual, "/testing/.dfmrc")
	})
}

func Test_cloneRepo(t *testing.T) {
	srcDir, err := filepath.Abs("./testing/src")
	fs.MkdirAll(srcDir, 0755)
	defer fs.RemoveAll("./testing")
	Convey("Should clone given repo to given source directory", t, func() {
		So(err, ShouldEqual, nil)
		cloneRepo("git@github.com:benjamincaldwell/public-test.git", srcDir)
		_, err := fs.Stat(path.Join(srcDir, "testing-file"))
		So(err, ShouldEqual, nil)
	})

	Convey("Should return an error if the clone failed", t, func() {
		So(err, ShouldEqual, nil)
		err := cloneRepo("git@github.com:benjamincaldwell/public-doesnt-exist.git", srcDir)
		So(err, ShouldNotEqual, nil)
	})
}

func Test_detectDefaultConfigFileLocation(t *testing.T) {
	tests := []struct {
		name  string
		files []string
		want  string
		err   error
	}{
		{
			"detect local dfm.yml",
			[]string{"dfm.yml", "$HOME/.dotfiles/dfm.yml", "$HOME/dfm.yml", "$HOME/.dfm.yml"},
			"dfm.yml",
			nil,
		},
		{
			"detect home dfm.yml in .dotfile folder",
			[]string{"$HOME/.dotfiles/dfm.yml", "$HOME/dfm.yml", "$HOME/.dfm.yml"},
			"/src/.dotfiles/dfm.yml",
			nil,
		},
		{
			"detect home dfm.yml in .dotfile folder",
			[]string{"$HOME/dotfiles/dfm.yml", "$HOME/dfm.yml", "$HOME/.dfm.yml"},
			"/src/dotfiles/dfm.yml",
			nil,
		},
		{
			"detect home dfm.yml",
			[]string{"$HOME/dfm.yml", "$HOME/.dfm.yml"},
			"/src/dfm.yml",
			nil,
		},
		{
			"detect home .dfm.yml",
			[]string{"$HOME/.dfm.yml"},
			"/src/.dfm.yml",
			nil,
		},
		{
			"return error is none found",
			[]string{},
			"",
			ErrNoConfigFile,
		},
	}
	for _, tt := range tests {
		Convey(tt.name, t, func() {
			fs = afero.NewMemMapFs()
			os.Setenv("HOME", "/src")
			fs.MkdirAll("/src", 0755)
			for _, file := range tt.files {
				file = os.ExpandEnv(file)
				fs.Create(file)
			}
			file, err := detectDefaultConfigFileLocation()
			So(file, ShouldEqual, tt.want)
			So(err, ShouldEqual, tt.err)
		})
	}
}

func Test_detectConfigFile(t *testing.T) {
	type args struct {
		configFileFlag string
		homeDir        string
	}
	tests := []struct {
		name           string
		args           args
		wantConfigFile string
		wantErr        bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotConfigFile, err := detectConfigFile(tt.args.configFileFlag, tt.args.homeDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("detectConfigFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotConfigFile != tt.wantConfigFile {
				t.Errorf("detectConfigFile() = %v, want %v", gotConfigFile, tt.wantConfigFile)
			}
		})
	}
}
