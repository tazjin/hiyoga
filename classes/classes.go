package classes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/fatih/color"
	"github.com/tazjin/hiyoga/util"
	"github.com/urfave/cli"
)

const CLASSES_URL string = "https://www.hiyoga.no/sats-api/no/classes?regions=%s&dates=%s"
const URL_DATE_FORMAT string = "20060102"
const PRETTY_PRINT_DATE_FORMAT = "Monday at 15:04"

// Right now HiYoga Majorstuen is the only center
const MAJORSTUEN string = "94c207f7-fdc0-4de2-8ca4-aa42e8387b60"

type Class struct {
	Name               string    `json:"name"`
	Id                 string    `json:"id"`
	ClassTypeId        string    `json:"classTypeId"`
	CenterFilterId     string    `json:"centerFilterId"`
	RegionId           string    `json:"regionId"`
	InstructorId       string    `json:"instructorId"`
	StartTime          time.Time `json:"startTime"`
	DurationInMinutes  int       `json:"durationInMinutes"`
	BookedPersonsCount int       `json:"bookedPersonsCount"`
	MaxPersonsCount    int       `json:"maxPersonsCount"`
	WaitingListCount   int       `json:"waitingListCount"`
	ClassCategoryIds   []string  `json:"classCategoryIds"`
}

type ClassResponse struct {
	Classes     []Class `json:"classes"`
	UserId      string  `json:"userId"`
	ContainsAll bool    `json:"resultContainsAll"`
}

func nextNDays(count int) string {
	today := time.Now()
	var buffer bytes.Buffer

	for i := 0; i < count; i++ {
		next := today.AddDate(0, 0, i)
		buffer.WriteString(next.Format(URL_DATE_FORMAT))

		// Intercalate semicolons (no trailing semicolon)
		if i != (count - 1) {
			buffer.WriteString(";")
		}
	}

	return buffer.String()
}

// Calls the HiYoga API to find all scheduled classes within the specified timeframe
func ListClasses(center string, daysFromNow int) (ClassResponse, error) {
	url := fmt.Sprintf(CLASSES_URL, center, nextNDays(daysFromNow))

	var c ClassResponse

	resp, err := http.Get(url)

	if err != nil {
		return c, err
	}

	body, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &c)

	return c, err
}

func PrettyPrintClassResponse(dayCount int, response *ClassResponse) {
	color.White("Classes at %s for the next %d days (including today):\n\n", "HiYoga Majorstuen", dayCount)

	day := time.Now().Day()

	for _, c := range response.Classes {
		// Extra newline between days
		if c.StartTime.Day() != day {
			fmt.Println("")
			day = c.StartTime.Day()
		}

		fmt.Printf("%s: %s with %s (%s)\n",
			prettyPrintClassTime(&c.StartTime),
			color.MagentaString(c.Name),
			c.InstructorId,
			prettyPrintAttendance(&c))

	}
}

func prettyPrintClassTime(classtime *time.Time) string {
	return color.BlueString(classtime.Format(PRETTY_PRINT_DATE_FORMAT))
}

func prettyPrintAttendance(c *Class) string {
	available := c.MaxPersonsCount - c.BookedPersonsCount

	if available > 10 {
		return color.GreenString("%d of %d spots taken", c.BookedPersonsCount, c.MaxPersonsCount)
	}

	if available == 0 {
		return color.RedString("all %d spots taken (%d in queue)", c.MaxPersonsCount, c.WaitingListCount)
	}

	return color.YellowString("%d of %d spots taken", c.BookedPersonsCount, c.MaxPersonsCount)
}

func ListClassesCommand() cli.Command {
	return cli.Command{
		Name:    "list-classes",
		Usage:   "list upcoming yoga classes",
		Aliases: []string{"lc"},
		Action: func(c *cli.Context) error {
			days := c.Int("days")
			cl, err := ListClasses(MAJORSTUEN, days)

			if err != nil {
				util.Fail(err)
			}

			PrettyPrintClassResponse(days, &cl)
			return nil
		},
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "days, d",
				Usage: "number of days to list (including today)",
				Value: 3,
			},
		},
	}
}
