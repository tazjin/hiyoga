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
	"os"
	"text/tabwriter"
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
	color.White("Booked classes:\n")

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	for _, b := range list {
		fmt.Fprintf(w, "%s\t%s\t%s\t(%s)\t\n",
			classes.PrettyPrintClassTime(&b.Class.StartTime),
			color.MagentaString(b.Class.Name),
			b.Class.InstructorId,
			prettyPrintBookingStatus(&b),
		)
	}

	w.Flush()
}

func prettyPrintBookingStatus(b *Booking) string {
	if b.Status == "confirmed" {
		return color.GreenString(b.Status)
	}

	if b.PositionInQueue >= 0 {
		return color.YellowString("%d in queue", b.PositionInQueue)
	}

	return color.YellowString(b.Status)
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
