package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	subdomain "github.com/IwantHappiness/subscriptions/internal/domain/subscription"
	subUseCase "github.com/IwantHappiness/subscriptions/internal/usecase/subscription"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type SubHandler struct {
	usecase subUseCase.Usecase
	logger  *slog.Logger
}

func NewSubHandler(usecase subUseCase.Usecase, logger *slog.Logger) *SubHandler {
	if logger == nil {
		logger = slog.Default()
	}

	return &SubHandler{
		usecase: usecase,
		logger:  logger,
	}
}

func (h *SubHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req SubMutationDTO
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	sub, err := h.usecase.Create(r.Context(), subUseCase.CreateInput{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   req.StartDate.Time,
		EndDate:     req.EndDate.Time,
	})
	if err != nil {
		h.writeUsecaseError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, newSubscriptionDTO(sub))
}

func (s *SubHandler) GetById(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	sub, err := s.usecase.GetById(r.Context(), id)
	if err != nil {
		s.writeUsecaseError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, newSubscriptionDTO(sub))
}

func (s *SubHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var req SubMutationDTO
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	updated, err := s.usecase.Update(r.Context(), id, subUseCase.UpdateInput{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		StartDate:   req.StartDate.Time,
		EndDate:     req.EndDate.Time,
	})
	if err != nil {
		s.writeUsecaseError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, newSubscriptionDTO(updated))

}

func (s *SubHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if err := s.usecase.Delete(r.Context(), id); err != nil {
		s.writeUsecaseError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *SubHandler) List(w http.ResponseWriter, r *http.Request) {
	subscriptions, err := s.usecase.List(r.Context())
	if err != nil {
		s.writeUsecaseError(w, err)
		return
	}

	response := make([]SubscriptionDTO, 0, len(subscriptions))
	for i, _ := range subscriptions {
		response = append(response, newSubscriptionDTO(&subscriptions[i]))
	}

	writeJSON(w, http.StatusOK, response)
}

func (s *SubHandler) GetTotalPrice(w http.ResponseWriter, r *http.Request) {
	priceFilter, err := getTotalPriceInputFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	totalPrice, err := s.usecase.GetTotalPrice(r.Context(), priceFilter)
	if err != nil {
		s.writeUsecaseError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, NewTotalPriceDTO(totalPrice))
}

func getTotalPriceInputFromRequest(r *http.Request) (subUseCase.GetTotalPriceInput, error) {
	query := r.URL.Query()

	rawUserID := strings.TrimSpace(query.Get("user_id"))
	if rawUserID == "" {
		return subUseCase.GetTotalPriceInput{}, errors.New("missing user_id")
	}

	userID, err := uuid.Parse(rawUserID)
	if err != nil {
		return subUseCase.GetTotalPriceInput{}, errors.New("invalid user_id")
	}

	serviceName := strings.TrimSpace(query.Get("service_name"))
	if serviceName == "" {
		return subUseCase.GetTotalPriceInput{}, errors.New("missing service_name")
	}

	rawFrom := strings.TrimSpace(query.Get("from"))
	if rawFrom == "" {
		return subUseCase.GetTotalPriceInput{}, errors.New("missing from")
	}

	from, err := parseMonthYear(rawFrom)
	if err != nil {
		return subUseCase.GetTotalPriceInput{}, err
	}

	rawTo := strings.TrimSpace(query.Get("to"))
	if rawTo == "" {
		return subUseCase.GetTotalPriceInput{}, errors.New("missing to")
	}

	to, err := parseMonthYear(rawTo)
	if err != nil {
		return subUseCase.GetTotalPriceInput{}, err
	}

	return subUseCase.GetTotalPriceInput{
		UserID:      userID,
		ServiceName: serviceName,
		From:        from,
		To:          &to,
	}, nil
}

func getIDFromRequest(r *http.Request) (int64, error) {
	rawID := mux.Vars(r)["id"]
	if rawID == "" {
		return 0, errors.New("missing subscription id")
	}

	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		return 0, errors.New("invalid subscription id")
	}

	if id <= 0 {
		return 0, errors.New("invalid subscription id")
	}

	return id, nil
}

func decodeJSON(r *http.Request, v any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(v); err != nil {
		return err
	}

	return nil
}

func (s *SubHandler) writeUsecaseError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, subdomain.ErrNotFound):
		writeError(w, http.StatusNotFound, err)
	case errors.Is(err, subUseCase.ErrInvalidInput):
		writeError(w, http.StatusBadRequest, err)
	default:
		s.logger.Error("subscription usecase failed", "error", err)
		writeError(w, http.StatusInternalServerError, err)
	}
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(payload)
}
