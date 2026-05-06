package handlers

import (
	"encoding/json"
	"fmt"
	"time"

	subdomain "github.com/IwantHappiness/subscriptions/internal/domain/subscription"
	"github.com/google/uuid"
)

const subscriptionDateLayout = "01-2006"

type MonthYear struct {
	time.Time
}

func parseMonthYear(raw string) (time.Time, error) {
	parsed, err := time.Parse(subscriptionDateLayout, raw)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format, expected MM-YYYY: %w", err)
	}

	return parsed, nil
}

func (m *MonthYear) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("date must be a string: %w", err)
	}

	parsed, err := parseMonthYear(raw)
	if err != nil {
		return err
	}

	m.Time = parsed
	return nil
}

func (m MonthYear) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Format(subscriptionDateLayout))
}

type NullMonthYear struct {
	Time *time.Time
}

func (n *NullMonthYear) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Time = nil
		return nil
	}

	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("date must be a string or null: %w", err)
	}

	parsed, err := parseMonthYear(raw)
	if err != nil {
		return err
	}

	n.Time = &parsed
	return nil
}

func (n NullMonthYear) MarshalJSON() ([]byte, error) {
	if n.Time == nil {
		return []byte("null"), nil
	}

	return json.Marshal(n.Time.Format(subscriptionDateLayout))
}

type SubMutationDTO struct {
	ServiceName string        `json:"service_name"`
	Price       int           `json:"price"`
	UserID      uuid.UUID     `json:"user_id"`
	StartDate   MonthYear     `json:"start_date"`
	EndDate     NullMonthYear `json:"end_date"`
}

type SubscriptionDTO struct {
	ID          int64         `json:"id"`
	ServiceName string        `json:"service_name"`
	Price       int           `json:"price"`
	UserID      uuid.UUID     `json:"user_id"`
	StartDate   MonthYear     `json:"start_date"`
	EndDate     NullMonthYear `json:"end_date"`
}

func newSubscriptionDTO(sub *subdomain.Subscription) SubscriptionDTO {
	dto := SubscriptionDTO{
		ID:          sub.ID,
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID,
		StartDate:   MonthYear{Time: sub.StartDate},
	}

	if sub.EndDate != nil {
		dto.EndDate = NullMonthYear{Time: sub.EndDate}
	}

	return dto
}

type TotalPriceDTO struct {
	TotalPrice  int               `json:"total_price"`
	UserID      uuid.UUID         `json:"user_id"`
	ServiceName *string           `json:"service_name"`
	List        []SubscriptionDTO `json:"list"`
}

func NewTotalPriceDTO(sub *subdomain.TotalPriceSubscription) TotalPriceDTO {
	list := make([]SubscriptionDTO, 0, len(sub.Subscription))
	for _, s := range sub.Subscription {
		list = append(list, newSubscriptionDTO(&s))
	}

	return TotalPriceDTO{
		TotalPrice:  sub.TotalPrice,
		UserID:      sub.UserID,
		ServiceName: &sub.ServiceName,
		List:        list,
	}
}
