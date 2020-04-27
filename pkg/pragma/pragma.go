package pragma

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"runtime"
	"strings"
	"unicode"

	"github.com/bcaldwell/go-printer"
)

// This package add support for commenting in or out lines based on a pragma
// to activate a pragma use a comment followed by @name pragma. See example below
// @dfm os=linux
// some line that get un-commented if os is linux
// supported pragmas:
// - host --> based on hostname
// - os --> based on os
// - env --> based on env
// - start/end -> allows for a whole block to be commented out
// The comment method is inferred by looking at what the line that contains the first @name is

var (
	PragmaName    = "dfm"
	CommentString = ""
)

type parsedPragma map[string]string

func NewFile(fileContents string) *File {
	f := &File{
		FileContents:  fileContents,
		PragmaName:    PragmaName,
		CommentString: CommentString,
	}

	return f
}

func ProcessFile(file string, vars map[string]string) error {
	fileStat, err := os.Stat(file)
	if err != nil {
		return err
	}

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	p := NewFile(string(b))
	p.Vars = vars

	s, err := p.Process()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(file, []byte(s), fileStat.Mode())
}

// File represents a file with pragmas that can be processed
type File struct {
	FileContents string

	PragmaName    string
	CommentString string

	Vars map[string]string

	pragmaLineRegex *regexp.Regexp

	hostname string
	os       string
}

func (f *File) Process() (string, error) {
	err := f.setupFileForProcessing()
	if err != nil {
		return "", err
	}

	lines := strings.Split(f.FileContents, "\n")

	// whether or not to comment next line
	commentNextLine := false
	uncommentNextLine := false
	// whether or not we are current in a comment block
	commentBlock := false

	for i, line := range lines {
		if ok, pragma := f.getPragmaForLine(line); ok {
			nextLine, blockStart, blockEnd := f.processPragma(pragma)

			commentNextLine = nextLine
			uncommentNextLine = !nextLine

			if blockStart {
				commentBlock = true
			} else if blockEnd {
				commentBlock = false

				commentNextLine = false
				uncommentNextLine = false
			}

			// set comment to current line comment if unset
			if f.CommentString == "" {
				f.CommentString = strings.Fields(line)[0]
			}

			continue
		}

		// fmt.Println(commentNextLine, uncommentNextLine, commentBlock)

		if commentNextLine {
			lines[i] = f.comment(line)

			if !commentBlock {
				commentNextLine = false
			}

			continue
		} else if uncommentNextLine {
			lines[i] = f.unComment(line)

			if !commentBlock {
				uncommentNextLine = false
			}

			continue
		}

		// if commentBlock {
		// 	lines[i] = f.comment(line)
		// }
	}

	return strings.Join(lines, "\n"), nil
}

// this isnt working if the comment is the start of the line
func (f *File) unComment(line string) string {
	strippedLine := strings.TrimLeftFunc(line, unicode.IsSpace)
	padding := len(line) - len(strippedLine)

	if strings.HasPrefix(strippedLine, f.CommentString) {
		// remove comment by removing the first x characters and a space where x is the length of the comment
		commentLength := len(f.CommentString) + 1
		line = line[:padding] + strippedLine[commentLength:]
	}

	return line
}

func (f *File) comment(line string) string {
	strippedLine := strings.TrimLeftFunc(line, unicode.IsSpace)
	padding := len(line) - len(strippedLine)

	if !strings.HasPrefix(strippedLine, f.CommentString) {
		line = fmt.Sprintf("%s%s %v", line[:padding], f.CommentString, strippedLine)
	}

	return line
}

func (f *File) setupFileForProcessing() error {
	if f.os == "" {
		f.os = runtime.GOOS
	}

	if f.hostname == "" {
		h, err := os.Hostname()
		if err != nil {
			return err
		}

		f.hostname = h
	}

	r, err := regexp.Compile(`^\W+ @` + f.PragmaName + ` (.+)@?`)
	if err != nil {
		return err
	}

	f.pragmaLineRegex = r

	return nil
}

func (f *File) getPragmaForLine(line string) (bool, parsedPragma) {
	matches := f.pragmaLineRegex.FindStringSubmatch(line)

	if len(matches) == 0 {
		return false, nil
	}

	pragmaMap := make(parsedPragma)

	pragmaParts := strings.Fields(matches[0])

	// skip the first 2 pragma parts since it'll be comment string and @name so we don't care about those
	for _, p := range pragmaParts[2:] {
		// split by the first equal sign
		i := strings.Index(p, "=")

		if i < 0 {
			pragmaMap[p] = ""
		} else {
			pragmaMap[p[:i]] = p[i+1:]
		}
	}

	return true, pragmaMap
}

// processPragma process a parsedPragma and returns if the next line or block should be commented
func (f *File) processPragma(pragma parsedPragma) (commentLine bool, commentBlockStart bool, commentBlockEnd bool) {
	pragmaParsed := false
	commentLine = true
	commentBlockStart = false
	commentBlockEnd = false

	fmt.Printf("%#v", pragma)

	for k, v := range pragma {
		switch strings.ToLower(k) {
		case "start":
			commentBlockStart = true

		case "end":
			commentBlockEnd = true

		case "host":
			if strings.EqualFold(f.hostname, v) {
				commentLine = false
			}

			pragmaParsed = true

		case "env":
			envParts := strings.Split(v, "=")
			if len(envParts) != 2 {
				printer.Warning("failed to get parse env pragma %v", v)
				continue
			}

			if strings.EqualFold(os.Getenv(envParts[0]), envParts[1]) {
				commentLine = false
			}

			pragmaParsed = true

		case "var":
			parts := strings.Split(v, "=")
			if len(parts) != 2 {
				printer.Warning("failed to get parse var pragma %v", v)
				continue
			}

			if strings.EqualFold(f.Vars[parts[0]], parts[1]) {
				commentLine = false
			}

			pragmaParsed = true

		case "os":
			if strings.EqualFold(f.os, v) {
				commentLine = false
			}

			pragmaParsed = true

		default:
			printer.Warning("Unknown pragam found: %v=%v", k, v)
		}
	}

	// if nothing was parsed return false for both
	if !pragmaParsed {
		return false, commentBlockStart, commentBlockEnd
	}

	// only enable commentBlock if commentLine is false, aka the other pragmas in the line were true
	// commentBlockStart = commentBlockStart && commentLine
	// commentBlockEnd = commentBlockEnd && commentLine

	if commentBlockEnd {
		commentLine = false
	}

	return commentLine, commentBlockStart, commentBlockEnd
}
