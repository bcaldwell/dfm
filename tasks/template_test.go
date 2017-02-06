package tasks

// import "testing"

// func Test_template_isDefined(t *testing.T) {
// 	type fields struct {
// 		Dest           string
// 		Files          []string
// 		TemplateString string
// 		Glob           string
// 		Vars           map[string]string
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		want   bool
// 	}{}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			t := &template{
// 				Dest:           tt.fields.Dest,
// 				Files:          tt.fields.Files,
// 				TemplateString: tt.fields.TemplateString,
// 				Glob:           tt.fields.Glob,
// 				Vars:           tt.fields.Vars,
// 			}
// 			if got := t.isDefined(); got != tt.want {
// 				t.Errorf("template.isDefined() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func Test_processTemplate(t *testing.T) {
// 	type args struct {
// 		tmpl template
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
// 			if err := processTemplate(tt.args.tmpl); (err != nil) != tt.wantErr {
// 				t.Errorf("processTemplate() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }
