package temperatura

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/ronaldognunes/lab-desafio-01/internal/entity"
)

type TemperaturaService interface {
	ConsultarTemperatura(localidade string) (entity.Current, error)
}

type TemperaturaServiceImpl struct {
	Url string
	key string
}

func NewTemperaturaService(url string, key string) TemperaturaService {
	return &TemperaturaServiceImpl{Url: url, key: key}
}
func (c *TemperaturaServiceImpl) ConsultarTemperatura(localidade string) (entity.Current, error) {
	// 6c71bade554742f0b46142657250703
	// Consulta temperatura
	tratarLocalidade := url.QueryEscape(localidade)

	response, err := http.Get(c.Url + tratarLocalidade + "&key=" + c.key)
	if err != nil {
		return entity.Current{}, err
	}

	body, err := io.ReadAll(response.Body)
	defer response.Body.Close()

	if err != nil {
		return entity.Current{}, err
	}
	var temperatura entity.Weather
	json.Unmarshal(body, &temperatura)
	return temperatura.Current, nil
}
