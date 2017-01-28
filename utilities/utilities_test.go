package utilities

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var testError = errors.New("error for testing")

func TestStringInSlice(t *testing.T) {
	Convey("Should return true if string is in slice", t, func() {
		So(StringInSlice("hi", []string{"hi", "bye", "person"}), ShouldEqual, true)
		So(StringInSlice("hi", []string{"hi", "hi", "person"}), ShouldEqual, true)
		So(StringInSlice("hi", []string{"bye", "hi", "person", "friend"}), ShouldEqual, true)
	})

	Convey("Should return false if string is not in slice", t, func() {
		So(StringInSlice("friend", []string{"hi", "bye", "person"}), ShouldEqual, false)
		So(StringInSlice("poppie", []string{"hi", "hi", "hi"}), ShouldEqual, false)
	})
}

func TestUniqueSliceTransform(t *testing.T) {
	Convey("Should return an array of only unique strings", t, func() {
		tests := []struct {
			arg  []string
			want []string
		}{
			{
				[]string{"hi", "hi", "hi"},
				[]string{"hi"},
			},
			{
				[]string{"hi", "bye", "hi"},
				[]string{"hi", "bye"},
			},
			{
				[]string{"hi", "bye", "dog"},
				[]string{"hi", "bye", "dog"},
			},
		}
		for _, tt := range tests {
			So(UniqueSliceTransform(tt.arg), ShouldResemble, tt.want)
		}
	})
}

func TestAppendIfUnique(t *testing.T) {
	Convey("Should only append if string ins't already in the slice", t, func() {
		type args struct {
			a []string
			i string
		}
		tests := []struct {
			arg  args
			want []string
		}{
			{
				args{
					[]string{"hi", "hi", "hi"},
					"sup",
				},
				[]string{"hi", "hi", "hi", "sup"},
			},
			{
				args{
					[]string{"hi", "bye", "hi"},
					"dog",
				},
				[]string{"hi", "bye", "hi", "dog"},
			},
		}
		for _, tt := range tests {
			So(AppendIfUnique(tt.arg.a, tt.arg.i), ShouldResemble, tt.want)
		}
	})
}

func TestRunFatalErrorCheck(t *testing.T) {
	fmt.Println("sup")
	if os.Getenv("TEST_FATAL_ERROR_CHECK") == "1" {
		fmt.Println("hiii")
		var err error
		val := os.Getenv("TEST_FATAL_ERROR_CHECK_ERROR_VALUE")
		switch val {
		case "":
			err = nil
		default:
			err = errors.New(val)
		}
		FatalErrorCheck(err, "testing error")
		return
	}
}

func TestFatalErrorCheck(t *testing.T) {
	Convey("Shouldn't crash if error is nil", t, func() {
		cmd := exec.Command(os.Args[0], "-test.run=TestRunFatalErrorCheck")
		cmd.Env = append(os.Environ(), "TEST_FATAL_ERROR_CHECK=1")
		err := cmd.Run()
		_, ok := err.(*exec.ExitError)
		So(ok, ShouldEqual, false)
	})

	Convey("Should crash if error is not nil", t, func() {
		cmd := exec.Command(os.Args[0], "-test.run=TestRunFatalErrorCheck")
		cmd.Env = append(os.Environ(), "TEST_FATAL_ERROR_CHECK=1", "TEST_FATAL_ERROR_CHECK_ERROR_VALUE=erroring")
		err := cmd.Run()
		e, ok := err.(*exec.ExitError)
		fmt.Printf("%+v, %+v", e, ok)
		So(ok, ShouldEqual, true)
		So(e.Success(), ShouldEqual, false)
	})
}

func TestErrorCheck(t *testing.T) {
	Convey("Should return the bool of if an error occur and print out message if it did", t, func() {

		type args struct {
			err     error
			message string
		}
		tests := []struct {
			args args
			want bool
		}{
			{
				args{
					testError,
					"error should exist",
				},
				true,
			},
			{
				args{
					nil,
					"error should not exist",
				},
				false,
			},
		}
		for _, tt := range tests {
			So(ErrorCheck(tt.args.err, tt.args.message), ShouldEqual, tt.want)
		}
	})
}

// func TestAskForConfirmation(t *testing.T) {
// 	type args struct {
// 		s string
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want bool
// 	}{
// 	// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := AskForConfirmation(tt.args.s); got != tt.want {
// 				t.Errorf("AskForConfirmation() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestHTTPDownload(t *testing.T) {
// 	type args struct {
// 		uri string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    []byte
// 		wantErr bool
// 	}{
// 	// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := HTTPDownload(tt.args.uri)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("HTTPDownload() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("HTTPDownload() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestWriteFile(t *testing.T) {
// 	type args struct {
// 		dst string
// 		d   []byte
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 	// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if err := WriteFile(tt.args.dst, tt.args.d); (err != nil) != tt.wantErr {
// 				t.Errorf("WriteFile() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

// func TestDownloadToFile(t *testing.T) {
// 	type args struct {
// 		uri string
// 		dst string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 	// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if err := DownloadToFile(tt.args.uri, tt.args.dst); (err != nil) != tt.wantErr {
// 				t.Errorf("DownloadToFile() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }
