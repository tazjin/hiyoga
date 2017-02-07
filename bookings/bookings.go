package bookings

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/fatih/color"
	"github.com/tazjin/hiyoga/auth"
	"github.com/tazjin/hiyoga/classes"
	"github.com/tazjin/hiyoga/util"
	"github.com/urfave/cli"
	"time"
)

type Booking struct {
	Id              string        `json:"id"`
	Status          string        `json:"status"`
	Class           classes.Class `json:"class"`
	PositionInQueue int           `json:"positionInQueue"`
	CenterId        string        `json:"centerId"`
}

type BookingListResult struct {
	Bookings []Booking `json:"bookings"`
}

const BookingUrl = "https://www.hiyoga.no/sats-api/no/bookings"

func listBookings() []Booking {
	resp, err := auth.AuthenticatedGet(BookingUrl)

	if err != nil {
		util.Fail(err)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var result BookingListResult

	err = json.Unmarshal(body, &result)

	if err != nil {
		util.Fail(err)
	}

	return result.Bookings
}

func prettyPrintBookings(list []Booking) {
	color.White("Booked classes:")

	day := time.Now().Day()

	for _, b := range list {
		// Extra newline between days
		if b.Class.StartTime.Day() != day {
			fmt.Println("")
			day = b.Class.StartTime.Day()
		}

		fmt.Printf("%s: %s with %s (%s)\n",
			b.Class.StartTime.Format(classes.PRETTY_PRINT_DATE_FORMAT),
			b.Class.Name,
			b.Class.InstructorId,
			prettyPrintBookingStatus(b.Status),
		)
	}
}

func prettyPrintBookingStatus(status string) string {
	if status == "confirmed" {
		return color.GreenString(status)
	}

	return color.YellowString(status)
}

func ListBookingsCommand() cli.Command {
	return cli.Command{
		Name:    "list-bookings",
		Usage:   "List current bookings (requires auth)",
		Aliases: []string{"lb"},
		Action: func(c *cli.Context) error {
			prettyPrintBookings(listBookings())
			return nil
		},
	}
}
