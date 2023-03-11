package manejador

import (
	"descubrir_discos_bandcamp_cli/utilidades"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"

	"atomicgo.dev/keyboard/keys"
	"github.com/PuerkitoBio/goquery"
	"github.com/pterm/pterm"
)

func buscarDiscosEnPaginaMusic(cliente http.Client, urlInfo *url.URL) []InfoDisco {
	var coleccionDiscos []InfoDisco

	urlMusic := "https://" + urlInfo.Host + "/music"

	peticion, peticionError := http.NewRequest("GET", urlMusic, nil)
	if peticionError != nil {
		return coleccionDiscos
	}
	utilidades.IncorporarCabeceras(peticion)

	respuesta, respuestaError := cliente.Do(peticion)
	if respuestaError != nil {
		return coleccionDiscos
	}
	defer respuesta.Body.Close()
	if respuesta.StatusCode != 200 {
		return coleccionDiscos
	}

	documento, documentoError := goquery.NewDocumentFromReader(respuesta.Body)
	if documentoError != nil {
		log.Fatalln(documentoError)
	}

	documento.Find("li.music-grid-item").Each(func(i int, s *goquery.Selection) {
		var infoDisco InfoDisco

		url, urlExistencia := s.Find("a").Attr("href")
		if !urlExistencia {
			return
		}
		infoDisco.Url = "urlPrincipal" + url

		titulo := s.Find(".title").Text()
		if titulo == "" {
			return
		}
		infoDisco.Titulo = utilidades.MinificarCadena(titulo)

		coleccionDiscos = append(coleccionDiscos, infoDisco)
	})

	return coleccionDiscos

}

func buscarDiscosEnPaginaPrincipal(documento *goquery.Document, urlPrincipal string) []InfoDisco {

	var coleccionDiscos []InfoDisco

	documento.Find("li.music-grid-item").Each(func(i int, s *goquery.Selection) {
		var infoDisco InfoDisco

		url, urlExistencia := s.Find("a").Attr("href")
		if !urlExistencia {
			return
		}
		infoDisco.Url = urlPrincipal + url

		titulo := s.Find(".title").Text()
		if titulo == "" {
			return
		}
		infoDisco.Titulo = utilidades.MinificarCadena(titulo)

		coleccionDiscos = append(coleccionDiscos, infoDisco)
	})

	return coleccionDiscos
}

func buscarDiscosEnPaginaDisco(documento *goquery.Document, urlPrincipal string) []InfoDisco {

	var coleccionDiscos []InfoDisco

	discografia := documento.Find("#discography")
	discografia.Find("li").Each(func(i int, s *goquery.Selection) {
		var infoDisco InfoDisco

		url, urlExistencia := s.Find("a.thumbthumb").Attr("href")
		if !urlExistencia {
			return
		}
		infoDisco.Url = urlPrincipal + url

		titulo := s.Find("div.trackTitle > a").Text()
		if titulo == "" {
			return
		}
		infoDisco.Titulo = utilidades.MinificarCadena(titulo)

		coleccionDiscos = append(coleccionDiscos, infoDisco)

	})

	return coleccionDiscos
}

func (m *Manejador) BuscarDiscos() {
	peticion, peticionError := http.NewRequest("GET", m.Url, nil)
	if peticionError != nil {
		log.Fatalln(peticionError)
	}
	utilidades.IncorporarCabeceras(peticion)

	respuesta, respuestaError := m.Cliente.Do(peticion)
	if respuestaError != nil {
		log.Fatalln(respuestaError)
	}
	defer respuesta.Body.Close()
	if respuesta.StatusCode != 200 {
		mensajeError := fmt.Sprintf("status code de %v incorrecto: %v", m.Url, respuesta.Status)
		log.Fatalln(mensajeError)
	}

	documento, documentoError := goquery.NewDocumentFromReader(respuesta.Body)
	if documentoError != nil {
		log.Fatalln(documentoError)
	}

	var coleccionDiscos []InfoDisco

	/*
		Con la página principal de una cuenta de tipo artista suele mostrarse en mosaico todas las ediciones
		-> buscarDiscosEnPaginaPrincipal
		Ej:	https://newretrowave.bandcamp.com/

		Hay páginas que redirigen automáticamente a una edición
		-> buscarDiscosEnPaginaDisco
		Ej: https://lamoda.bandcamp.com/"

		Las cuentas de sellos se redirigen a /music?
		-> buscarDiscosEnPaginaMusic
		Ej: https://relapserecords.bandcamp.com, https://hellsheadbangers.bandcamp.com/

	*/

	coleccionDiscos = buscarDiscosEnPaginaPrincipal(documento, m.Url)
	if len(coleccionDiscos) == 0 {
		coleccionDiscos = buscarDiscosEnPaginaDisco(documento, m.Url)
	}

	if len(coleccionDiscos) == 0 {
		coleccionDiscos = buscarDiscosEnPaginaMusic(m.Cliente, m.UrlInfo)
	}

	if len(coleccionDiscos) == 0 {
		log.Fatalln("no se ha encontrado ningún disco en la URL introducida")
	}

	m.Discos = coleccionDiscos

}

func (m *Manejador) EscogerDiscos() {
	var opciones []string

	for i, v := range m.Discos {
		indiceHumano := i + 1
		texto := fmt.Sprintf("%d - %v", indiceHumano, v.Titulo)
		opciones = append(opciones, texto)
	}

	pintarOpciones := pterm.DefaultInteractiveMultiselect.WithOptions(opciones)
	pintarOpciones.DefaultText = "\nEscoge los discos que se analizarán para crear la lista de reproduccion"
	pintarOpciones.Filter = false
	pintarOpciones.KeyConfirm = keys.Enter
	pintarOpciones.KeySelect = keys.Space
	pintarOpciones.MaxHeight = 20
	opcionesEscogidas, opcionesEsogidasError := pintarOpciones.Show()

	if opcionesEsogidasError != nil {
		log.Fatalln(opcionesEsogidasError)
	}

	var discosAnalizar []string
	for _, v := range opcionesEscogidas {
		expRegIndice := regexp.MustCompile(`(^[0-9]{1,3}).*`)
		indiceOpcion := expRegIndice.ReplaceAllString(v, "$1")
		indice, indiceError := strconv.Atoi(indiceOpcion)
		if indiceError != nil {
			fmt.Println("se omite la elección del disco", v, "por no haber podido localizar el indice")
			continue
		}
		discosAnalizar = append(discosAnalizar, m.Discos[indice-1].Url)
	}

	if len(discosAnalizar) == 0 {
		log.Println("Ningún discos seleccionado")
		os.Exit(0)
	}

	m.DiscosSeleccionados = discosAnalizar

	aleatoridad := pterm.DefaultInteractiveConfirm
	aleatoridad.ConfirmText = "s"
	aleatoridad.RejectText = "n"
	aleatoridad.DefaultText = "¿Quieres que las canciones se ordenen aleatoriamente?"
	aleatoridad.DefaultValue = true

	pintarAleatoridad, pintarAleatoridadError := aleatoridad.Show()
	if pintarAleatoridadError != nil {
		log.Fatalln(pintarAleatoridad)
	}

	m.Azar = pintarAleatoridad
}
