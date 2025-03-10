package webserver

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ronaldognunes/lab-desafio-01/internal/infra/service/cep"
	"github.com/ronaldognunes/lab-desafio-01/internal/infra/service/temperatura"
	"github.com/stretchr/testify/assert"
)

func TestConsultarCepHandler(t *testing.T) {
	servicoCep := cep.NewCepService("https://viacep.com.br/ws/")
	servicoTemperatura := temperatura.NewTemperaturaService("http://api.weatherapi.com/v1/current.json?q=", "6c71bade554742f0b46142657250703")

	handler := NewConsultaHandler(servicoCep, servicoTemperatura)

	t.Run("Deve retornar 422 se Cep for inválido", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/consulta-temperatura-por-cep?cep=2503553244", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler.ConsultarCepHandler(rr, req)

		assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
		assert.Equal(t, "invalid zipcode\n", rr.Body.String())
	})

	t.Run("deve retornar 404 quando cep não for encontrado ou não existir", func(t *testing.T) {

		req, err := http.NewRequest("GET", "/consulta-temperatura-por-cep?cep=25035531", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler.ConsultarCepHandler(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Equal(t, "can not find zipcode\n", rr.Body.String())
	})

	t.Run("deve retornar 200 em caso de sucesso", func(t *testing.T) {

		req, err := http.NewRequest("GET", "/consulta-temperatura-por-cep?cep=24921692", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler.ConsultarCepHandler(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var response ResponseDto
		err = json.NewDecoder(rr.Body).Decode(&response)
		assert.NoError(t, err)
		assert.NotNil(t, response.TempC)
		assert.NotNil(t, response.TempF)
		assert.NotNil(t, response.TempK)
	})

}
