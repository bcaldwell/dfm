package templates

import (
	"html/template"
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCreationAndSetOptions(t *testing.T) {
	Convey("creating and seting options on Tpl struct", t, func() {
		funcMap := template.FuncMap{
			"hi": func() {},
		}
		tmp, err := New(
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
				AppendVariable("second", "false"),
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
	})
}

func Test_loadTemplate(t *testing.T) {

}
