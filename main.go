package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
	"strconv"
)

const CLASSES_URL string = "https://www.hiyoga.no/sats-api/no/classes?regions=%s&dates=%s"
const MAJORSTUEN string = "94c207f7-fdc0-4de2-8ca4-aa42e8387b60"
const URL_DATE_FORMAT string = "20060102"
const PRETTY_PRINT_DATE_FORMAT = "Monday at 15:04"

type classResponse struct {
	Classes     []class `json:"classes"`
	UserId      string  `json:"userId"`
	ContainsAll bool    `json:"resultContainsAll"`
}

type class struct {
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

func main() {
	days := getDaysArg()

	if days > 7 {
		fail("Can not print more than one week in advance!")
	}

	classes, err := listClasses(days)

	if err != nil {
		fail(fmt.Sprintf("Could not list classes: %v\n", err))
	}

	prettyPrintClasses(days, &classes)
}

func getDaysArg() int {
	args := os.Args[1:]

	if len(args) == 1 {
		days, err := strconv.Atoi(args[0])

		if err != nil {
			fail("Usage: hiyoga [number of days]")
		}

		return days
	}

	return 3
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

func listClasses(daysFromNow int) (classResponse, error) {
	url := fmt.Sprintf(CLASSES_URL, MAJORSTUEN, nextNDays(daysFromNow))

	var c classResponse

	resp, err := http.Get(url)

	if err != nil {
		return c, err
	}

	body, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &c)

	return c, err
}

func prettyPrintClasses(dayCount int, response *classResponse) {
	color.Green("Classes at %s for the next %d days (including today):\n\n", "HiYoga Majorstuen", dayCount)

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

func prettyPrintAttendance(c *class) string {
	available := c.MaxPersonsCount - c.BookedPersonsCount

	if available > 10 {
		return color.GreenString("%d of %d spots taken", c.BookedPersonsCount, c.MaxPersonsCount)
	}

	if available == 0 {
		return color.RedString("all %d spots taken (%d in queue)", c.MaxPersonsCount, c.WaitingListCount)
	}

	return color.YellowString("%d of %d spots taken", c.BookedPersonsCount, c.MaxPersonsCount)
}

func fail(msg string) {
	fmt.Fprintf(os.Stderr, color.RedString(msg))
	os.Exit(1)
}
