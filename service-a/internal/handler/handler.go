package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Eliezer2000/weather-system/service-a/internal/model"
	"github.com/Eliezer2000/weather-system/service-a/internal/service"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type WeatherHandler struct {
	service *service.WeatherService
}

func NewWeatherHandler(service *service.WeatherService) *WeatherHandler {
	return &WeatherHandler{service: service}
}

func (h *WeatherHandler) PostCEP(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("service-a")
	ctx, span := tracer.Start(r.Context(), "PostCEP")
	defer span.End()

	var req model.CEPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"message": "invalid zipcode"}`, http.StatusUnprocessableEntity)
		span.SetAttributes(attribute.String("error", "invalid JSON"))
		return
	}
	if !h.service.IsValidCEP(req.CEP) {
		http.Error(w, `{"message": "invalid zipcode"}`, http.StatusUnprocessableEntity)
		span.SetAttributes(attribute.String("error", "invalid zipcode"))
		return
	}
	span.SetAttributes(attribute.String("cep", req.CEP))
	response, err := h.service.ForwardToServiceB(ctx, req.CEP)
	if err != nil {
		// Check if it's a specific error from Service B
		if err.Error() == "can not find zipcode" {
			http.Error(w, `{"message": "can not find zipcode"}`, http.StatusNotFound)
			span.SetAttributes(attribute.String("error", err.Error()))
			return
		}
		http.Error(w, `{"message": "error communicating with Service B"}`, http.StatusInternalServerError)
		span.SetAttributes(attribute.String("error", err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
