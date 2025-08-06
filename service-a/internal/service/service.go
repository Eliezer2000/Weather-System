package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Eliezer2000/weather-system/service-a/internal/config"
	"github.com/Eliezer2000/weather-system/service-a/internal/model"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type WeatherService struct {
	cfg *config.Config
}

func NewWeatherService(cfg *config.Config) *WeatherService {
	return &WeatherService{cfg: cfg}
}

func (s *WeatherService) IsValidCEP(cep string) bool {
	if len(cep) != 8 {
		return false
	}
	for _, char := range cep {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}

func (s *WeatherService) ForwardToServiceB(ctx context.Context, cep string) (map[string]interface{}, error) {
	tracer := otel.Tracer("service-a")
	_, span := tracer.Start(ctx, "ForwardToServiceB")
	defer span.End()
	
	url := s.cfg.ServiceBURL + "/weather/" + cep
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		span.SetAttributes(attribute.String("error", err.Error()))
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		span.SetAttributes(attribute.String("error", err.Error()))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResp model.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err == nil {
			span.SetAttributes(attribute.String("error", errorResp.Message))
			return nil, fmt.Errorf(errorResp.Message)
		}
		span.SetAttributes(attribute.String("error", "unexpected response from Service B"))
		return nil, fmt.Errorf("unexpected response from Service B: %d", resp.StatusCode)
	}
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		span.SetAttributes(attribute.String("error", err.Error()))
		return nil, err
	}
	return result, nil
}