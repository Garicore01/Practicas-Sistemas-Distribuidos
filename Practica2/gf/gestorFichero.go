package gf

import (
	"fmt"
	"io/ioutil"
	"os"
)

func LeerFichero(fichero string) (string){
	dat, err := ioutil.ReadFile(fichero)
	if err != nil {
		fmt.Println("Error al leer el fichero")
		os.Exit(1)
	}
	return string(dat)
}


func EscribirFichero(fichero string, contenido string){
	f, err :=os.OpenFile(fichero, os.O_APPEND | os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Error al abrir el fichero")
		os.Exit(1)
	}
	defer f.Close()
	_, err = f.WriteString(contenido)
	if err != nil {
		fmt.Println("Error al escribir en el fichero")
		os.Exit(1)
	}
}

func CrearFichero(fichero string){
	_, err := os.Create(fichero)
	if err != nil {
		fmt.Println("Error al crear el fichero")
		os.Exit(1)
	}

}