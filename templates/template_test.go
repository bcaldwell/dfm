package templates

import (
	"bytes"
	"errors"
	"html/template"
	"reflect"
	"runtime"
	"testing"

	"fmt"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCreationAndSetOptions(t *testing.T) {
	Convey("creating and seting options on Tpl struct", t, func() {
		funcMap := template.FuncMap{
			"hi": func() {},
		}
		tmp, err := New(
			Name("hi"),
			TemplateString("tpl string"),
			Files([]string{"fav.go", "silly.go"}),
			Glob("*.go"),
			Variables(map[string]string{
				"test": "true",
			}),
			FuncMap(funcMap),
		)
		Convey("New should create a new tmp struct with options applied", func() {
			expected := &Tpl{
				Name:           "hi",
				TemplateString: "tpl string",
				Files:          []string{"fav.go", "silly.go"},
				Glob:           "*.go",
				Variables: map[string]string{
					"test": "true",
				},
				FuncMap: funcMap,
			}
			So(err, ShouldEqual, nil)
			So(tmp, ShouldResemble, expected)
		})
		Convey("Seting options using append methods", func() {
			expected := &Tpl{
				Name:           "hi",
				TemplateString: "tpl string",
				Files:          []string{"fav.go", "silly.go", "extra.go"},
				Glob:           "*.go",
				Variables: map[string]string{
					"test":   "true",
					"second": "false",
				},
				FuncMap: funcMap,
			}
			extraFunc := func() {}
			expected.FuncMap["extra"] = extraFunc
			tmp.SetOptions(
				AppendFiles("extra.go"),
				AppendVariables("second", "false"),
				AppendFuncMap("extra", extraFunc),
			)
			So(err, ShouldEqual, nil)
			So(tmp, ShouldResemble, expected)
		})
		Convey("New should set default values", func() {
			tmp, err = New()
			So(err, ShouldEqual, nil)
			funcMapKeys := reflect.ValueOf(tmp.FuncMap).MapKeys()
			So(len(funcMapKeys), ShouldBeGreaterThan, 0)
			So(tmp.FuncMap, ShouldContainKey, "render")
		})
		Convey("Should return error if one occurs", func() {
			someErr := errors.New("yay")
			_, err := tmp.SetOptions(func(_ *Tpl) (err error) {
				return someErr
			})
			So(err, ShouldEqual, someErr)
		})
	})
}

func Test_Tpl_Execute(t *testing.T) {
	Convey("Parse and execute template string and printer result to a writer", t, func() {
		tests := []struct {
			name         string
			options      []Option
			expect       string
			newError     func(actual interface{}, expected ...interface{}) (message string)
			executeError func(actual interface{}, expected ...interface{}) (message string)
		}{
			{
				"Basic templateString",
				[]Option{TemplateString("<h1>{{.os}}</h1>")},
				fmt.Sprintf("<h1>%s</h1>", runtime.GOOS),
				ShouldEqual,
				ShouldEqual,
			},
			{
				"Basic template file",
				[]Option{AppendFiles("../resources/testing.tpl")},
				"HELLO!HELLO!HELLO!HELLO!HELLO!",
				ShouldEqual,
				ShouldEqual,
			},
			{
				"Basic Glob",
				[]Option{Glob("../resources/*.tpl")},
				"HELLO!HELLO!HELLO!HELLO!HELLO!",
				ShouldEqual,
				ShouldEqual,
			},
			{
				"Bad glob returns error",
				[]Option{Glob("../fghjk/*.tpl")},
				"",
				ShouldEqual,
				ShouldNotEqual,
			},
			{
				"error if template not found",
				[]Option{AppendFiles("fghjk.tpl")},
				"",
				ShouldEqual,
				ShouldNotEqual,
			},
			{
				"Default function map is set",
				[]Option{TemplateString("{{ \"hello!\" | upper | repeat 5 }}")},
				"HELLO!HELLO!HELLO!HELLO!HELLO!",
				ShouldEqual,
				ShouldEqual,
			},
			{
				"Using render function",
				[]Option{TemplateString("<h1>{{render \"../resources/testing.tpl\"}}</h1>")},
				"<h1>HELLO!HELLO!HELLO!HELLO!HELLO!</h1>",
				ShouldEqual,
				ShouldEqual,
			},
		}

		for _, tt := range tests {
			Convey(tt.name, func() {
				var result bytes.Buffer
				temp, err := New(tt.options...)
				So(err, tt.newError, nil)
				err = temp.Execute(&result)
				So(err, tt.executeError, nil)
				So(result.String(), ShouldEqual, tt.expect)
			})
		}
	})
}
