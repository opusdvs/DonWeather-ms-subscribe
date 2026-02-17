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
	"strconv"
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
	"github.com/redis/go-redis/v9"
)

func main() {
	appCtx, appCancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer appCancel()

	redisTTL := os.Getenv("REDIS_TTL")
	if redisTTL == "" {
		log.Fatal("REDIS_TTL environment variable is required")
	}
	redisTTLInt, err := strconv.Atoi(redisTTL)
	if err != nil {
		log.Fatal("REDIS_TTL environment variable is required")
	}
	redisTTLDuration := time.Duration(redisTTLInt) * time.Minute
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		log.Fatal("REDIS_HOST environment variable is required")
	}
	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		log.Fatal("REDIS_PORT environment variable is required")
	}
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
	if err := db.PingContext(appCtx); err != nil {
		log.Fatal("Failed to connect to the database: %w", err)
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", redisHost, redisPort),
		DB:   0,
	})
	defer redisClient.Close()
	if err := redisClient.Ping(appCtx).Err(); err != nil {
		log.Fatal("Failed to connect to the redis: %w", err)
	}
	fmt.Println("Connected to the database and redis")
	if err := RunMigrations(dsn); err != nil {
		log.Fatal(err)
	}

	subscribeRepository := repository.NewPostgresqlSubscribeRepository(db)
	pendingSubscribeRepository := repository.NewRedisSubscribeRepository(redisClient, redisTTLDuration)

	subscribeService := usecase.NewSubscribeService(subscribeRepository, pendingSubscribeRepository)
	subscribeHandlers := delivery.NewSubscribeHandlers(*subscribeService, appCtx)

	healthHandler := delivery.NewHealthHandler(db)
	healthMux := http.NewServeMux()
	healthMux.HandleFunc("/health/liveness", healthHandler.LivenessProbe)
	healthMux.HandleFunc("/health/readiness", healthHandler.ReadinessProbe)

	apiMux := http.NewServeMux()
	apiMux.HandleFunc("/api/v1/subscribe/create", subscribeHandlers.CreateSubscribe)
	apiMux.HandleFunc("/api/v1/subscribe/create/pending", subscribeHandlers.CreatePendingSubscribe)
	handlerMiddleware := middleware.MiddlewareChain(apiMux, middleware.TraceMiddleware, middleware.CorsMiddleware)

	mainMux := http.NewServeMux()
	mainMux.Handle("/", handlerMiddleware)
	mainMux.Handle("/health/", healthMux)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mainMux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func(s *http.Server) {
		err := s.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
			return
		}
		fmt.Println("Server listen and serve success")
	}(server)

	fmt.Println("Server started on port 8080")
	<-appCtx.Done()

	fmt.Println("Server stopped")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Failed to shutdown server: %w", ctx, err)
	}
	fmt.Println("Server shutdown successfully")
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
