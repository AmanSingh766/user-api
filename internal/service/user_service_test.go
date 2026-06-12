package service

import (
	"testing"
	"time"
)

func TestCalculateAge(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name string
		dob  time.Time
		want int
	}{
		{
			name: "birthday already passed this year",
			// Born exactly 25 years ago yesterday
			dob:  time.Date(now.Year()-25, now.Month(), now.Day()-1, 0, 0, 0, 0, time.UTC),
			want: 25,
		},
		{
			name: "birthday is today",
			dob:  time.Date(now.Year()-30, now.Month(), now.Day(), 0, 0, 0, 0, time.UTC),
			want: 30,
		},
		{
			name: "birthday not yet this year",
			// Born exactly 22 years from tomorrow
			dob:  time.Date(now.Year()-22, now.Month(), now.Day()+1, 0, 0, 0, 0, time.UTC),
			want: 21,
		},
		{
			name: "newborn (same day)",
			dob:  time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC),
			want: 0,
		},
		{
			name: "classic birthday – 1990-05-10",
			dob:  time.Date(1990, 5, 10, 0, 0, 0, 0, time.UTC),
			want: func() int {
				n := now.Year() - 1990
				if now.Month() < 5 || (now.Month() == 5 && now.Day() < 10) {
					n--
				}
				return n
			}(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := CalculateAge(tc.dob)
			if got != tc.want {
				t.Errorf("CalculateAge(%v) = %d, want %d", tc.dob, got, tc.want)
			}
		})
	}
}
