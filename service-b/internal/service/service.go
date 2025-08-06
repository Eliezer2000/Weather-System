package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Eliezer2000/weather-system/service-b/internal/config"
	"github.com/Eliezer2000/weather-system/service-b/internal/model"
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

func (s *WeatherService) GetWeather(ctx context.Context, cep string) (*model.WeatherResponse, error) {
	tracer := otel.Tracer("service-b")
	ctx, span := tracer.Start(ctx, "GetWeatherService")
	defer span.End()

	ctx, viaSpan := tracer.Start(ctx, "ViaCEPRequest")
	viaCEPURL := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
	resp, err := http.Get(viaCEPURL)
	if err != nil {
		viaSpan.SetAttributes(attribute.String("error", err.Error()))
		viaSpan.End()
		return nil, err
	}
	defer resp.Body.Close()
	viaSpan.End()

	var viaCEP model.ViaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&viaCEP); err != nil {
		span.SetAttributes(attribute.String("error", err.Error()))
		return nil, err
	}

	if viaCEP.Erro != "" {
		span.SetAttributes(attribute.String("error", "can not find zipcode"))
		return nil, fmt.Errorf("can not find zipcode")
	}

	ctx, weatherSpan := tracer.Start(ctx, "WeatherAPIRequest")
	encodedCity := url.QueryEscape(viaCEP.Localidade)
	weatherURL := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s", s.cfg.WeatherAPIKey, encodedCity)
	resp, err = http.Get(weatherURL)
	if err != nil {
		weatherSpan.SetAttributes(attribute.String("error", err.Error()))
		weatherSpan.End()
		return nil, err
	}
	defer resp.Body.Close()
	weatherSpan.End()

	var weatherAPI model.WeatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherAPI); err != nil {
		span.SetAttributes(attribute.String("error", err.Error()))
		return nil, err
	}

	tempC := weatherAPI.Current.TempC
	tempF := (tempC * 1.8) + 32
	tempK := tempC + 273

	return &model.WeatherResponse{
		City:  viaCEP.Localidade,
		TempC: tempC,
		TempF: tempF,
		TempK: tempK,
	}, nil
}
