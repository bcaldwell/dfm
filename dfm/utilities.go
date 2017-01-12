package dfm

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/benjamincaldwell/devctl/printer"
)

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func uniqueSliceTransform(a []string) (output []string) {

	for _, s := range a {
		output = appendIfUnique(output, s)
	}

	return output
}

func appendIfUnique(slice []string, i string) []string {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}

func fatalErrorCheck(e error, s string) {
	if e != nil {
		printer.Error("%s: %s", s, e)
		os.Exit(1)
	}
}

func errorCheck(err error, message string) bool {
	if err != nil {
		printer.Fail("%s failed with %s", message, err)
		return true
	}
	return false
}

// AskForConfirmation asks the user for confirmation. A user must type in "yes" or "no" and
// then press enter. It has fuzzy matching, so "y", "Y", "yes", "YES", and "Yes" all count as
// confirmations. If the input is not recognized, it will ask again. The function does not return
// until it gets a valid response from the user.
func AskForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		printer.InfoBar("%s?%s %s [y/n]: ", printer.Green+printer.Bold, printer.Nc, s)
		fmt.Printf("\033[1A\033[%dC", len(s)+12)
		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}

func HTTPDownload(uri string) ([]byte, error) {
	res, err := http.Get(uri)
	errorCheck(err, "downloading "+uri)

	defer res.Body.Close()
	d, err := ioutil.ReadAll(res.Body)
	errorCheck(err, "reading "+uri)

	return d, err
}

func WriteFile(dst string, d []byte) error {
	err := ioutil.WriteFile(dst, d, 0444)
	errorCheck(err, "writing "+dst)

	return err
}

func DownloadToFile(uri string, dst string) error {
	d, err := HTTPDownload(uri)
	if err == nil {
		return WriteFile(dst, d)
	}
	return err
}
