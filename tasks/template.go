package tasks

import (
	"errors"
	"os"

	"path/filepath"

	"github.com/bcaldwell/dfm/templates"
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

	glob := tmpl.Glob
	if glob != "" {
		glob = absPath(glob, SrcDir)
	}

	tpl, err := templates.New(
		templates.TemplateString(tmpl.TemplateString),
		templates.Glob(glob),
		templates.MergeVariables(tmpl.Vars),
	)
	if err != nil {
		return err
	}

	for _, file := range tmpl.Files {
		file = absPath(file, SrcDir)
		tpl, err = tpl.SetOptions(
			templates.AppendFiles(file),
		)

		if err != nil {
			return err
		}
	}

	dest := absPath(tmpl.Dest, DestDir)

	err = os.MkdirAll(filepath.Dir(dest), 0755)
	if err != nil {
		return err
	}

	f, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer f.Close()

	err = tpl.Execute(f)

	return err
}
