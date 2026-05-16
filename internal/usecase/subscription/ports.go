package subscription

import (
	"context"
	"time"

	subdomain "github.com/IwantHappiness/subscriptions/internal/domain/subscription"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, sub *subdomain.Subscription) (*subdomain.Subscription, error)
	GetById(ctx context.Context, id int64) (*subdomain.Subscription, error)
	Update(ctx context.Context, sub *subdomain.Subscription) (*subdomain.Subscription, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, input subdomain.ListFilter) ([]subdomain.Subscription, error)
	GetTotalPrice(ctx context.Context, input *subdomain.TotalCostFilter) (*subdomain.TotalPriceSubscription, error)
}

type Usecase interface {
	Create(ctx context.Context, input CreateInput) (*subdomain.Subscription, error)
	GetById(ctx context.Context, id int64) (*subdomain.Subscription, error)
	Update(ctx context.Context, id int64, input UpdateInput) (*subdomain.Subscription, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, input ListInput) ([]subdomain.Subscription, error)
	GetTotalPrice(ctx context.Context, input GetTotalPriceInput) (*subdomain.TotalPriceSubscription, error)
}

type CreateInput struct {
	ServiceName string     `json:"service_name"`
	Price       int        `json:"price"`
	UserID      uuid.UUID  `json:"user_id"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
}

type UpdateInput struct {
	ServiceName string     `json:"service_name"`
	Price       int        `json:"price"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
}

type GetTotalPriceInput struct {
	UserID      uuid.UUID  `json:"user_id"`
	ServiceName string     `json:"service_name"`
	From        time.Time  `json:"from"`
	To          *time.Time `json:"to"`
}

type ListInput struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
