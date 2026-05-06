package transporthttp

import (
	"log/slog"
	"net/http"

	swaggerdocs "github.com/IwantHappiness/subscriptions/internal/transport/http/docs"
	httphandlers "github.com/IwantHappiness/subscriptions/internal/transport/http/handlers"
	"github.com/gorilla/mux"
)

func NewRouter(subHandler *httphandlers.SubHandler, docsHandler *swaggerdocs.Handler, logger *slog.Logger) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.Use(requestLogger(logger))

	router.HandleFunc("/swagger/openapi.json", docsHandler.ServeSpec).Methods(http.MethodGet)
	router.HandleFunc("/swagger/", docsHandler.ServeUI).Methods(http.MethodGet)
	router.HandleFunc("/swagger", docsHandler.RedirectToUI).Methods(http.MethodGet)

	api := router.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/subscriptions", subHandler.Create).Methods(http.MethodPost)
	api.HandleFunc("/subscriptions", subHandler.List).Methods(http.MethodGet)
	api.HandleFunc("/subscriptions/total-price", subHandler.GetTotalPrice).Methods(http.MethodGet)
	api.HandleFunc("/subscriptions/{id:[0-9]+}", subHandler.GetById).Methods(http.MethodGet)
	api.HandleFunc("/subscriptions/{id:[0-9]+}", subHandler.Update).Methods(http.MethodPut)
	api.HandleFunc("/subscriptions/{id:[0-9]+}", subHandler.Delete).Methods(http.MethodDelete)

	return router
}
