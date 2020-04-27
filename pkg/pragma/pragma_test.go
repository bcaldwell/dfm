package pragma

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"testing"
)

func TestNewFile(t *testing.T) {
	type args struct {
		fileContents string
	}

	tests := []struct {
		name string
		args args
		want *File
	}{
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFile(tt.args.fileContents); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFile_Process(t *testing.T) {
	type fields struct {
		FileContents    string
		PragmaName      string
		CommentString   string
		pragmaLineRegex *regexp.Regexp
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &File{
				FileContents:    tt.fields.FileContents,
				PragmaName:      tt.fields.PragmaName,
				CommentString:   tt.fields.CommentString,
				pragmaLineRegex: tt.fields.pragmaLineRegex,
			}
			if _, err := p.Process(); (err != nil) != tt.wantErr {
				t.Errorf("File.Process() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFile_GenerateRegex(t *testing.T) {
	type fields struct {
		FileContents    string
		PragmaName      string
		CommentString   string
		pragmaLineRegex *regexp.Regexp
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &File{
				FileContents:    tt.fields.FileContents,
				PragmaName:      tt.fields.PragmaName,
				CommentString:   tt.fields.CommentString,
				pragmaLineRegex: tt.fields.pragmaLineRegex,
			}
			if err := f.setupFileForProcessing(); (err != nil) != tt.wantErr {
				t.Errorf("File.generateRegex() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFile_getPragmaForLine(t *testing.T) {
	type fields struct {
		FileContents    string
		PragmaName      string
		CommentString   string
		pragmaLineRegex *regexp.Regexp
	}

	type args struct {
		line string
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
		want1  map[string]string
	}{
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &File{
				FileContents:    tt.fields.FileContents,
				PragmaName:      tt.fields.PragmaName,
				CommentString:   tt.fields.CommentString,
				pragmaLineRegex: tt.fields.pragmaLineRegex,
			}
			got, got1 := f.getPragmaForLine(tt.args.line)
			if got != tt.want {
				t.Errorf("File.getPragmaForLine() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("File.getPragmaForLine() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestFile_getPragmaForLine2(t *testing.T) {
	f := NewFile("")
	fmt.Println(f.setupFileForProcessing())
	fmt.Println(f.getPragmaForLine("# test"))
	fmt.Println(f.getPragmaForLine("# @dfm"))
	fmt.Println(f.getPragmaForLine("# @dfm start"))
	fmt.Println(f.getPragmaForLine("// @dfm start"))
	fmt.Println(f.getPragmaForLine("# @dfm host=test"))
	fmt.Println(f.getPragmaForLine("# @dfm host=test start"))
	fmt.Println(f.getPragmaForLine("# @dfm env=test=test start"))

	f = NewFile(`
	// some comment
	// @dfm start
	something
	again
	// @dfm end
	more
		`)
	fmt.Println(f.Process())

	f = NewFile(`
		// some comment
		// @dfm os=linux
		something
		again
		// @dfm end
		more
			`)
	fmt.Println(f.Process())

	f = NewFile(`
		// some comment
		// @dfm os=linux
		something
		// @dfm os=darwin
		again
		more
			`)
	fmt.Println(f.Process())

	f = NewFile(`
		// some comment
		// @dfm os=linux
		// something
		// @dfm os=darwin
		// again
		more
			`)
	fmt.Println(f.Process())

	f = NewFile(`
		// some comment
		// @dfm os=linux start
		something
		// @dfm os=darwin
		again
		// @dfm end
		more
			`)
	fmt.Println(f.Process())

	f = NewFile(`
		# some comment
		# @dfm os=darwin start
		something
		# @dfm os=linux
		again
		# @dfm end
		more
			`)
	fmt.Println(f.Process())

	f = NewFile(`
	# @dfm os=linux
	# again
	more
		`)
	fmt.Println(f.Process())

	f = NewFile(`
	# @dfm os=linux start
	[mergetool "Kaleidoscope"]
	  # cmd = ksdiff --merge --output \"$MERGED\" --base \"$BASE\" -- \"$LOCAL\" --snapshot \"$REMOTE\" --snapshot
	  # trustexitcode = true
	# [merge]
	  # tool = Kaleidoscope
	# [difftool "Kaleidoscope"]
	  # cmd = ksdiff --partial-changeset --relative-path \"$MERGED\" -- \"$LOCAL\" \"$REMOTE\"
	# @dfm end
	
	[color]
	  diff = auto
	  status = auto
	  branch = auto
	[core]
	  whitespace = fix
		`)
	fmt.Println(f.Process())
}

func TestFile_processPragma(t *testing.T) {
	type fields struct {
		FileContents    string
		PragmaName      string
		CommentString   string
		pragmaLineRegex *regexp.Regexp
		hostname        string
		os              string
	}

	type args struct {
		pragma parsedPragma
	}

	tests := []struct {
		name                  string
		fields                fields
		args                  args
		wantCommentLine       bool
		wantCommentBlockStart bool
		wantCommentBlockEnd   bool
	}{
		{
			name:                  "host",
			fields:                fields{hostname: "test-host", os: "linux"},
			args:                  args{pragma: parsedPragma{"host": "test-host"}},
			wantCommentLine:       false,
			wantCommentBlockStart: false,
			wantCommentBlockEnd:   false,
		},
		{
			name:                  "host and os",
			fields:                fields{hostname: "test-host", os: "linux"},
			args:                  args{pragma: parsedPragma{"host": "test-host", "os": "darwin"}},
			wantCommentLine:       true,
			wantCommentBlockStart: false,
			wantCommentBlockEnd:   false,
		},
		{
			name:                  "host and os negative",
			fields:                fields{hostname: "test-host", os: "linux"},
			args:                  args{pragma: parsedPragma{"host": "test-host", "os": "linux"}},
			wantCommentLine:       false,
			wantCommentBlockStart: false,
			wantCommentBlockEnd:   false,
		},
		{
			name:                  "host conditional start",
			fields:                fields{hostname: "test-host", os: "linux"},
			args:                  args{pragma: parsedPragma{"host": "test-host", "start": ""}},
			wantCommentLine:       false,
			wantCommentBlockStart: true,
			wantCommentBlockEnd:   false,
		},
		{
			name:                  "env set",
			fields:                fields{hostname: "test-host", os: "linux"},
			args:                  args{pragma: parsedPragma{"env": "TEST=set"}},
			wantCommentLine:       false,
			wantCommentBlockStart: false,
			wantCommentBlockEnd:   false,
		},
		{
			name:                  "env unset",
			fields:                fields{hostname: "test-host", os: "linux"},
			args:                  args{pragma: parsedPragma{"env": "TEST=unset"}},
			wantCommentLine:       true,
			wantCommentBlockStart: false,
			wantCommentBlockEnd:   false,
		},
		{
			name:                  "end",
			fields:                fields{hostname: "test-host", os: "linux"},
			args:                  args{pragma: parsedPragma{"end": ""}},
			wantCommentLine:       false,
			wantCommentBlockStart: false,
			wantCommentBlockEnd:   true,
		},
		{
			name:                  "host conditional end",
			fields:                fields{hostname: "test-host", os: "linux"},
			args:                  args{pragma: parsedPragma{"host": "test-host", "end": ""}},
			wantCommentLine:       false,
			wantCommentBlockStart: false,
			wantCommentBlockEnd:   true,
		},
		{
			name:                  "negative host conditional end",
			fields:                fields{hostname: "test-host", os: "linux"},
			args:                  args{pragma: parsedPragma{"host": "test-host-wrong", "end": ""}},
			wantCommentLine:       true,
			wantCommentBlockStart: false,
			wantCommentBlockEnd:   false,
		},
	}

	os.Setenv("TEST", "set")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &File{
				FileContents:    tt.fields.FileContents,
				PragmaName:      tt.fields.PragmaName,
				CommentString:   tt.fields.CommentString,
				pragmaLineRegex: tt.fields.pragmaLineRegex,
				hostname:        tt.fields.hostname,
				os:              tt.fields.os,
			}
			gotCommentLine, gotCommentBlockStart, gotCommentBlockEnd := f.processPragma(tt.args.pragma)
			if gotCommentLine != tt.wantCommentLine {
				t.Errorf("File.processPragma() gotCommentLine = %v, want %v", gotCommentLine, tt.wantCommentLine)
			}
			if gotCommentBlockStart != tt.wantCommentBlockStart {
				t.Errorf("File.processPragma() gotCommentBlockStart = %v, want %v", gotCommentBlockStart, tt.wantCommentBlockStart)
			}
			if gotCommentBlockEnd != tt.wantCommentBlockEnd {
				t.Errorf("File.processPragma() gotCommentBlockEnd = %v, want %v", gotCommentBlockEnd, tt.wantCommentBlockEnd)
			}
		})
	}
}
