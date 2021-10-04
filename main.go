package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/asahnoln/event-notifier/pkg"
	"google.golang.org/api/option"
)

// TODO: Test no just happy path - have to tell all those vars are necessary
// TODO If Tg is not set up correctly - there is no error
func main() {
	tomorrow := flag.Bool("tomorrow", false, "Use this flag if you need tomorrow events")
	flag.Parse()

	cal := pkg.NewGCalStore(os.Getenv("GCALID"), option.WithCredentialsFile(os.Getenv("GCALCRED")))

	var (
		es  []pkg.Event
		err error
	)
	if *tomorrow {
		es, err = pkg.TomorrowEvents(cal)
	} else {
		es, err = pkg.TodayEvents(cal)
	}
	if err != nil {
		log.Fatal(err)
	}

	sdr := pkg.NewTg(os.Getenv("TGKEY"), os.Getenv("TGCHATID"))
	err = pkg.Send(es, sdr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Sent succesfully!")
}
