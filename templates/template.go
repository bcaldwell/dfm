package templates

import (
	"bytes"
	"html/template"
	"path/filepath"
	"runtime"

	"io"

	"github.com/Masterminds/sprig"
	"github.com/benjamincaldwell/dfm/utilities"
)

// Tpl is the template struct the contains template options
type Tpl struct {
	Name           string
	TemplateString string
	Files          []string
	Glob           string
	Variables      map[string]string
	FuncMap        template.FuncMap
}

// Option is an option function used by SetOptions to configure the struct
type Option func(*Tpl) error

// New returns a new Tpl struct in`itialized with the option arguements pass in
func New(options ...Option) (*Tpl, error) {
	tmpl := &Tpl{}
	defaultOptions := []Option{DefaultVariables, DefaultFuncMap}
	tmpl.SetOptions(defaultOptions...)
	return tmpl.SetOptions(options...)
}

// SetOptions runs the passed in option functions and returns an error is one occurs
func (tmpl *Tpl) SetOptions(options ...Option) (*Tpl, error) {
	for _, opt := range options {
		if err := opt(tmpl); err != nil {
			return tmpl, err
		}
	}
	return tmpl, nil
}

// Name sets name parameter of the tmp struct. Name is the name of the template
func Name(name string) Option {
	return func(tmpl *Tpl) (err error) {
		tmpl.Name = name
		return nil
	}
}

// TemplateString sets parameter of the tmp struct. TemplateString is the parsed as a string template
func TemplateString(templateString string) Option {
	return func(tmpl *Tpl) (err error) {
		tmpl.TemplateString = templateString
		return nil
	}
}

// Files sets parameter of the tmp struct. Files are the files parsed as templates
func Files(files []string) Option {
	return func(tmpl *Tpl) (err error) {
		tmpl.Files = files
		return nil
	}
}

// AppendFiles appends file parameter of the tmp struct. Files are the files parsed as templates
func AppendFiles(file string) Option {
	return func(tmpl *Tpl) (err error) {
		tmpl.Files = append(tmpl.Files, file)
		return nil
	}
}

// Glob sets parameter of the tmp struct. Files identified by the pattern are used as the templates
func Glob(files string) Option {
	return func(tmpl *Tpl) (err error) {
		tmpl.Glob = files
		return nil
	}
}

// Variables sets parameter of the tmp struct. Variables are passed into the template
func Variables(variable map[string]string) Option {
	return func(tmpl *Tpl) (err error) {
		tmpl.Variables = variable
		return nil
	}
}

// AppendVariables appends the variable parameter of the tmp struct. Variables are passed into the template
func AppendVariables(name, value string) Option {
	return func(tmpl *Tpl) (err error) {
		tmpl.Variables[name] = value
		return nil
	}
}

// MergeVariables merges the variable parameter of the tmp struct with the passed in map. Variables are passed into the template
func MergeVariables(m map[string]string) Option {
	return func(tmpl *Tpl) (err error) {
		for k, v := range m {
			tmpl.Variables[k] = v
		}
		return nil
	}
}

// FuncMap sets the funcMap parameter of the tmp struct. FuncMap functions are passed into the template
func FuncMap(funcMap template.FuncMap) Option {
	return func(tmpl *Tpl) (err error) {
		tmpl.FuncMap = funcMap
		return nil
	}
}

// AppendFuncMap appends the funcMap parameter of the tmp struct. FuncMap functions are passed into the template
func AppendFuncMap(name string, value interface{}) Option {
	return func(tmpl *Tpl) (err error) {
		tmpl.FuncMap[name] = value
		return nil
	}
}

// DefaultVariables is an option function that sets the default values of Variables map
func DefaultVariables(tmpl *Tpl) (err error) {
	tmpl.Variables = map[string]string{
		"os": runtime.GOOS,
	}
	return nil
}

// DefaultFuncMap is an option function that sets the default values of the FuncMap
func DefaultFuncMap(tmpl *Tpl) (err error) {
	funcMap := template.FuncMap{}
	funcMap["render"] = loadAndExecuteTemplate
	// append sprig functions to funcMap
	sprigMap := sprig.FuncMap()
	for k, v := range sprigMap {
		funcMap[k] = v
	}
	tmpl.FuncMap = funcMap
	return nil
}

// Execute parses and Execute the template as configured by the options.
// The output is written to the passed in writter
func (tmpl *Tpl) Execute(wr io.Writer) (err error) {
	name := tmpl.Name
	if name == "" {
		if tmpl.TemplateString != "" {
			name = "template"
		} else if len(tmpl.Files) > 0 {
			name = filepath.Base(tmpl.Files[0])
		} else if tmpl.Glob != "" {
			if filenames, err := filepath.Glob(tmpl.Glob); err == nil && len(filenames) > 0 {
				name = filepath.Base(filenames[0])
			}
		}
	}

	temp := template.New(name).Funcs(tmpl.FuncMap)
	if tmpl.TemplateString != "" {
		temp, err = temp.Parse(tmpl.TemplateString)
	} else if len(tmpl.Files) > 0 {
		temp, err = temp.ParseFiles(tmpl.Files...)
	} else if tmpl.Glob != "" {
		temp, err = temp.ParseGlob(tmpl.Glob)
	}

	if utilities.ErrorCheck(err, "Parse template") {
		return err
	}

	err = temp.Funcs(tmpl.FuncMap).Execute(wr, tmpl.Variables)
	return err
}

func loadAndExecuteTemplate(files ...string) string {
	var doc bytes.Buffer
	if temp, err := New(Files(files)); err == nil {
		temp.Execute(&doc)
	}
	return doc.String()
}
