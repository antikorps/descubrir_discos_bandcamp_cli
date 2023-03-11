package main

import "descubrir_discos_bandcamp_cli/manejador"

func main() {
	manejadorDescargas := manejador.Iniciar()
	manejadorDescargas.BuscarDiscos()
	manejadorDescargas.EscogerDiscos()
	manejadorDescargas.BuscarCanciones()
	if manejadorDescargas.Azar {
		manejadorDescargas.Aleatorizar()
	}
	manejadorDescargas.CrearListaReproduccion()
	manejadorDescargas.Resumir()
}
