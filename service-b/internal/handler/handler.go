package handler

import (
    "encoding/json"
    "net/http"
    "github.com/Eliezer2000/weather-system/service-b/internal/service"
    "github.com/gorilla/mux"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
)

type WeatherHandler struct {
	service *service.WeatherService
}

func NewWeatherHandler(service *service.WeatherService) *WeatherHandler {
	return &WeatherHandler{service: service}
}

func (h *WeatherHandler) GetWeather(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("service-b")
	ctx, span := tracer.Start(r.Context(), "GetWeather")
	defer span.End()

	vars := mux.Vars(r)
	cep := vars["cep"]
	span.SetAttributes(attribute.String("cep", cep))

	if !h.service.IsValidCEP(cep) {
		http.Error(w, `{"message": "invalid zipcode"}`, http.StatusUnprocessableEntity)
		span.SetAttributes(attribute.String("error", "invalid zipcode"))
		return
	}
	
	weather, err := h.service.GetWeather(ctx, cep)
	if err != nil {
		span.SetAttributes(attribute.String("error", err.Error()))
		if err.Error() == "can not find zipcode" {
			http.Error(w, `{"message": "can not find zipcode"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"message": "internal server error"}`, http.StatusInternalServerError)
		return	
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(weather)

}