package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	_ "github.com/go-chi/chi/v5"
	"github.com/ronaldognunes/lab-desafio-01/internal/infra/service/cep"
	"github.com/ronaldognunes/lab-desafio-01/internal/infra/service/temperatura"
	web "github.com/ronaldognunes/lab-desafio-01/internal/infra/webserver"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.9.0"
)

func main() {
	tp := InitTracer()
	defer tp.Shutdown(context.Background())

	servicoCep := cep.NewCepService("https://viacep.com.br/ws/")
	servicoTemperatura := temperatura.NewTemperaturaService("http://api.weatherapi.com/v1/current.json?q=", "6c71bade554742f0b46142657250703")
	consultaHandler := web.NewConsultaHandler(servicoCep, servicoTemperatura)

	router := chi.NewRouter()
	router.Get("/consulta-temperatura-por-cep", consultaHandler.ConsultarCepHandler)

	http.ListenAndServe(":8080", otelhttp.NewHandler(router, "webapi02"))
}

func InitTracer() *trace.TracerProvider {
	exp, err := otlptracehttp.New(context.Background(),
		otlptracehttp.WithEndpoint("otel-collector:4318"),
		otlptracehttp.WithInsecure())
	if err != nil {
		log.Fatalf("Falha ao criar exportador: %v", err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("webapi02"),
		)),
	)

	otel.SetTracerProvider(tp)
	return tp
}
