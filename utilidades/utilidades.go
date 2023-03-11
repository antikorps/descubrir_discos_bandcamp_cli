package utilidades

import (
	"net/http"
	"regexp"
	"strings"
)

func IncorporarCabeceras(r *http.Request) {
	cabeceras := map[string]string{
		"User-Agent":                "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/109.0",
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
		"Accept-Language":           "en-US,en;q=0.5",
		"Upgrade-Insecure-Requests": "1",
		"Sec-Fetch-Dest":            "document",
		"Sec-Fetch-Mode":            "navigate",
		"Sec-Fetch-Site":            "none",
		"Sec-Fetch-User":            "?1",
	}
	for clave, valor := range cabeceras {
		r.Header.Add(clave, valor)
	}
}

func MinificarCadena(cadena string) string {
	modificaciones := strings.ReplaceAll(cadena, "\n", "")
	modificaciones = strings.ReplaceAll(modificaciones, "\r", "")
	modificaciones = strings.ReplaceAll(modificaciones, "\t", "")
	expRegEspacios := regexp.MustCompile(`\s+`)
	modificaciones = expRegEspacios.ReplaceAllString(modificaciones, " ")
	modificaciones = strings.TrimSpace(modificaciones)
	return modificaciones
}
