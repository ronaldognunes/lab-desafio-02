package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

func main() {
	tp := InitTracer()
	defer tp.Shutdown(context.Background())
	router := chi.NewRouter()
	router.Get("/consulta-cep-servidor/{cep}", consultarServidor)
	http.ListenAndServe(":8081", otelhttp.NewHandler(router, "webapi01"))
}

type ResponseDto struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_c"`
	TempF float64 `json:"temp_f"`
	TempK float64 `json:"temp_k"`
}

func consultarServidor(w http.ResponseWriter, r *http.Request) {
	trace := otel.Tracer("servico A")
	ctx, span := trace.Start(r.Context(), "--------> Início processamento serviço A <--------")
	defer span.End()

	cep := chi.URLParam(r, "cep")

	if cepvalido := len(cep) != 8; cepvalido {
		http.Error(w, " invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	ctx, span1 := trace.Start(ctx, "-------> Início consumo da api serviço B <---------------")
	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	req, _ := http.NewRequestWithContext(ctx, "GET", "http://webapi02:8080/consulta-temperatura-por-cep?cep="+cep, nil)

	response, err := client.Do(req)
	span1.End()
	time.Sleep(5 * time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var ResponseDto ResponseDto
	err = json.Unmarshal(body, &ResponseDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponseDto)
	return

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
			semconv.ServiceNameKey.String("webapi01"),
		)),
	)

	otel.SetTracerProvider(tp)
	return tp
}
