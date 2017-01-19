package templates

import (
	"bytes"
	"html/template"
	"path"

	"io"

	"github.com/Masterminds/sprig"
)

var funcMap = template.FuncMap{}

func init() {
	funcMap["loadTemplate"] = loadTemplate
	sprigMap := sprig.FuncMap()
	for k, v := range sprigMap {
		funcMap[k] = v
	}
}

func ExecuteTemplate(wr io.Writer, files ...string) {
	tpl := template.Must(template.New(path.Base(files[0])).Funcs(funcMap).ParseFiles(files...))
	tplVars := map[string]string{}
	tpl.Execute(wr, tplVars)
}

func loadTemplate(files ...string) string {
	// t, err := template.New(path.Base(files[0])).Funcs(funcMap).ParseFiles(files...)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	var doc bytes.Buffer
	// if err := t.Execute(&doc, nil); err != nil {
	// 	log.Fatal(err)
	// }
	ExecuteTemplate(&doc, files...)
	return doc.String()
}
