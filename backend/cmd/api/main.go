// Comando api levanta el servidor HTTP del backend.
package main

import (
	"log"
	"net/http"
	"os"
	"time"

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

	repoProyectos := postgres.NuevoRepositorioProyectos(db)
	repoHistorias := postgres.NuevoRepositorioHistorias(db)
	repoSprints := postgres.NuevoRepositorioSprints(db)
	repoEsfuerzo := postgres.NuevoRepositorioEsfuerzo(db)
	proyectos := dominio.NuevoServicioProyectos(repoProyectos)
	backlog := dominio.NuevoServicioBacklog(repoProyectos, repoHistorias, repoSprints, repoEsfuerzo)
	repoIntegrantes := postgres.NuevoRepositorioIntegrantes(db)
	integrantes := dominio.NuevoServicioIntegrantes(repoProyectos, repoIntegrantes)
	// La fecha de hoy del estado del proyecto (spec 03) es la fecha local del servidor.
	estado := dominio.NuevoServicioEstado(repoProyectos, repoHistorias, repoSprints, time.Now)
	sprints := dominio.NuevoServicioSprints(repoProyectos, repoHistorias, repoSprints)
	poker := dominio.NuevoServicioPlanningPoker(repoProyectos, repoHistorias, repoIntegrantes, postgres.NuevoRepositorioVotaciones(db), backlog)
	metricas := dominio.NuevoServicioMetricas(repoProyectos, repoHistorias, repoSprints)
	// La fecha de hoy del registro de esfuerzo (spec 09) también es la fecha local del servidor.
	esfuerzo := dominio.NuevoServicioEsfuerzo(repoProyectos, repoHistorias, repoIntegrantes, repoEsfuerzo, time.Now)
	// El Dashboard (spec 13) reutiliza el estado y las métricas, con la misma fecha de hoy.
	dashboard := dominio.NuevoServicioDashboard(estado, metricas, time.Now)

	addr := ":" + port
	log.Printf("API escuchando en http://localhost%s", addr)
	if err := http.ListenAndServe(addr, api.NewRouter(proyectos, backlog, integrantes, estado, sprints, poker, metricas, esfuerzo, dashboard)); err != nil {
		log.Fatal(err)
	}
}
