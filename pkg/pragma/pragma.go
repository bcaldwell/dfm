package pragma

import (
	"os"
	"regexp"
	"strings"

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

// File represents a file with pragmas that can be processed
type File struct {
	FileContents string

	PragmaName    string
	CommentString string

	pragmaLineRegex *regexp.Regexp

	hostname string
	os       string
}

func (f *File) Process() error {
	return nil
}

func (f *File) generateRegex() error {
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

	// skip the first pragma part since itll be @name so we dont care about that one
	for _, p := range pragmaParts[1:] {
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

	for k, v := range pragma {
		switch strings.ToLower(k) {
		case "start":
			commentBlockStart = true
			pragmaParsed = true

		case "end":
			commentBlockEnd = true
			pragmaParsed = true

		case "host":
			if f.hostname != v {
				commentLine = false
			}

			pragmaParsed = true

		case "env":
			envParts := strings.Split(v, "=")
			if len(envParts) != 2 {
				printer.Warning("failed to get parse env pragma %v", v)
				continue
			}

			if os.Getenv(envParts[0]) != envParts[1] {
				commentLine = false
			}

			pragmaParsed = true

		case "os":
			if f.os != v {
				commentLine = false
			}

			pragmaParsed = true

		default:
			printer.Warning("Unknown pragam found: %v=k", v, k)
		}
	}

	// if nothing was parsed return false for both
	if !pragmaParsed {
		return false, false, false
	}

	// only enable commentBlock if commentLine is true, aka the other pragmas in the line were true
	commentBlockStart = commentBlockStart && commentLine
	commentBlockEnd = commentBlockEnd && commentLine

	if commentBlockEnd {
		commentLine = false
	}

	return commentLine, commentBlockStart, commentBlockEnd
}
