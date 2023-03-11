package manejador

import (
	"encoding/xml"
	"io"
	"log"
	"os"
	"path/filepath"
)

func (m *Manejador) CrearListaReproduccion() {
	// XML Shareable Playlist Format (XSPF)
	// JSON https://www.xspf.org/jspf

	var listaReproduccion ListaReproduccion
	listaReproduccion.Version = "1"
	listaReproduccion.Xmlns = "http://xspf.org/ns/0/"

	for _, v := range m.Canciones {
		var track Track
		track.Album = v.disco
		track.Creator = v.grupo
		track.Location = v.mp3
		track.Title = v.titulo
		listaReproduccion.TrackList.Track = append(listaReproduccion.TrackList.Track, track)
	}

	rutaEjecucion, rutaEjecucionError := os.Getwd()
	if rutaEjecucionError != nil {
		log.Fatalln(rutaEjecucionError)
	}

	rutaLista := filepath.Join(rutaEjecucion, "playlist_bandcamp.xspf")
	m.RutaLista = rutaLista

	archivoJSPF, archivoJSPFError := os.Create(rutaLista)
	if archivoJSPFError != nil {
		log.Fatalln(archivoJSPFError)
	}
	defer archivoJSPF.Close()

	listaJson, listaJsonError := xml.Marshal(listaReproduccion)
	if listaJsonError != nil {
		log.Fatalln(listaJsonError)
	}

	_, errorEscritura := io.WriteString(archivoJSPF, string(listaJson))
	if errorEscritura != nil {
		log.Fatalln(errorEscritura)
	}
}
