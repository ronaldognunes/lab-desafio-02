package cep

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/ronaldognunes/lab-desafio-01/internal/entity"
)

type CepService interface {
	ConsultarCep(cep string) (entity.Cep, error)
}

type CepServiceImpl struct {
	url string
}

func NewCepService(url string) CepService {
	return &CepServiceImpl{url: url}
}

func (c *CepServiceImpl) ConsultarCep(cep string) (entity.Cep, error) {
	tratarCep := url.QueryEscape(cep)
	response, err := http.Get(c.url + tratarCep + "/json/")
	fmt.Println(c.url + tratarCep + "/json/")
	if err != nil {
		return entity.Cep{}, err
	}
	fmt.Println(err)
	body, err := io.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return entity.Cep{}, err
	}
	var cepResponse entity.Cep
	json.Unmarshal(body, &cepResponse)
	return cepResponse, nil
}
