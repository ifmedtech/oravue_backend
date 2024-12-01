package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"log/slog"
	"net/http"
	"oravue_backend/db"
	"oravue_backend/internal/config"
	"oravue_backend/internal/http/handlers/user"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// load config
	cfg := config.MustLoad()

	//database setup
	database, errDb := db.New(cfg)
	if errDb != nil {
		log.Fatalf("unable to connect database %s", errDb)
	}
	slog.Info("storage initialized database")

	//setup router
	router := mux.NewRouter()
	api := router.PathPrefix("/api/v1").Subrouter()

	userRepo := &user.UserRepoStruct{Db: database}
	user.Routes(api, userRepo, cfg)

	//setup server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}
	fmt.Printf("server started %s", cfg.Addr)

	//grace-full shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-done

	slog.Info("shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	err := server.Shutdown(ctx)

	if err != nil {
		slog.Error("Failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("server stopped")
}
