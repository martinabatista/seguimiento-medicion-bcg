// Comando api levanta el servidor HTTP del backend.
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/api"
	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
	"github.com/martinabatista/seguimiento-medicion-bcg/backend/postgres"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	url := os.Getenv("DATABASE_URL")
	if url == "" {
		log.Fatal("falta la variable DATABASE_URL (ver .env.example)")
	}
	db, err := postgres.Abrir(url)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	proyectos := dominio.NuevoServicioProyectos(postgres.NuevoRepositorioProyectos(db))

	addr := ":" + port
	log.Printf("API escuchando en http://localhost%s", addr)
	if err := http.ListenAndServe(addr, api.NewRouter(proyectos)); err != nil {
		log.Fatal(err)
	}
}
