package dfm

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestExecute(t *testing.T) {
	tests := []struct {
		name string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Execute()
		})
	}
}

func Test_detectHomeDir(t *testing.T) {
	Convey("Should return value of HOME if it exists", t, func() {
		os.Setenv("HOME", "testing_dir")
		So(detectHomeDir(), ShouldEqual, "testing_dir")
		// os.Setenv("HOME", "")
		// So(detectHomeDir(), ShouldEqual, "testing_dir")
	})
}

func Test_printFlagOptions(t *testing.T) {
	tests := []struct {
		name string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			printFlagOptions()
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
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotConfigFile := detectConfigFile(tt.args.configFileFlag, tt.args.homeDir); gotConfigFile != tt.wantConfigFile {
				t.Errorf("detectConfigFile() = %v, want %v", gotConfigFile, tt.wantConfigFile)
			}
		})
	}
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
	type args struct {
		homeDir string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := determineRcFile(tt.args.homeDir); got != tt.want {
				t.Errorf("determineRcFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cloneRepo(t *testing.T) {
	type args struct {
		repo   string
		srcDir string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := cloneRepo(tt.args.repo, tt.args.srcDir); (err != nil) != tt.wantErr {
				t.Errorf("cloneRepo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_detectDefaultConfigFileLocation(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := detectDefaultConfigFileLocation(); got != tt.want {
				t.Errorf("detectDefaultConfigFileLocation() = %v, want %v", got, tt.want)
			}
		})
	}
}
