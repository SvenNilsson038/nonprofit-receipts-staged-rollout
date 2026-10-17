package main

import (
	"fmt"
	"log"
	"os"

	"nonprofit-rollout/infrai"
)

func main() {
	c, err := infrai.NewFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	enabled, err := c.IsEnabled("nonprofit-receipts")
	if err != nil {
		log.Fatal(err)
	}
	donations := []Donation{{Donor: "Asha", AmountCents: 5000, Campaign: "winter"}, {Donor: "Luis", AmountCents: 2500, Campaign: "winter"}}
	for _, d := range donations {
		receipt, reminder := prepare(d, enabled)
		fmt.Println(receipt)
		if reminder != nil {
			fmt.Printf("reminder queued for %s\n", reminder.Donor)
		}
	}
	job, err := c.CreateCron("0 9 * * 1", os.Getenv("PUBLIC_REMINDER_URL"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("weekly reminder job: %s\nreport: %+v\n", job, report(donations))
}
