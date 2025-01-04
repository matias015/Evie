package utils

import (
	"log"
	"os"
)

func ReadFile(path string) string {
	datosComoBytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	// convertir el arreglo a string
	datosComoString := string(datosComoBytes)
	// imprimir el string
	return datosComoString
}

func AddExtension(filename string) string {
	return filename + ".ev"
}
