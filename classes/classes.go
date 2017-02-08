package classes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/tazjin/hiyoga/util"
	"github.com/urfave/cli"
	"os"
	"sync"
	"text/tabwriter"
)

const CLASSES_URL string = "https://www.hiyoga.no/sats-api/no/classes?regions=%s&dates=%s"
const URL_DATE_FORMAT string = "20060102"

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

	resp, err := util.HiyogaGet(url)

	if err != nil {
		return c, err
	}

	err = json.Unmarshal(resp, &c)

	return c, err
}

func PrettyPrintClassResponse(dayCount int, response *ClassResponse) {
	color.White("Classes at %s for the next %d days (including today):\n\n", "HiYoga Majorstuen", dayCount)

	// Prevent extra newline from being printed if no classes happen today by flipping this after the first class
	var once sync.Once
	firstOut := false

	day := time.Now().Day()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	for _, c := range response.Classes {
		// Extra newline between days
		if c.StartTime.Day() != day {
			if firstOut {
				fmt.Fprintln(w, "")
			}
			day = c.StartTime.Day()
		}

		once.Do(func() { firstOut = true })

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t\n",
			PrettyPrintClassTime(&c),
			color.MagentaString(c.Name),
			c.InstructorId,
			prettyPrintAttendance(&c))
	}

	w.Flush()
}

func PrettyPrintClassTime(c *Class) string {
	endTime := c.StartTime.Add(time.Duration(c.DurationInMinutes) * time.Minute)

	return color.BlueString(
		"%s-%s",
		c.StartTime.Format("Monday\t15:04"),
		endTime.Format("15:04"),
	)
}

func prettyPrintAttendance(c *Class) string {
	available := c.MaxPersonsCount - c.BookedPersonsCount
	taken := fmt.Sprintf("%2d of %2d spots taken", c.BookedPersonsCount, c.MaxPersonsCount)

	if available > 10 {
		return color.GreenString(taken)
	}

	if available == 0 {
		return color.RedString("%s (%d in queue)", taken, c.WaitingListCount)
	}

	return color.YellowString(taken)
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
