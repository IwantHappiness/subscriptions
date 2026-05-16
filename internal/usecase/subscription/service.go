package subscription

import (
	"context"
	"fmt"
	"strings"

	subdomain "github.com/IwantHappiness/subscriptions/internal/domain/subscription"
	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (*subdomain.Subscription, error) {
	validInput, err := isValidCreateInput(input)
	if err != nil {
		return nil, err
	}

	model := &subdomain.Subscription{
		ServiceName: validInput.ServiceName,
		Price:       validInput.Price,
		UserID:      validInput.UserID,
		StartDate:   validInput.StartDate,
		EndDate:     validInput.EndDate,
	}

	created, err := s.repo.Create(ctx, model)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *Service) GetById(ctx context.Context, id int64) (*subdomain.Subscription, error) {
	subscription, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	return subscription, nil
}

func (s *Service) Update(ctx context.Context, id int64, input UpdateInput) (*subdomain.Subscription, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	validUpdateInput, err := isValidUpdateInput(input)
	if err != nil {
		return nil, err
	}

	model := &subdomain.Subscription{
		ID:          id,
		ServiceName: validUpdateInput.ServiceName,
		Price:       validUpdateInput.Price,
		StartDate:   validUpdateInput.StartDate,
		EndDate:     validUpdateInput.EndDate,
	}

	updated, err := s.repo.Update(ctx, model)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

func (s *Service) List(ctx context.Context, input ListInput) ([]subdomain.Subscription, error) {
	validInput, err := isValidListInput(input)
	if err != nil {
		return nil, err
	}

	filter := subdomain.ListFilter{
		Limit:  validInput.Limit,
		Offset: validInput.Offset,
	}

	subscriptions, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	return subscriptions, nil
}

func (s *Service) GetTotalPrice(ctx context.Context, input GetTotalPriceInput) (*subdomain.TotalPriceSubscription, error) {
	validInput, err := isValidGetTotalPriceInput(input)
	if err != nil {
		return nil, err
	}

	req := &subdomain.TotalCostFilter{
		UserID:      validInput.UserID,
		ServiceName: validInput.ServiceName,
		From:        validInput.From,
		To:          validInput.To,
	}

	result, err := s.repo.GetTotalPrice(ctx, req)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func isValidCreateInput(input CreateInput) (CreateInput, error) {
	input.ServiceName = strings.TrimSpace(input.ServiceName)

	if input.ServiceName == "" {
		return CreateInput{}, fmt.Errorf("%w: service name is required", ErrInvalidInput)
	}

	if input.Price < 0 {
		return CreateInput{}, fmt.Errorf("%w: price must be greater than or equal to 0", ErrInvalidInput)
	}

	if input.UserID == uuid.Nil {
		return CreateInput{}, fmt.Errorf("%w: user id is required", ErrInvalidInput)
	}

	if input.EndDate != nil && input.StartDate.After(*input.EndDate) {
		return CreateInput{}, fmt.Errorf("%w: end date must be after start date", ErrInvalidInput)
	}

	if input.StartDate.IsZero() {
		return CreateInput{}, fmt.Errorf("%w: start date is required", ErrInvalidInput)
	}

	return input, nil
}

func isValidGetTotalPriceInput(input GetTotalPriceInput) (GetTotalPriceInput, error) {
	input.ServiceName = strings.TrimSpace(input.ServiceName)

	if input.UserID == uuid.Nil {
		return GetTotalPriceInput{}, fmt.Errorf("%w: user id is required", ErrInvalidInput)
	}

	if input.ServiceName == "" {
		return GetTotalPriceInput{}, fmt.Errorf("%w: service name is required", ErrInvalidInput)
	}

	if input.From.IsZero() {
		return GetTotalPriceInput{}, fmt.Errorf("%w: from date is required", ErrInvalidInput)
	}

	if input.To == nil {
		return GetTotalPriceInput{}, fmt.Errorf("%w: to date is required", ErrInvalidInput)
	}

	if input.To.IsZero() {
		return GetTotalPriceInput{}, fmt.Errorf("%w: to date is required", ErrInvalidInput)
	}

	if input.From.After(*input.To) {
		return GetTotalPriceInput{}, fmt.Errorf("%w: to date must be after from date", ErrInvalidInput)
	}

	return input, nil
}

func isValidListInput(input ListInput) (ListInput, error) {
	if input.Limit <= 0 {
		return ListInput{}, fmt.Errorf("%w: limit must be positive", ErrInvalidInput)
	}

	if input.Limit > 100 {
		return ListInput{}, fmt.Errorf("%w: limit must be less than or equal to 100", ErrInvalidInput)
	}

	if input.Offset < 0 {
		return ListInput{}, fmt.Errorf("%w: offset must be greater than or equal to 0", ErrInvalidInput)
	}

	return input, nil
}

func isValidUpdateInput(input UpdateInput) (UpdateInput, error) {
	input.ServiceName = strings.TrimSpace(input.ServiceName)

	if input.ServiceName == "" {
		return UpdateInput{}, fmt.Errorf("%w: service name is required", ErrInvalidInput)
	}

	if input.Price < 0 {
		return UpdateInput{}, fmt.Errorf("%w: price must be greater than or equal to 0", ErrInvalidInput)
	}

	if input.StartDate.IsZero() {
		return UpdateInput{}, fmt.Errorf("%w: start date is required", ErrInvalidInput)
	}

	if input.EndDate != nil && input.StartDate.After(*input.EndDate) {
		return UpdateInput{}, fmt.Errorf("%w: end date must be after start date", ErrInvalidInput)
	}

	return input, nil
}
