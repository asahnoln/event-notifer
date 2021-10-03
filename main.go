package main

import (
	"fmt"
	"log"
	"os"

	"github.com/asahnoln/event-notifier/pkg"
	"google.golang.org/api/option"
)

// TODO: Test no just happy path - have to tell all those vars are necessary
// TODO If Tg is not set up correctly - there is no error
func main() {
	cal := pkg.NewGCalStore(os.Getenv("GCALID"), option.WithCredentialsFile(os.Getenv("GCALCRED")))
	es, err := pkg.TomorrowEvents(cal)
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
