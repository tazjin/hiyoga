package util

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/polydawn/meep"
	"io/ioutil"
	"net/http"
)

type HiyogaError struct {
	meep.AllTraits
}

func Fail(cause error) {
	hiyogaError := meep.New(
		&HiyogaError{},
		meep.Cause(cause),
	)
	fmt.Fprintf(os.Stderr, color.RedString("%v\n", hiyogaError))
	os.Exit(1)
}

const UserAgent string = "HiYoga CLI (https://github.com/tazjin/hiyoga)"

// GET a URL with the correct user agent set
func HiyogaGet(url string) (body []byte, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	req.Header.Set("user-agent", UserAgent)

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)

	return
}
