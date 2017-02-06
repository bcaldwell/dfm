package tasks

import (
	"errors"
	"os"
	"path"

	"github.com/benjamincaldwell/dfm/templates"
	"github.com/benjamincaldwell/dfm/utilities"
)

var errNoDest = errors.New("template save destination was not specified")

type Template struct {
	Dest           string
	Files          []string
	TemplateString string
	Glob           string
	Vars           map[string]string
}

func (t *Template) isDefined() bool {
	return len(t.Files) > 0 || t.TemplateString != "" || t.Glob != ""
}

func processTemplate(tmpl Template) error {
	if tmpl.Dest == "" {
		return errNoDest
	}
	tpl, err := templates.New(
		templates.Files(tmpl.Files),
		templates.TemplateString(tmpl.TemplateString),
		templates.Glob(tmpl.Glob),
		templates.MergeVariables(tmpl.Vars),
	)
	utilities.ErrorCheck(err, "parsing template arguments")
	// todo...do this right
	dest := path.Join(DestDir, tmpl.Dest)
	f, err := os.Create(dest)
	defer f.Close()
	if err != nil {
		return err
	}
	err = tpl.Execute(f)
	return err
}
