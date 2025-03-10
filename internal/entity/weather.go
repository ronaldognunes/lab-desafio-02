package entity

import (
	"unicode"

	"golang.org/x/text/unicode/norm"
)

type Weather struct {
	Current Current `json:"current"`
}

type Current struct {
	TempC float64 `json:"temp_c"`
	TempF float64
	TempK float64
}

func (c *Current) CalcularTemperaturas() {
	if c.TempC == 0 {
		return
	}
	c.TempF = c.TempC*1.8 + 32
	c.TempK = c.TempC + 273
}

func RemoveAcentos(s string) string {
	t := norm.NFD.String(s)
	result := make([]rune, 0, len(t))
	for _, r := range t {
		if unicode.IsMark(r) {
			continue
		}
		result = append(result, r)
	}
	return string(result)
}
