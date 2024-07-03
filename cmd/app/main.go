package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"

	v1 "spycat/api/v1"
	"spycat/internal/cats"
	"spycat/internal/missions"
	"spycat/internal/repository"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Application failed: %v", err)
	}
}

func run() error {
	cfg, err := LoadConfig("./config/config.json")
	if err != nil {
		return err
	}

	db, err := sql.Open("postgres", cfg.Database.ConnectionString)
	if err != nil {
		return err
	}
	defer db.Close()

	if err = goose.Up(db, "./migrations"); err != nil {
		log.Fatal("Migration failed:", err)
	}

	catsRepo := repository.NewCatsRepository(db)
	missionRepo := repository.NewMissionRepository(db)
	targetsRepo := repository.NewTargetRepository(db)

	catsService := cats.NewCatService(catsRepo)
	missionsService := missions.NewMissionService(missionRepo, targetsRepo)

	apiV1 := v1.NewAPI(
		catsService,
		missionsService,
	)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: apiV1.SetupRouter(),
	}

	log.Printf("Starting server on %s", server.Addr)
	return server.ListenAndServe()
}
