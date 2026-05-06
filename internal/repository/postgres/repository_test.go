package postgres

import (
	"testing"
	"time"

	subdomain "github.com/IwantHappiness/subscriptions/internal/domain/subscription"
)

func TestTotalPriceForPeriod(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		sub  subdomain.Subscription
		from time.Time
		to   time.Time
		want int
	}{
		{
			name: "active subscription covers full selected period",
			sub: subdomain.Subscription{
				Price:     100,
				StartDate: month(2024, time.January),
			},
			from: month(2024, time.January),
			to:   month(2024, time.March),
			want: 300,
		},
		{
			name: "subscription starts inside selected period",
			sub: subdomain.Subscription{
				Price:     100,
				StartDate: month(2024, time.February),
			},
			from: month(2024, time.January),
			to:   month(2024, time.March),
			want: 200,
		},
		{
			name: "subscription ends inside selected period",
			sub: subdomain.Subscription{
				Price:     100,
				StartDate: month(2024, time.January),
				EndDate:   ptr(month(2024, time.February)),
			},
			from: month(2024, time.January),
			to:   month(2024, time.March),
			want: 200,
		},
		{
			name: "subscription does not overlap selected period",
			sub: subdomain.Subscription{
				Price:     100,
				StartDate: month(2024, time.April),
			},
			from: month(2024, time.January),
			to:   month(2024, time.March),
			want: 0,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := totalPriceForPeriod(tt.sub, tt.from, tt.to)
			if got != tt.want {
				t.Fatalf("totalPriceForPeriod() = %d, want %d", got, tt.want)
			}
		})
	}
}

func month(year int, month time.Month) time.Time {
	return time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
}

func ptr[T any](v T) *T {
	return &v
}
