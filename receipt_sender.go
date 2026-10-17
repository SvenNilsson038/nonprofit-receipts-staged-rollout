package main

import "fmt"

type Donation struct {
	Donor       string
	AmountCents int
	Campaign    string
}
type Reminder struct {
	Donor    string
	Campaign string
}
type Report struct {
	Campaign   string
	TotalCents int
	Donors     int
}

func prepare(d Donation, enabled bool) (string, *Reminder) {
	if !enabled {
		return "legacy receipt for " + d.Donor, nil
	}
	return fmt.Sprintf("receipt for %s: %d cents", d.Donor, d.AmountCents), &Reminder{Donor: d.Donor, Campaign: d.Campaign}
}

func report(ds []Donation) []Report {
	by := map[string]Report{}
	for _, d := range ds {
		r := by[d.Campaign]
		r.Campaign = d.Campaign
		r.TotalCents += d.AmountCents
		r.Donors++
		by[d.Campaign] = r
	}
	out := make([]Report, 0, len(by))
	for _, r := range by {
		out = append(out, r)
	}
	return out
}
