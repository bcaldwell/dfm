package templates

import (
	"bytes"
	"html/template"

	"github.com/Masterminds/sprig"
)

// Tpl is the template struct the contains template options
type Tpl struct {
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

// AppendVariable appends the variable parameter of the tmp struct. Variables are passed into the template
func AppendVariable(name, value string) Option {
	return func(tmpl *Tpl) (err error) {
		tmpl.Variables[name] = value
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

// ExecuteTemplate loads the parses the passed in files as the template.
// The template contents are writtent to wr. And error is returned if any occur
// func ExecuteTemplate(wr io.Writer, files ...string) (err error) {
// 	Tpl, err := template.New(path.Base(files[0])).Funcs(funcMap).ParseFiles(files...)
// 	if utilities.ErrorCheck(err, "Load template "+files[0]) {
// 		return err
// 	}
// 	TplVars := map[string]string{}
// 	err = Tpl.Execute(wr, TplVars)
// 	return err
// }

func loadAndExecuteTemplate(files ...string) string {
	var doc bytes.Buffer
	// ExecuteTemplate(&doc, files...)
	return doc.String()
}
