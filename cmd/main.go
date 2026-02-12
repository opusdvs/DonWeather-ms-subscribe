package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	delivery "github.com/opusdvs/DonWeather-ms-subscribe/internal/delivery/http"
	"github.com/opusdvs/DonWeather-ms-subscribe/internal/delivery/middleware"
	"github.com/opusdvs/DonWeather-ms-subscribe/internal/repository"
	"github.com/opusdvs/DonWeather-ms-subscribe/internal/usecase"
)

func main() {
	pingCtx, pingCancel := context.WithTimeout(context.Background(), 10*time.Second)
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		log.Fatal("DB_PASSWORD environment variable is required")
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		log.Fatal("DB_USER environment variable is required")
	}
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		log.Fatal("DB_HOST environment variable is required")
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		log.Fatal("DB_PORT environment variable is required")
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		log.Fatal("DB_NAME environment variable is required")
	}
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	defer pingCancel()
	if err := db.PingContext(pingCtx); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to the database")

	subscribeRepository := repository.NewPostgresqlSubscribeRepository(db)
	subscribeService := usecase.NewSubscribeService(subscribeRepository)
	subscribeHandlers := delivery.NewSubscribeHandlers(*subscribeService)
	if err := RunMigrations(dsn); err != nil {
		log.Fatal(err)
	}

	healthHandler := delivery.NewHealthHandler(db)
	healthMux := http.NewServeMux()
	healthMux.HandleFunc("/health/liveness", healthHandler.LivenessProbe)
	healthMux.HandleFunc("/health/readiness", healthHandler.ReadinessProbe)

	apiMux := http.NewServeMux()
	apiMux.HandleFunc("/api/v1/subscribe-create", subscribeHandlers.CreateSubscribe)
	apiMux.HandleFunc("/api/v1/subscribe-get-all", subscribeHandlers.GetAllSubscribes)
	apiMux.HandleFunc("/api/v1/subscribe-get-by-id", subscribeHandlers.GetSubscribeById)
	apiMux.HandleFunc("/api/v1/subscribe-update", subscribeHandlers.UpdateSubscribe)
	apiMux.HandleFunc("/api/v1/subscribe-delete", subscribeHandlers.DeleteSubscribe)
	apiMux.HandleFunc("/api/v1/set-telegram-id", subscribeHandlers.SetTelegramID)
	handlerMiddleware := middleware.MiddlewareChain(apiMux, middleware.TraceMiddleware, middleware.CorsMiddleware)

	mainMux := http.NewServeMux()
	mainMux.Handle("/", handlerMiddleware)
	mainMux.Handle("/health", healthMux)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mainMux,
	}

	go func(s *http.Server) {
		err := s.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
			return
		}
		fmt.Println("Server listen and serve success")
	}(server)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Server stopped")
}

func RunMigrations(dsn string) error {
	m, err := migrate.New(
		"file://migrations",
		dsn,
	)
	if err != nil {
		return fmt.Errorf("create migrate instance: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}
		return fmt.Errorf("run migrations: %w", err)
	}

	return nil
}
