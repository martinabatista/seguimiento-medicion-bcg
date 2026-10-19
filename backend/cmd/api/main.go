// Comando api levanta el servidor HTTP del backend.
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/api"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port
	log.Printf("API escuchando en http://localhost%s", addr)
	if err := http.ListenAndServe(addr, api.NewRouter()); err != nil {
		log.Fatal(err)
	}
}
