package manejador

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

func Iniciar() Manejador {
	fmt.Println("Introduce la URL principal del bandcamp:")
	var bandcampUrl string
	fmt.Scanln(&bandcampUrl)

	url, urlError := url.Parse(bandcampUrl)
	if urlError != nil {
		log.Fatalln(urlError)
	}

	bandcampUrl = "https://" + url.Host

	return Manejador{
		Url:     bandcampUrl,
		UrlInfo: url,
		Cliente: http.Client{
			Timeout: 7 * time.Second,
		},
	}
}
