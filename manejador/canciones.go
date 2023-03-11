package manejador

import (
	"descubrir_discos_bandcamp_cli/utilidades"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

func analizarDisco(cliente http.Client, url string, canal chan CanalCanciones) {
	defer wg.Done()
	peticion, peticionError := http.NewRequest("GET", url, nil)
	if peticionError != nil {
		canal <- CanalCanciones{
			Error: peticionError,
		}
		return
	}
	utilidades.IncorporarCabeceras(peticion)
	respuesta, respuestaError := cliente.Do(peticion)
	if respuestaError != nil {
		canal <- CanalCanciones{
			Error: respuestaError,
		}
		return
	}
	defer respuesta.Body.Close()

	if respuesta.StatusCode != 200 {
		mensajeError := fmt.Sprintf("la url %v ha devuelto un status code incorrecto: %v", url, respuesta.Status)
		canal <- CanalCanciones{
			Error: errors.New(mensajeError),
		}
		return
	}

	contenido, contenidoError := io.ReadAll(respuesta.Body)
	if contenidoError != nil {
		canal <- CanalCanciones{
			Error: contenidoError,
		}
		return
	}

	codigoHTML := string(contenido)
	codigoHTML = utilidades.MinificarCadena(codigoHTML)

	expRegTralbum := regexp.MustCompile(`.*?data-tralbum="(.*?)".*`)
	infoAlbum := expRegTralbum.ReplaceAllString(codigoHTML, "$1")
	infoAlbum = strings.ReplaceAll(infoAlbum, "&quot;", `"`)

	var bandcampTrackAlbum BandcampTrackAlbum
	jsonError := json.Unmarshal([]byte(infoAlbum), &bandcampTrackAlbum)
	if jsonError != nil {
		canal <- CanalCanciones{
			Error: jsonError,
		}
		return
	}

	var recopilacionCanciones []Cancion
	for _, v := range bandcampTrackAlbum.Trackinfo {
		nuevaCancion := Cancion{
			grupo:  bandcampTrackAlbum.Current.Artist,
			disco:  bandcampTrackAlbum.Current.Title,
			titulo: v.Title,
			mp3:    v.File.Mp3128,
		}

		if nuevaCancion.mp3 == "" {
			continue
		}

		recopilacionCanciones = append(recopilacionCanciones, nuevaCancion)
	}

	canal <- CanalCanciones{
		Canciones: recopilacionCanciones,
	}

}

func (m *Manejador) BuscarCanciones() {

	canal := make(chan CanalCanciones)

	for _, v := range m.DiscosSeleccionados {
		wg.Add(1)
		go analizarDisco(m.Cliente, v, canal)
	}

	go func() {
		wg.Wait()
		close(canal)
	}()

	for v := range canal {
		if v.Error != nil {
			m.Errores = append(m.Errores, v.Error)
			continue
		}
		m.Canciones = append(m.Canciones, v.Canciones...)
	}

	if len(m.Canciones) == 0 {
		log.Fatalln("no se ha encontrado ninguna canción tras el análisis")
	}

}

func (m *Manejador) Aleatorizar() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(m.Canciones), func(i, j int) { m.Canciones[i], m.Canciones[j] = m.Canciones[j], m.Canciones[i] })
}
