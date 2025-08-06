package main

import (
	"context"
    "log"
    "net/http"
    "github.com/Eliezer2000/weather-system/service-b/internal/config"
    "github.com/Eliezer2000/weather-system/service-b/internal/handler"
    "github.com/Eliezer2000/weather-system/service-b/internal/service"
    "github.com/gorilla/mux"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/zipkin"
    "go.opentelemetry.io/otel/sdk/resource"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.17.0"

)
func initTracer() *sdktrace.TracerProvider {
	exporter, err := zipkin.New("http://zipkin:9411/api/v2/spans")
	if err != nil {
		log.Fatalf("Failed to create Zipkin exporter: %v", err)
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("service-b"),
		)),
	)
	otel.SetTracerProvider(tp)

	return tp

}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	tp := initTracer()
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	weatherService := service.NewWeatherService(cfg)
	weatherHandler := handler.NewWeatherHandler(weatherService)

	router := mux.NewRouter()
	router.HandleFunc("/weather/{cep}", weatherHandler.GetWeather).Methods("GET")

	log.Printf("Service B running on port %s", cfg.Port)
	if err := http.ListenAndServe(":" + cfg.Port, router); err != nil {
		log.Fatalf("Service B failed: %v", err)
	}
}