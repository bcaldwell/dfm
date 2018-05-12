package shMock

import (
	"strings"

	"github.com/benjamincaldwell/go-sh"
)

// list of commands that go run
var Commands []string

// top element is popped from this array and returned as the error value
var ErrorsToReturn []error

// top element is popped and returned from output function
var OutputToRetrun []string

var DefaultShInterface sh.SessionInterface
var MockShellInterface sh.SessionInterface

func init() {
	DefaultShInterface = sh.MainInterface
	MockShellInterface = new(sessionMock)
	UseMock()
}

func Reset() {
	Commands = Commands[:0]
	ErrorsToReturn = ErrorsToReturn[:0]
	OutputToRetrun = OutputToRetrun[:0]
}

func UseDefault() {
	sh.MainInterface = DefaultShInterface
}

func UseMock() {
	sh.MainInterface = MockShellInterface
}

type sessionMock struct {
	cmd   string
	dir   string
	input string
	env   map[string]string
}

func (c *sessionMock) New() sh.SessionInterface {
	return new(sessionMock)
}

func (c *sessionMock) Command(name string, arg ...string) sh.SessionInterface {
	args := append([]string{name}, arg...)
	c.cmd = strings.Join(args, " ")
	return c
}

func (c *sessionMock) SetInput(s string) sh.SessionInterface {
	c.input = s
	return c
}

func (c *sessionMock) SetDir(s string) sh.SessionInterface {
	c.dir = s
	return c
}

func (c *sessionMock) SetPath(path string) sh.SessionInterface {
	c.env["PATH"] = path
	return c
}

func (c *sessionMock) SetEnv(key, value string) sh.SessionInterface {
	c.env[key] = value
	return c
}

func (c *sessionMock) Output() ([]byte, error) {
	Commands = append(Commands, c.cmd)

	return make([]byte, 0), errorReturnValue()
}

func (c *sessionMock) PrintOutput() error {
	return errorReturnValue()
}

func (c *sessionMock) Run() error {
	Commands = append(Commands, c.cmd)

	return errorReturnValue()
}

func errorReturnValue() error {
	var returnErr error
	if len(ErrorsToReturn) > 0 {
		returnErr = ErrorsToReturn[0]
		ErrorsToReturn = ErrorsToReturn[1:]
	}
	return returnErr
}
