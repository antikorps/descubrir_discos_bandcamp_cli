package manejador

import (
	"encoding/xml"
	"net/http"
	"net/url"
)

type Manejador struct {
	Azar                bool
	Cliente             http.Client
	Discos              []InfoDisco
	DiscosSeleccionados []string
	Url                 string
	UrlInfo             *url.URL
	Errores             []error
	Canciones           []Cancion
	RutaLista           string
}

type InfoDisco struct {
	Titulo string
	Url    string
}

type CanalCanciones struct {
	Error     error
	Canciones []Cancion
}

type Cancion struct {
	grupo  string
	disco  string
	titulo string
	numero int
	mp3    string
}

type BandcampTrackAlbum struct {
	Current struct {
		Title       string `json:"title"`
		PublishDate string `json:"publish_date"`
		Artist      string `json:"artist"`
		About       string `json:"about"`
	} `json:"current"`
	Trackinfo []struct {
		File struct {
			Mp3128 string `json:"mp3-128"`
		} `json:"file"`
		Title string `json:"title"`
	} `json:"trackinfo"`
}

type ListaReproduccion struct {
	XMLName   xml.Name `xml:"playlist"`
	Text      string   `xml:",chardata"`
	Version   string   `xml:"version,attr"`
	Xmlns     string   `xml:"xmlns,attr"`
	TrackList struct {
		Text  string  `xml:",chardata"`
		Track []Track `xml:"track"`
	} `xml:"trackList"`
}

type Track struct {
	Text     string `xml:",chardata"`
	Location string `xml:"location"`
	Creator  string `xml:"creator"`
	Album    string `xml:"album"`
	Title    string `xml:"title"`
}
