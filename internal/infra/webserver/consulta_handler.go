package webserver

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ronaldognunes/lab-desafio-01/internal/entity"
	"github.com/ronaldognunes/lab-desafio-01/internal/infra/service/cep"
	"github.com/ronaldognunes/lab-desafio-01/internal/infra/service/temperatura"
	"go.opentelemetry.io/otel"
)

type ResponseDto struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_c"`
	TempF float64 `json:"temp_f"`
	TempK float64 `json:"temp_k"`
}

type ConsultaHandler struct {
	CepService         cep.CepService
	TemperaturaService temperatura.TemperaturaService
}

func NewConsultaHandler(cepService cep.CepService, temperaturaService temperatura.TemperaturaService) *ConsultaHandler {
	return &ConsultaHandler{CepService: cepService, TemperaturaService: temperaturaService}
}

func (c *ConsultaHandler) ConsultarCepHandler(w http.ResponseWriter, r *http.Request) {
	cep := r.URL.Query().Get("cep")
	if cep == "" {
		http.Error(w, "CEP não informado", http.StatusBadRequest)
		return
	}

	if cepValido := entity.ValidaCEP(cep); cepValido == false {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	fmt.Println("Consultando CEP")
	trace := otel.Tracer("servico B")
	_, span := trace.Start(r.Context(), "-------> Início consulta api CPF <---------------") // span consulta api cep

	dadosCep, err := c.CepService.ConsultarCep(cep)
	if err != nil {
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}
	if dadosCep.Cep == "" {
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}

	span.End() // span consulta api cep fim
	fmt.Println("Localidade:", dadosCep.Localidade)

	_, span2 := trace.Start(r.Context(), "----------> Início consulta api temperatura<-----------") // span consulta api cep
	temp, err := c.TemperaturaService.ConsultarTemperatura(entity.RemoveAcentos(dadosCep.Localidade))

	if err != nil {
		http.Error(w, "Erro ao consultar temperatura", http.StatusBadRequest)
		return
	}

	if temp.TempC == 0 {
		http.Error(w, "Temperatura não encontrada", http.StatusBadRequest)
		return
	}

	fmt.Println("Temperatura:", temp)
	span2.End() // span consulta api temperatura fim

	temp.CalcularTemperaturas()
	response := ResponseDto{
		City:  dadosCep.Localidade,
		TempC: temp.TempC,
		TempF: temp.TempF,
		TempK: temp.TempK,
	}
	fmt.Println("Temperatura", response)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	return

}
