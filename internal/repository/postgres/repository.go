package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	subdomain "github.com/IwantHappiness/subscriptions/internal/domain/subscription"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, sub *subdomain.Subscription) (*subdomain.Subscription, error) {
	const query = `
	INSERT INTO subscriptions (user_id, service_name, price, start_date, end_date)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, user_id, service_name, price, start_date, end_date`

	row := r.pool.QueryRow(ctx, query, sub.UserID, sub.ServiceName, sub.Price, sub.StartDate, sub.EndDate)
	created, err := scanSubscription(row)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (r *Repository) GetById(ctx context.Context, id int64) (*subdomain.Subscription, error) {
	const query = `
	SELECT id, user_id, service_name, price, start_date, end_date
	FROM subscriptions WHERE id = $1`

	row := r.pool.QueryRow(ctx, query, id)
	sub, err := scanSubscription(row)
	if err != nil {
		return nil, err
	}

	return sub, nil
}

func (r *Repository) Update(ctx context.Context, sub *subdomain.Subscription) (*subdomain.Subscription, error) {
	const query = `
	UPDATE subscriptions
	SET service_name = $1, price = $2, start_date = $3, end_date = $4
	WHERE id = $5
	RETURNING id, user_id, service_name, price, start_date, end_date`

	row := r.pool.QueryRow(ctx, query, sub.ServiceName, sub.Price, sub.StartDate, sub.EndDate, sub.ID)
	updated, err := scanSubscription(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, subdomain.ErrNotFound
		}
		return nil, err
	}

	return updated, nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	const query = `DELETE FROM subscriptions WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return subdomain.ErrNotFound
	}

	return nil
}

func (r *Repository) List(ctx context.Context, input subdomain.ListFilter) ([]subdomain.Subscription, error) {
	const query = `
	SELECT id, user_id, service_name, price, start_date, end_date
	FROM subscriptions ORDER BY id ASC
	LIMIT $1 OFFSET $2`

	rows, err := r.pool.Query(ctx, query, input.Limit, input.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	subscriptions := make([]subdomain.Subscription, 0, input.Limit)
	for rows.Next() {
		subscription, err := scanSubscription(rows)
		if err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, *subscription)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return subscriptions, nil
}

func (r *Repository) GetTotalPrice(ctx context.Context, input *subdomain.TotalCostFilter) (*subdomain.TotalPriceSubscription, error) {
	if input.To == nil {
		return nil, fmt.Errorf("to date is required")
	}

	const query = `
	SELECT id, user_id, service_name, price, start_date, end_date
	FROM subscriptions WHERE user_id = $1
  	AND service_name = $2
   	AND start_date <= $4
   	AND (end_date IS NULL OR end_date >= $3)`

	rows, err := r.pool.Query(ctx, query, input.UserID, input.ServiceName, input.From, *input.To)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	subscriptions := make([]subdomain.Subscription, 0, 50)
	for rows.Next() {
		subscription, err := scanSubscription(rows)
		if err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, *subscription)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	totalPrice := 0
	for _, sub := range subscriptions {
		totalPrice += totalPriceForPeriod(sub, input.From, *input.To)
	}

	return &subdomain.TotalPriceSubscription{
		UserID:       input.UserID,
		ServiceName:  input.ServiceName,
		TotalPrice:   totalPrice,
		Subscription: subscriptions,
	}, nil
}

type SubcriptionsScanner interface {
	Scan(dest ...any) error
}

func scanSubscription(scanner SubcriptionsScanner) (*subdomain.Subscription, error) {
	var subscribe subdomain.Subscription

	if err := scanner.Scan(
		&subscribe.ID,
		&subscribe.UserID,
		&subscribe.ServiceName,
		&subscribe.Price,
		&subscribe.StartDate,
		&subscribe.EndDate,
	); err != nil {
		return nil, err
	}

	return &subscribe, nil
}

func totalPriceForPeriod(sub subdomain.Subscription, from, to time.Time) int {
	periodStart := firstDayOfMonth(from)
	periodEnd := firstDayOfMonth(to)
	subStart := firstDayOfMonth(sub.StartDate)
	subEnd := periodEnd
	if sub.EndDate != nil {
		subEnd = firstDayOfMonth(*sub.EndDate)
	}

	overlapStart := maxMonth(periodStart, subStart)
	overlapEnd := minMonth(periodEnd, subEnd)
	if overlapStart.After(overlapEnd) {
		return 0
	}

	return sub.Price * monthsInclusive(overlapStart, overlapEnd)
}

func firstDayOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
}

func maxMonth(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

func minMonth(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}

func monthsInclusive(from, to time.Time) int {
	return (to.Year()-from.Year())*12 + int(to.Month()-from.Month()) + 1
}
