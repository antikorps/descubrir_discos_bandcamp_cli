package manejador

import "fmt"

func (m *Manejador) Resumir() {

	if len(m.Errores) > 0 {
		fmt.Println(`
		
SE HAN PRODUCIDO LOS SIGUIENTES ERRORES NO CRÍTICOS:
====================================================`)

		for _, v := range m.Errores {
			fmt.Println(v.Error())
		}
	}

	fmt.Printf(`
		
FIN:
====
Lista de reproducción creada en: %v
`, m.RutaLista)

}
