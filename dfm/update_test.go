package dfm

import (
	"errors"
	"testing"

	"github.com/benjamincaldwell/go-sh/mock"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_updateAction(t *testing.T) {
	gitErr := errors.New("Not a git repo")

	type args struct {
		args   []string
		config *Configuration
	}
	arg := args{
		[]string{},
		&Configuration{
			SrcDir: "/src",
		},
	}

	tests := []struct {
		name           string
		args           args
		wantErr        error
		ErrorsToReturn []error
		shellCommands  []string
	}{
		{
			"git fetch fails",
			arg,
			gitErr,
			[]error{gitErr},
			[]string{"git fetch"},
		},
		{
			"git pull fails",
			arg,
			gitErr,
			[]error{nil, gitErr},
			[]string{"git fetch", "git pull"},
		},
		{
			"No errors",
			arg,
			nil,
			[]error{},
			[]string{"git fetch", "git pull"},
		},
	}
	for _, tt := range tests {
		Convey(tt.name, t, func() {
			shMock.Reset()
			shMock.ErrorsToReturn = tt.ErrorsToReturn
			err := updateAction(tt.args.args, tt.args.config)
			So(err, ShouldEqual, tt.wantErr)
			So(len(shMock.Commands), ShouldEqual, len(tt.shellCommands))
			So(shMock.Commands, ShouldResemble, tt.shellCommands)
		})
	}
}
