package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/2024-CMPU9010-GROUP-3/magpie/internal/handlers"
	"github.com/2024-CMPU9010-GROUP-3/magpie/internal/middleware"
	"github.com/2024-CMPU9010-GROUP-3/magpie/internal/routes"
	"github.com/2024-CMPU9010-GROUP-3/magpie/internal/util/dbutil"
	"github.com/2024-CMPU9010-GROUP-3/magpie/internal/util/env"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/cors"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	env.Load()

	err := env.Validate()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	port, _ := env.Get(env.EnvPort)

	allowedOriginsEnv, _ := env.Get(env.EnvCorsOrigins)
	allowedMethodsEnv, _ := env.Get(env.EnvCorsMethods)

	allowedOrigins := strings.Split(allowedOriginsEnv, " ")
	allowedMethods := strings.Split(allowedMethodsEnv, " ")

	corsManager := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   allowedMethods,
		AllowCredentials: true, // needed for login
	})

	middlewares := middleware.CreateMiddlewareStack(
		middleware.Logging,
		corsManager.Handler,
	)
	router := routes.Router

	dbPool, err := dbutil.SetupDatabase()
	if err != nil {
		log.Fatalf("Could not connect to database: %v\n", err)
		os.Exit(1)
	} else {
		log.Println("Successfully connected to database")
	}
	defer dbPool.Close()

	ctx := context.Background()

	handlers.RegisterDatabaseConnection(&ctx, dbPool)

	server := http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: middlewares(router),
	}

	// run server in separate thread for graceful shutdown
	go func() {
		log.Printf("Listening on port %v\n", port)
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %+v", err)
		}
		log.Println("Commencing shutdown...")
	}()

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
	// the following line pauses the execution of this thread until a signal is received
	// the server keeps running in a different thread
	<-signalChannel

	// tell server to shut down in time for the shutdown timeout at the latest
	// no new connections are accepted, existing connections will still be handled
	shutdownTimeout, _ := env.GetInt(env.EnvShutdownTimeout)
	shutdownContext, shutdownCancelFunc := context.WithTimeout(context.Background(), time.Duration(shutdownTimeout)*time.Second)
	defer shutdownCancelFunc()

	if err := server.Shutdown(shutdownContext); err != nil {
		log.Fatalf("Error during shutdown procedure: %+v", err)
	}
	log.Println("Shutdown complete")
}
