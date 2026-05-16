package subscription

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          int64      `json:"id"`
	ServiceName string     `json:"service_name"`
	Price       int        `json:"price"`
	UserID      uuid.UUID  `json:"user_id"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
}

type TotalCostFilter struct {
	UserID      uuid.UUID
	ServiceName string
	From        time.Time
	To          *time.Time
}

type ListFilter struct {
	Limit  int
	Offset int
}

type TotalPriceSubscription struct {
	UserID       uuid.UUID      `json:"user_id"`
	ServiceName  string         `json:"service_name"`
	TotalPrice   int            `json:"total_price"`
	Subscription []Subscription `json:"subscription"`
}
