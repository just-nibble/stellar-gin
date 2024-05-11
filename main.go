package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"bitgifty.com/stellar/database"
	"bitgifty.com/stellar/routers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type config struct {
	port int
	env  string
}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func runMigration() {

}

func loadDatabase() {
	database.Connect()
	// runMigration()
}

// @host      0.0.0.0:8080
// @BasePath  /api/v1

func main() {
	loadEnv()

	loadDatabase()

	var cfg config
	flag.IntVar(&cfg.port, "port", 8000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()
	// Logging to a file.
	var f, _ = os.Create("gin.log")
	// Use the following code if you need to write the logs to file and console at the same time.
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	r := routers.NewRouter()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.port),
		Handler: r,
	}

	go func() {
		// service connections
		log.Printf("starting %s server on %s", cfg.env, srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.

	<-ctx.Done()
	log.Println("timeout of 2 seconds.")
	log.Println("Server exiting")
}
