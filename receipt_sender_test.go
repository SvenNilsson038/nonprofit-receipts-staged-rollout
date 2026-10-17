package main

import "testing"

func TestPrepare(t *testing.T) {
	tests := []struct {
		name         string
		enabled      bool
		wantReminder bool
	}{{"staged", true, true}, {"held", false, false}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, r := prepare(Donation{Donor: "Nia", AmountCents: 100, Campaign: "spring"}, tt.enabled)
			if (r != nil) != tt.wantReminder {
				t.Fatalf("reminder=%v", r != nil)
			}
		})
	}
}

func TestReport(t *testing.T) {
	got := report([]Donation{{Campaign: "a", AmountCents: 10}, {Campaign: "a", AmountCents: 20}})
	if len(got) != 1 || got[0].TotalCents != 30 || got[0].Donors != 2 {
		t.Fatalf("unexpected report: %+v", got)
	}
}
