package classtypes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/fatih/color"
	"github.com/tazjin/hiyoga/util"
	"github.com/urfave/cli"
)

const CLASSTYPES_URL string = "https://www.hiyoga.no/sats-api/no/classtypes"

type ClassProfile struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type ClassType struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	// There are more fields in class types, but they're not relevant for a CLI application (yet).
}

type ClassTypeResponse struct {
	ClassTypes []ClassType `json:"classTypes"`
}

// Calls the HiYoga API to find all available yoga class types
func ListClassTypes() (ClassTypeResponse, error) {
	var ct ClassTypeResponse

	resp, err := http.Get(CLASSTYPES_URL)

	if err != nil {
		return ct, err
	}

	body, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &ct)

	return ct, err
}

func PrettyPrintClassTypeResponse(response *ClassTypeResponse) {
	color.Green("Available class types:\n\n")

	for _, ct := range response.ClassTypes {
		prettyPrintClassType(&ct)
	}
}

func prettyPrintClassType(classType *ClassType) {
	fmt.Printf("%s\n%s\n\n",
		color.MagentaString(classType.Name),
		strings.TrimSpace(classType.Description))
}

func ListClassTypesCommand() cli.Command {
	return cli.Command{
		Name:    "list-class-types",
		Usage:   "list available yoga class types",
		Aliases: []string{"lct"},
		Action: func(c *cli.Context) error {
			ct, err := ListClassTypes()

			if err != nil {
				util.Fail(err)
			}

			PrettyPrintClassTypeResponse(&ct)
			return nil
		},
	}
}
